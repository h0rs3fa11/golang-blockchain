package main

import (
	"golang-blockchain/src/part9-wallet/BLC"
	//"fmt"
	//"math/big"
)

func main() {
<<<<<<< HEAD

=======
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	blockchain := BLC.NewBlockChain()

	cli := BLC.CLI{blockchain}

	cli.Run()
}
