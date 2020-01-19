package main

import (
	"golang-blockchain/src/part10-signature/BLC"
	//"fmt"
	//"math/big"
)

func main() {
	blockchain := BLC.NewBlockChain()

	cli := BLC.CLI{blockchain}

	cli.Run()
}
