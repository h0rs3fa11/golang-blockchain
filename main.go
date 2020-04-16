package main

import (
	"blockchain/core"
	//"blockchain/client"
	"blockchain/rpcapi"
)

func main() {
    blockchain := core.NewBlockChain()

	rpcapi.StartRpc(blockchain)
}