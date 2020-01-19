package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

var blockNumber = 0

type Block struct {
	Height int
	// timestamp
	Timestamp int64
	//previous hash
	PrevBlockHash []byte
	//transaction data
	Transaction []*Transaction
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

func (b *Block) hashTransaction() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transaction {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
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

func NewBlock(transaction []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{height, time.Now().Unix(), prevBlockHash, transaction /*transaction*/, []byte{}, 0}
	//create pow object
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	//pow.run,create a block
	block.Hash = hash
	block.Nonce = nonce

	return block
}

func NewGenesisBlock(bc *Blockchain) *Block {
	// Initialize Genesis Transaction
	//需要一个初始账户，并把coinbase交易发送到该账户
	//调用createTransaction（coinbaseTransaction）
	coinbaseTx := createCoinbaseTx("system", "", bc.Params)
	return NewBlock([]*Transaction{coinbaseTx}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1)
}
