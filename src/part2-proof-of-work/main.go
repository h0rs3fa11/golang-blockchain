package main

import (
	"golang-blockchain/src/part1-Basic-Prototype/BLC"
	"fmt"
	"time"
)

func main() {
	blockchain := BLC.NewBlockChain()

	blockchain.AddBlock("Send 20 BTC To M");
	blockchain.AddBlock("Send 10 BTC To M");
	blockchain.AddBlock("Send 5 BTC To M");

	for _, block := range blockchain.Blocks {
		fmt.Printf("Data:%s \n", string(block.Data))
		fmt.Printf("Timestamp:%s \n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("PrevBlockHash:%x \n", block.PrevBlockHash)
		fmt.Printf("Hash:%x \n", block.Hash)
		fmt.Println("\n")
	}
	//fmt.Println(block)
}
