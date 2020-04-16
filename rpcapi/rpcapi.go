package rpcapi

import (
	"blockchain/core"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
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
	`)
	return nil
}

func (r *Rpc) PrintChain(args string, reply *string) error {
	r.bc.PrintChain()
	*reply = fmt.Sprintln("success")
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

		*reply = fmt.Sprintf("PrevBlockHash：%x \nTimestamp：%s \nHash：%x \nNonce：%d \n", block.PrevBlockHash, time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"), block.Hash, block.Nonce)

		for _, tx := range block.Transaction {
			fmt.Printf("Transaction id: %x\n", tx.ID)
			*reply += fmt.Sprintf("Transaction id: %x\n", tx.ID)
			for _, in := range tx.Vin {
				//fmt.Printf("Transaction Input:\ntxid: %x\nout index: %d\nscriptSig:%s\n", in.Txid, in.Vout, in.ScriptSig)
				*reply += fmt.Sprintf("Transaction Input:\ntxid: %x\nout index: %d\nscriptSig:%s\n", in.Txid, in.Vout, in.Signature)
			}
			//fmt.Println()
			for _, out := range tx.Vout {
				//fmt.Printf("Transaction Output:\nvalue: %d\nscriptPubKey: %s\n", out.Value, out.ScriptPubKey)
				*reply += fmt.Sprintf("Transaction Output:\nvalue: %d\nscriptPubKey: %s\n", out.Value, out.PubKeyHash)
			}
		}

		fmt.Println()

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