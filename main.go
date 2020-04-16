package main

import (
	"blockchain/core"
	//"blockchain/client"
	"blockchain/rpcapi"
)

// TODO:RPC接口怎么跟Blockchain对应
func main() {
    blockchain := core.NewBlockChain()
	
	//cli := client.CLI{blockchain}

    //cli.Run()

    rpcapi.StartRpc(blockchain)
}