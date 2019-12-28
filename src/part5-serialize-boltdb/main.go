package main

import (
	"golang-blockchain/src/part5-serialize-boltdb/BLC"
	"fmt"
)

func main() {
	blockchain := BLC.NewBlockChain()

	blockchain.AddBlock("Send 20 BTC To M");
	blockchain.AddBlock("Send 10 BTC To M");
	blockchain.AddBlock("Send 5 BTC To M");

	//fmt.Println(blockchain)
	fmt.Printf("blockchain tip block: %x\n", blockchain.Tip)
	blockchain.FindBlock(blockchain.Tip)
}
