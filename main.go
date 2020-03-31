package main

import (
    "blockchain/core"
    "blockchain/client"
)

func main() {
    blockchain := core.NewBlockChain()
	
	cli := client.CLI{blockchain}

	cli.Run()
}