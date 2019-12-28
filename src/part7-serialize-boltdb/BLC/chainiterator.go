package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type ChainIterator struct {
	CurrentHash []byte
	DB *bolt.DB
}

func (blockchain *Blockchain) Iterator() *ChainIterator {
	return &ChainIterator{blockchain.Tip, blockchain.Database}
}

func (chainiterator *ChainIterator) Next() *ChainIterator {
	var nextHash []byte

	err := chainiterator.DB.View(func(tx *bolt.Tx) error {

		// 获取表
		b := tx.Bucket([]byte(blocksBucket))

		// 通过当前的Hash获取Block
		currentBlockBytes := b.Get(chainiterator.CurrentHash)

		// 反序列化
		currentBlock := DeserializeBlock(currentBlockBytes)

		nextHash = currentBlock.PrevBlockHash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &ChainIterator{nextHash, chainiterator.DB}
}