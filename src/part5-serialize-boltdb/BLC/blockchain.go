package BLC

import (
  "github.com/boltdb/bolt"
  "fmt"
  "log"
)

const dbFile = "blockchain.db"

const blocksBucket = "blocks"

type Blockchain struct {
  Tip []byte
  Database *bolt.DB
}

//add new block
func (blockchain *Blockchain) AddBlock(data string) {
  //create new block
  block := NewBlock(data, blockchain.Tip)
  //add block to database
  AddBlockToDatabase(block)
  blockchain.Tip = block.Hash
}

//Create a blockchain with genesis block
func NewBlockChain() *Blockchain {
      genesisBlock := NewGenesisBlock()
      db := AddBlockToDatabase(genesisBlock)
      tip := genesisBlock.Hash

  return &Blockchain{tip, db}
}

func AddBlockToDatabase(block *Block) *bolt.DB {
  db,err := bolt.Open(dbFile, 0600, nil)
  if err !=nil {
    log.Panic(err)
  }
  defer db.Close()

  err = db.Update(func(tx *bolt.Tx) error {
    // 判断这一张表是否存在于数据库中
    b := tx.Bucket([]byte(blocksBucket))
    if b == nil {
      fmt.Println("No existing blockchain found. Creating a new one...")
      // CreateBucket 创建表
      b, err := tx.CreateBucket([]byte(blocksBucket))
      if err != nil {
        log.Panic(err)
      }

      blockData := block.Serialize()

      err = b.Put(block.Hash, blockData)
      if err != nil {
        log.Panic(err)
      }
      fmt.Println("Create Genesis Block\n")

    }
    return nil
  });

  return db
}

//func (blockchain *Blockchain) FindBlock(hash []byte) {
//  db := blockchain.Database
//  //get value from database
//  err = db.View(func(tx *bolt.Tx) error {
//    b := tx.Bucket([]byte(blocksBucket))
//    valueByte := b.Get(hash)
//
//    fmt.Printf("%s", valueByte)
//
//    return nil
//  });
//}