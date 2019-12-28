package BLC

import (
	"bytes"
	"time"
	//"strconv"
	//"fmt"
	//"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	//timestamp
	Timestamp int64
	//previous hash
	PrevBlockHash []byte
	//transaction data
	Data []byte
	//block hash
	Hash []byte

	Nonce int
}

func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeBlock(blockBytes []byte) *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), prevBlockHash, []byte(data), []byte{}, 0}

	//create pow object
	pow := NewProofOfWork(block)
	nonce,hash := pow.Run()
	//pow.run,create a block
	//fmt.Println(pow.Validate())
	//fmt.Println("\n")
	block.Hash = hash
	block.Nonce = nonce

	return block;
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}