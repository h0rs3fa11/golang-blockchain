package rpc

import (
	"net"
	"fmt"
	"net/rpc"
	"net/http"
)

type Rpc int

func StartRpcServer() {
	que := new(Rpc)
	rpc.Register(que)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", "127.0.0.1:8332")
	if err != nil {
		fmt.Printf("Listen error:%s", err)
	}

	go http.Serve(l, nil)
}