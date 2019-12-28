package main

import (
	"golang-blockchain/src/part8-transaction/BLC"
	//"fmt"
	//"math/big"
)

func main() {
	blockchain := BLC.NewBlockChain()

	cli := BLC.CLI{blockchain}

	cli.Run()
}
