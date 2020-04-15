package rpcapi

import (
	"blockchain/core"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

type Rpc struct {
	bc core.Blockchain
}

var RPC_ADDRESS = "127.0.0.1"
var RPC_PORT = "8332"

// TODO 接收命令行参数
// 如果没有参数就调用Help信息
// 有参数就解析参数，调用对应API

func StartRpc() {
	que := new(Rpc)
	rpc.Register(que)
	rpc.HandleHTTP()
	address := RPC_ADDRESS+":"+RPC_PORT
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Listen error:%s", err)
	}

	go http.Serve(l, nil)

	cli, err := InitClient(address)
	if err != nil {
		fmt.Printf("Listen error:%s", err)
	}
	cli.ParseCmdAndCall()

}