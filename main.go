package main

import (
	//"blockchain/core"
	//"blockchain/client"
	"blockchain/rpcapi"
)

func main() {
    //blockchain := core.NewBlockChain()
	
	//cli := client.CLI{blockchain}

    //cli.Run()
    
    rpcapi.StartRpc()
}