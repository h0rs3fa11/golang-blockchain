package BLC

import (
  "github.com/boltdb/bolt"
  "fmt"
  "log"
  "math/big"
  "encoding/hex"
)

const dbFile = "blockchain.db"

const blocksBucket = "blocks"

type Blockchain struct {
      Tip []byte
      Database *bolt.DB
}

//add new block
func (blockchain *Blockchain) AddBlock(tx []*Transaction) {
      //create new block
      block := NewBlock(tx, blockchain.Tip)

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

//Find Transaction contains UTXO of address
func (blockchain *Blockchain) findUnspentTX(address string) []Transaction {
        var unspentTXs []Transaction

        spentTXOs := make(map[string][]int) //存储已花费的UTXO
        bci := blockchain.Iterator()
        var hashInt big.Int

        for {
                err := blockchain.Database.View(func(tx *bolt.Tx) error {
                        // 获取表
                        b := tx.Bucket([]byte(blocksBucket))
                        // 通过Hash获取区块字节数组
                        blockBytes := b.Get(bci.CurrentHash)

                        block := DeserializeBlock(blockBytes)

                        //遍历区块交易
                        for _,transaction := range block.Transaction {
                                fmt.Printf("TransactionHash:%x\n", transaction.ID)
                                txID := hex.EncodeToString(transaction.ID)

                                //遍历交易输出
                                Outputs:
                                for outIdx, out := range transaction.Vout {
                                        //遍历spentTXOs
                                        if spentTXOs[txID] != nil {
                                                for _, spentOut := range spentTXOs[txID] {
                                                        if spentOut == outIdx{ //为什么比较金额，如果两笔金额一样的UTXO呢
                                                                continue Outputs
                                                        }
                                                }
                                        }

                                        if out.CanUnlock(address) {
                                                unspentTXs = append(unspentTXs, *transaction)
                                        }
                                }

                                if transaction.isCoinbase() == false {
                                        for _, in := range transaction.Vin {
                                                inTxID := hex.EncodeToString(in.Txid)
                                                spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
                                        }
                                }
                        }

                        fmt.Println();

                        return nil
                })

                if err != nil {
                        log.Panic(err)
                }

                bci = bci.Next()

                hashInt.SetBytes(bci.CurrentHash)

                if hashInt.Cmp(big.NewInt(0)) == 0 {
                        break;
                }
        }

        return unspentTXs
}

//Find UTXO
func (blockchain *Blockchain) findUnspentOutput(address string, amount int) (int, map[string][]int) {
        unspentOutputs := make(map[string][]int)
        unspentTx := blockchain.findUnspentTX(address)
        accumulated := 0

        Work:
        for _, tx := range unspentTx {
                txID := hex.EncodeToString(tx.ID)

                for outIdx, out := range tx.Vout {
                        if out.CanUnlock(address) && accumulated < amount {
                                accumulated += out.Value
                                unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

                                if accumulated >= amount {
                                        break Work
                                }
                        }
                }
        }
        return accumulated, unspentOutputs
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
