package rpcapi

import (
	"blockchain/core"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/boltdb/bolt"
	jsoniter "github.com/json-iterator/go"
)

type Transfer struct {
	From string
	To string
	Amount int
}

func NewTransfer(from string, to string, amount int) (*Transfer) {
	tx := Transfer{
		From: from,
		To: to,
		Amount: amount,
	}
	return &tx
}

func (r *Rpc) Help(arg string, reply *string) error {
	*reply = fmt.Sprintln(`Usage:
	getbalance -address ADDRESS - Get balance of ADDRESS 
	printchain - Print all the blocks of the blockchain
	sendmany -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO
	listaddress list all address from wallet
	createaddress create a new address
	pubkeyhashtoaddress -pubkeyhash PUBLIC_KEY_HASH convert public key to address
	`)
	return nil
}

func (r *Rpc) PrintChain(args string, reply *string) error {
	blockchainIterator := r.bc.Iterator()
	blockInfo := &core.Block{}
	for {
		blockInfo = blockchainIterator.Next()
		
		blockJSON, err := jsoniter.MarshalIndent(blockInfo, "", "    ")
		if err != nil {
			*reply = fmt.Sprintln("Json marshal failed");
			return nil
		}
		*reply += fmt.Sprintln(string(blockJSON))

		var hashInt big.Int
		hashInt.SetBytes(blockInfo.PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return nil
}

func (r *Rpc) Getblock(args string, reply *string) error {

	hashByte, err := hex.DecodeString(args)

	if err != nil {
		log.Panic(err)
	}

	err = r.bc.Database.View(func(tx *bolt.Tx) error {
		// 获取表
		b := tx.Bucket([]byte(core.BlocksBucket))
		// 通过Hash获取区块字节数组
		blockBytes := b.Get(hashByte)

		block := core.DeserializeBlock(blockBytes)

		blockJSON, err := jsoniter.MarshalIndent(block, "", "    ")
		if err != nil {
			*reply = fmt.Sprintln("Json marshal failed");
			return nil
		}
		*reply += fmt.Sprintln(string(blockJSON))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return nil
}

func (r *Rpc) Sendmany(args *Transfer, reply *string) error {
	tx, err := core.CreateTransaction(args.From, args.To, args.Amount, &r.bc, "")
	if err != nil {
		*reply = fmt.Sprintln(err)
	}

	r.bc.AddBlock([]*core.Transaction{tx})
	return nil
}

func (r *Rpc) Getbalance(args string, reply *string) error {
	txs := r.bc.FindUnspentTX(args)
	pubKey, err := core.GetPublickey(args)
	if err != nil {
		*reply = fmt.Sprintln(err)
		return nil
	}
	balance := 0

	for _, tx := range txs {
		for _, out := range tx.Vout {
			if out.IsLockWithKey(core.HashPubKey(pubKey)) {
				balance += out.Value
			}
		}
	}

	*reply = strconv.Itoa(balance)
	return nil
}

func (r *Rpc) Listaddress(args string, reply *string) error {
	wallets, _ := core.NewWallets()
	for address, _ := range wallets.WalletsMap {
		*reply += fmt.Sprintf("%s\n", address)
	}
	return nil
}

func (r *Rpc) Createaddress(args string, reply *string) error {
	wallets, _ := core.NewWallets()
	*reply = wallets.CreateNewWallet()
	
	return nil
}

func (r *Rpc) PubkeyHashToAddress(args string, reply *string) error {
	*reply = string(core.GetAddressFromPubkey([]byte(args)))

	return nil
}