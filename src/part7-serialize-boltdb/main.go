package main

import (
	"golang-blockchain/src/part7-serialize-boltdb/BLC"
	//"fmt"
	//"math/big"
)

func main() {
	blockchain := BLC.NewBlockChain()

	cli := BLC.CLI{blockchain}

	cli.Run()
}
