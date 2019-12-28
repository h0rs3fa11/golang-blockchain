package BLC

import (
  "github.com/boltdb/bolt"
  "fmt"
  "log"
  "time"
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

      err := blockchain.Database.Update(func(tx *bolt.Tx) error {
          b := tx.Bucket([]byte(blocksBucket))

          err := b.Put(block.Hash, block.Serialize())
          if err != nil {
              log.Panic(err)
          }

          err = b.Put([]byte("l"), block.Hash)
          if err != nil {
              log.Panic(err)
          }
          blockchain.Tip = block.Hash
          return nil
      })
      if err != nil {
          log.Panic(err)
      }
}

func (blockchain *Blockchain) FindBlock(hash []byte){
        //get value from database
        err := blockchain.Database.View(func(tx *bolt.Tx) error {
                b := tx.Bucket([]byte(blocksBucket))
                valueByte := b.Get(hash)

                block := DeserializeBlock(valueByte)

                fmt.Printf("Block Data:\nHash:%x\nTimestamp:%s\nPrevious Block Hash: %x\nData:%s\nNonce:%d\n", block.Hash, time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"), block.PrevBlockHash, block.Data, block.Nonce)
                return nil
        });

        if err != nil {
                log.Panic(err)
        }
}

//Create a blockchain with genesis block
func NewBlockChain() *Blockchain {
      var tip []byte

      db, err := bolt.Open(dbFile, 0600, nil)
      if err != nil {
          log.Panic(err)
      }

      err = db.Update(func(tx *bolt.Tx) error {
          b := tx.Bucket([]byte(blocksBucket))

          if b == nil {
              fmt.Println("No existing blockchain found. Create a new one ...")
              genesisBlock := NewGenesisBlock()

              b, err := tx.CreateBucket([]byte(blocksBucket))
              if err != nil {
                log.Panic(err)
              }

              err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
              if err != nil {
                  log.Panic(err)
              }

              err = b.Put([]byte("l"), genesisBlock.Hash)
              if err != nil {
                log.Panic(err)
              }

              tip = genesisBlock.Hash
          } else {
            tip = b.Get([]byte("l"))
          }

          return nil
      })
      if err != nil {
        log.Panic(err)
      }

      return &Blockchain{tip, db}
}
