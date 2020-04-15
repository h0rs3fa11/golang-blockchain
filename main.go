package main

import (
    //"blockchain/core"
    //"blockchain/client"
    "blockchain/rpc"
)

func main() {
    //blockchain := core.NewBlockChain()
	
	//cli := client.CLI{blockchain}

    //cli.Run()
    
    rpc.StartRpcServer()
}