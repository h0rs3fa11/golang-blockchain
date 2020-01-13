package BLC

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

type CLI struct {
	Chain *Blockchain
}

func (cli *CLI) printUsage() {

	fmt.Println("Usage:")
	fmt.Println("\tgetbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("\tcreateblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("\tprintchain - Print all the blocks of the blockchain:")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")

}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
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

			fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
			fmt.Printf("Timestamp：%s \n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
			fmt.Printf("Hash：%x \n", block.Hash)
			fmt.Printf("Nonce：%d \n", block.Nonce)
			for tx := range block.Transaction {
				fmt.Printf("Transaction id: %x\n", block.Transaction[tx].ID)
			}

			fmt.Println()

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
			break
		}
	}
}

func (cli *CLI) addBlock(data string) {
	tx := createTransaction("system", "mqx", 5, cli.Chain)
	cli.Chain.AddBlock([]*Transaction{tx})
}

func (cli *CLI) getBlock(hash string) {

	hashByte, err := hex.DecodeString(hash)

	if err != nil {
		log.Panic(err)
	}

	err = cli.Chain.Database.View(func(tx *bolt.Tx) error {
		// 获取表
		b := tx.Bucket([]byte(blocksBucket))
		// 通过Hash获取区块字节数组
		blockBytes := b.Get(hashByte)

		block := DeserializeBlock(blockBytes)

		fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
		fmt.Printf("Timestamp：%s \n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x \n", block.Hash)
		fmt.Printf("Nonce：%d \n", block.Nonce)
		for _, tx := range block.Transaction {
			fmt.Printf("Transaction id: %x\n", tx.ID)

			for _, in := range tx.Vin {
				fmt.Printf("Transaction Input:\ntxid: %x\nout index: %d\nscriptSig:%s\n", in.Txid, in.Vout, in.ScriptSig)
			}
			fmt.Println()
			for _, out := range tx.Vout {
				fmt.Printf("Transaction Output:\nvalue: %d\nscriptPubKey: %s\n", out.Value, out.ScriptPubKey)
			}
		}

		fmt.Println()

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

func (cli *CLI) sendMany(from string, to string, amount int) { //进行中
	fmt.Println("from:")
	fmt.Println(from)
	fmt.Println("to:")
	fmt.Println(to)
	fmt.Println("amount:")
	fmt.Println(amount)

	for index, f := range from {
		num, err = strconv.Atoi(amount[index])
		if err != nil {
			log.Panic(err)
		}
		tx := createTransaction(f, to[index], num, cli.Chain)
		cli.Chain.AddBlock([]*Transaction{tx})
	}
}

func (cli *CLI) getBalance(address string) {
	//遍历区块？
	fmt.Println("Not finish yet")
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBlockCmd := flag.NewFlagSet("getblock", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendManyCmd := flag.NewFlagSet("sendmany", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")
	getBlockData := getBlockCmd.String("hash", "", "Block hash")
	getBalanceData := getBalanceCmd.String("address", "", "address")
	sendFrom := sendManyCmd.String("from", "", "from address")
	sendTo := sendManyCmd.String("to", "", "to address")
	sendAmount := sendManyCmd.String("amount", "", "the amount you want to send")

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
	case "getblock":
		err := getBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "sendmany":
		err := sendManyCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		cli.printUsage()
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

	if getBlockCmd.Parsed() {
		cli.getBlock(*getBlockData)
	}
	if getbalanceCmd.Parsed() {
		cli.getBalance(*getbalanceData)
	}

	if sendManyCmd.Parsed() {
		fromAddress := JSONToArray(*sendFrom)
		toAddress := JSONToArray(*sendTo)
		sendAmounts := JSONToArray(*sendAmount)

		if len(fromAddress) == len(toAddress) && len(fromAddress) == len(sendAmounts) {
			cli.sendMany(fromAddress, toAddress, sendAmounts)
		} else {
			fmt.Println("Arguments error")
		}
	}
}
