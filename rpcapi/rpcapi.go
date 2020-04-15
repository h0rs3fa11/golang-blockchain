package rpcapi

import (
	"blockchain/core"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

func (r *Rpc) Help(arg string, reply *string) error {
	*reply = fmt.Sprintln(`Usage:
		getbalance -address ADDRESS - Get balance of ADDRESS 
		createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS
		printchain - Print all the blocks of the blockchain
		sendmany -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO
		listaddress list all address from wallet
		createaddress create a new address
	`)
	return nil
}

func (r *Rpc) printChain(arg string, reply *string) error {
	r.bc.PrintChain()
	*reply = fmt.Sprintln("success")
	return nil
}

func (r *Rpc) getBlock(arg string, reply *string) error {

	hashByte, err := hex.DecodeString(hash)

	if err != nil {
		log.Panic(err)
	}

	err = r.bc.Database.View(func(tx *bolt.Tx) error {
		// 获取表
		b := tx.Bucket([]byte(core.BlocksBucket))
		// 通过Hash获取区块字节数组
		blockBytes := b.Get(hashByte)

		block := core.DeserializeBlock(blockBytes)

		*reply := fmt.Sprintf("PrevBlockHash：%x \nTimestamp：%s \nHash：%x \nNonce：%d \n", block.PrevBlockHash, time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"), block.Hash, block.Nonce)
		// TODO 怎么输出完整的reply
		for _, tx := range block.Transaction {
			fmt.Printf("Transaction id: %x\n", tx.ID)

			// for _, in := range tx.Vin {
			// 	fmt.Printf("Transaction Input:\ntxid: %x\nout index: %d\nscriptSig:%s\n", in.Txid, in.Vout, in.ScriptSig)
			// }
			// fmt.Println()
			// for _, out := range tx.Vout {
			// 	fmt.Printf("Transaction Output:\nvalue: %d\nscriptPubKey: %s\n", out.Value, out.ScriptPubKey)
			// }
		}

		fmt.Println()

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}