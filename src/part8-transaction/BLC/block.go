package BLC

import (
  "bytes"
  "time"
  "encoding/gob"
  "log"
)

var blockNumber = 0

type Block struct {
        //block index
        Index int
        // timestamp
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
	block := &Block{blockNumber - 1, time.Now().Unix(), prevBlockHash, []byte(data), []byte{}, 0}
        //create pow object
        pow := NewProofOfWork(block)
        nonce,hash := pow.Run()
        //pow.run,create a block
        block.Hash = hash
        block.Nonce = nonce
        blockNumber++
	return block;
}

func NewGenesisBlock() *Block {
        blockNumber++
	return NewBlock("Genesis Block", []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}
