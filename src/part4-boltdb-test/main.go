package main

import(
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

//database name
const dbFile = "blockchain.db"

//bucket
const blocksBucket = "blocks"

func main() {

	//open or create a database
	db,err := bolt.Open(dbFile, 0600, nil)
	if err !=nil {
		log.Panic(err)
	}
	defer db.Close()

	//update data
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

			err = b.Put([]byte("mqx"), []byte("123456"))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("ljt"), []byte("12345678"))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	});

	//get value from database
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		valueByte := b.Get([]byte("mqx"))

		fmt.Printf("%s", valueByte)

		return nil
	});
}
