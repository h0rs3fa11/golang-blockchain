package BLC

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"crypto/ecdsa"
	"github.com/boltdb/bolt"
	"bytes"
	"time"
)

const dbFile = "blockchain.db"

const blocksBucket = "blocks"

type Blockchain struct {
	Tip      []byte
	Database *bolt.DB
	Params   Chainparams
}

func (blockchain *Blockchain) Iterator() *ChainIterator {
	return &ChainIterator{blockchain.Tip, blockchain.Database}
}

//add new block
func (blockchain *Blockchain) AddBlock(txs []*Transaction) {
	var txToBlock []*Transaction

	//create coinbase transction	
	txToBlock = append(txToBlock, createCoinbaseTx(HashPubKey(getPublickey(blockchain.Params.Miner)), "", blockchain.Params))
	for _, transaction := range txs {
		txToBlock = append(txToBlock, transaction)
	}

	var block *Block

	blockchain.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil 
	})

	for _,tx := range txs  {

		if blockchain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	//create new block
	block = NewBlock(txToBlock, block.Hash, block.Height + 1)

	err := blockchain.Database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			log.Panic(err)
		}
		blockchain.Tip = block.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

//Find Transaction contains UTXO of address
func (blockchain *Blockchain) findUnspentTX(address string) []Transaction {
	var unspentTXs []Transaction

	spentTXOs := make(map[string][]int) //存储已花费的UTXO
	bci := blockchain.Iterator()
	var hashInt big.Int

	pubKey := getPublickey(address)
	for {
		block := bci.Next()

		//遍历区块交易
		for _, transaction := range block.Transaction {
			//fmt.Printf("TransactionHash:%x\n", transaction.ID)
			txID := hex.EncodeToString(transaction.ID)

		Outputs:
			//遍历交易输出
			for outIdx, out := range transaction.Vout {
				//检查当前transaction.Vout中的output有没有已花费的(spentTXOs)
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs //跳过当前vout
						}
					}
				}
				//如果已花费就跳过，没有就检查unlock，然后添加到unspentTXs中
				if out.IsLockWithKey(HashPubKey(pubKey)) {
					unspentTXs = append(unspentTXs, *transaction)
				}
			}

			//如果当前交易不是coinbase，遍历transaction.Vin，将input添加到spentTXOs中
			if transaction.isCoinbase() == false {
				for _, in := range transaction.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		hashInt.SetBytes(bci.CurrentHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return unspentTXs
}

//Find UTXO
func (blockchain *Blockchain) findUnspentOutput(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTx := blockchain.findUnspentTX(address)
	accumulated := 0
	pubKey := getPublickey(address)

Work:
	for _, tx := range unspentTx {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockWithKey(HashPubKey(pubKey)) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}

func (blockchain *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	if tx.isCoinbase() {
		return
	}

	prevTxs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTx, err := blockchain.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	tx.Sign(privKey, prevTxs)
}

func (blockchain *Blockchain) FindTransaction(id []byte) (Transaction, error) {
	bci := blockchain.Iterator()

	for {
		block := bci.Next()

		for _,tx := range block.Transaction {
			if bytes.Compare(tx.ID, id) == 0 {
				return *tx, nil
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return Transaction{}, nil
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTxs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTx, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	return tx.Verify(prevTxs)
}

//Create a blockchain with genesis block
func NewBlockChain() *Blockchain {
	var tip []byte

	//暂时 后面写成配置文件
	params := Chainparams{}
	params.init()

	wallets, err := newWallets()
	
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Create a new one ...")

			wallets.createNewWallet()

			params.setCoinbase()

			blockchain := Blockchain{nil, nil, params}
			genesisBlock := NewGenesisBlock(&blockchain)

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			tip = genesisBlock.Hash
		} else {
			params.setCoinbase()
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db, params}
}

func (chain *Blockchain) PrintChain() {

	blockchainIterator := chain.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height:%d\n", block.Height)
		fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
		fmt.Printf("Timestamp：%s \n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x \n", block.Hash)
		fmt.Printf("Nonce：%d \n", block.Nonce)
		for _, tx := range block.Transaction {
			fmt.Printf("\nTransaction id: %x\n", tx.ID)
			//遍历vin
			fmt.Println("----------transaction input----------")
			for _, txin := range tx.Vin {
				fmt.Printf("Vin transaction ID: %x\n", txin.Txid)
				fmt.Printf("Vin Vout: %d\n", txin.Vout)
				//fmt.Printf("Script Sig: %s\n", txin.ScriptSig)
			} 
			fmt.Println("----------transaction output----------")
			//遍历vout
			fmt.Println("Vouts:")
			for _, txout := range tx.Vout {
				fmt.Println(txout.Value)
				fmt.Printf("%s",GetAddressFromPubkey(txout.PubKeyHash))
			}
		}
		fmt.Println()
		fmt.Println("---------------------------------------")
		
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
