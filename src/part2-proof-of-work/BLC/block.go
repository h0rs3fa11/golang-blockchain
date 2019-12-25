package BLC

import (
  "bytes"
  "time"
  "strconv"
  //"fmt"
  "crypto/sha256"
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

func (block *Block)SetHash() {
  //int64 => string
  timeString := strconv.FormatInt(block.Timestamp, 2)
  //string => byte
  timestamp := []byte(timeString)

  headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})

  hash := sha256.Sum256(headers)

  block.Hash = hash[:]
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
