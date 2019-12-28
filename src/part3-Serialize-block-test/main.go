package main

import(
	"golang-blockchain/src/part3-Serialize-block/BLC"
	"fmt"
)

func main() {
	block := BLC.NewBlock("Test", []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	fmt.Printf("%d\n", block.Nonce)
	fmt.Printf("%x\n", block.Hash)

	bytes := block.Serialize()

	fmt.Println(bytes)

	block =BLC.DeserializeBlock(bytes)

	fmt.Printf("%d\n", block.Nonce)
	fmt.Printf("%x\n", block.Hash)
}
