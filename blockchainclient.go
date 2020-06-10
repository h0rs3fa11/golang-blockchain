package main

import (
	"fmt"
	"net/rpc"
)

var serverAddress string = "127.0.0.1"

// func main() {
// 	client, err := rpc.DialHTTP("tcp", serverAddress + ":8332")
// 	if err != nil {
// 		fmt.Printf("dialing error:%s", err)
// 	}

// 	args := ""
// 	var reply string
// 	err = client.Call("Rpc.Help", args, &reply)
// 	if err != nil {
// 		fmt.Printf("Call RPC API error:%s", err)
// 	} else {
// 		fmt.Println(reply)
// 	}
// }


func clientConnect() {
	client, err := rpc.DialHTTP("tcp", serverAddress + ":8332")
	if err != nil {
		fmt.Printf("dialing error:%s", err)
	}

	args := ""
	var reply string
	err = client.Call("Rpc.Help", args, &reply)
	if err != nil {
		fmt.Printf("Call RPC API error:%s", err)
	} else {
		fmt.Println(reply)
	}
}