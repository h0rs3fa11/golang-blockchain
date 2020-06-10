package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"

const BlocksBucket = "blocks"

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
	if blockchain.Params.Miner == "" {
		blockchain.Params.setCoinbase()
	}
	pubkey, err := GetPublickey(blockchain.Params.Miner)
	if err != nil {
		fmt.Println(err)
	}

	txToBlock = append(txToBlock, createCoinbaseTx(HashPubKey(pubkey), "", blockchain.Params))
	for _, transaction := range txs {
		txToBlock = append(txToBlock, transaction)
	}

	var block *Block

	blockchain.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
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

	err = blockchain.Database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

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

// TODO 修复findunspentTX bug
//Find Transaction contains UTXO of address
func (blockchain *Blockchain) FindUnspentTX(address string) []Transaction {
	var unspentTXs []Transaction

	spentTXOs := make(map[string][]int) //存储已花费的UTXO
	bci := blockchain.Iterator()
	var hashInt big.Int

	pubKey, err := GetPublickey(address)
	if err != nil {
		fmt.Println(err)
	}

	for {
		block := bci.Next()

		//遍历区块交易
		for _, transaction := range block.Transaction {
			//fmt.Printf("TransactionHash:%x\n", transaction.ID)
			txID := hex.EncodeToString(transaction.ID)

			//如果当前交易不是coinbase，遍历transaction.Vin，将input添加到spentTXOs中
			if transaction.isCoinbase() == false {
				for _, in := range transaction.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}

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
	unspentTx := blockchain.FindUnspentTX(address)
	accumulated := 0
	pubKey, err := GetPublickey(address)
	if err != nil {
		fmt.Println(err)
	}

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
	if accumulated < amount {
		return -1, nil
	}
	return accumulated, unspentOutputs
}

//find All UTXO
func (blockchain *Blockchain) findAllUnspentOutput(address string) int {
	unspentTx := blockchain.FindUnspentTX(address)
	accumulated := 0
	pubKey, err := GetPublickey(address)
	if err != nil {
		fmt.Println(err)
	}

	for _, tx := range unspentTx {

		for _, out := range tx.Vout {
			if out.IsLockWithKey(HashPubKey(pubKey)) {
				accumulated += out.Value
			}
		}
	}

	return accumulated
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
			break
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
	params := GetConfig()
	
	wallets, err := NewWallets()
	
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Create a new one ...")

			coinbase := wallets.CreateNewWallet()

			params.Updateparams(params.Miner, coinbase)

			blockchain := Blockchain{nil, nil, *params}
			genesisBlock := NewGenesisBlock(&blockchain)

			b, err := tx.CreateBucket([]byte(BlocksBucket))
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
			//params.setCoinbase()
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db, *params}
}

func DbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
