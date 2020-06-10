package main

import (
	"blockchain/core"
	//"blockchain/client"
	"blockchain/rpcapi"
)

func main() {
    blockchain := core.NewBlockChain()

	//TODO:这个程序主要是接受RPC请求再处理，另开发一个发送RPC请求的客户端
	rpcapi.StartRpc(blockchain)
}