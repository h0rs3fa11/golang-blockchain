package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type ChainIterator struct {
	CurrentHash []byte
	DB *bolt.DB
}

func (chainiterator *ChainIterator) Next() *Block {
	var block *Block
	err := chainiterator.DB.View(func(tx *bolt.Tx) error {

		// 获取表
		b := tx.Bucket([]byte(blocksBucket))

		// 通过当前的Hash获取Block
		currentBlockBytes := b.Get(chainiterator.CurrentHash)

		// 反序列化
		block = DeserializeBlock(currentBlockBytes)

		chainiterator.CurrentHash = block.PrevBlockHash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block
}