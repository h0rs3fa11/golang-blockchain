package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
	"math/big"
	"github.com/boltdb/bolt"
	"time"
)

type CLI struct {
	Chain *Blockchain
}

func (cli *CLI) printUsage()  {

	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("\taddblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("\tprintchain - print all the blocks of the blockchain")

}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2{
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	var blockchainIterator ChainIterator

	blockchainIterator = *cli.Chain.Iterator()

	var hashInt big.Int

	for {

		err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
			// 获取表
			b := tx.Bucket([]byte(blocksBucket))
			// 通过Hash获取区块字节数组
			blockBytes := b.Get(blockchainIterator.CurrentHash)

			block := DeserializeBlock(blockBytes)

			fmt.Printf("Data：%s \n",string(block.Data))
			fmt.Printf("PrevBlockHash：%x \n",block.PrevBlockHash)
			fmt.Printf("Timestamp：%s \n",time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM") )
			fmt.Printf("Hash：%x \n",block.Hash)
			fmt.Printf("Nonce：%d \n",block.Nonce)

			fmt.Println();

			return nil
		})

		if err != nil {
			log.Panic(err)
		}

		// 获取下一个迭代器
		blockchainIterator = *blockchainIterator.Next()

		// 将迭代器中的hash存储到hashInt
		hashInt.SetBytes(blockchainIterator.CurrentHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}
	}
}

func (cli *CLI) addBlock(data string) {
	cli.Chain.AddBlock(data)
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")
	//fmt.Println("CLI Run \n")
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}