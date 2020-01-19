package main

import (
	"golang-blockchain/src/part9-wallet/BLC"
	//"fmt"
	//"math/big"
)

func main() {
	blockchain := BLC.NewBlockChain()

	cli := BLC.CLI{blockchain}

	cli.Run()
}
