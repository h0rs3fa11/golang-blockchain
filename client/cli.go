package client

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"blockchain/core"
)

type CLI struct {
	Chain *core.Blockchain
}

func (cli *CLI) printUsage() {

	fmt.Println("Usage:")
	fmt.Println("\tgetbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("\tcreateblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("\tprintchain - Print all the blocks of the blockchain:")
	fmt.Println("\tsendmany -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
	fmt.Println("\tlistaddress list all address from wallet")
	fmt.Println("\tcreateaddress create a new address")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	if core.DbExists() == false {
		fmt.Println("Database not exist")
		os.Exit(1)
	}

	cli.Chain.PrintChain()
}

// func (cli *CLI) addBlock(data string, address string) {
// 	tx := createTransaction("system", "mqx", 5, cli.Chain)
// 	cli.Chain.AddBlock([]*Transaction{tx}, address)
// }

func (cli *CLI) getBlock(hash string) {

	hashByte, err := hex.DecodeString(hash)

	if err != nil {
		log.Panic(err)
	}

	err = cli.Chain.Database.View(func(tx *bolt.Tx) error {
		// 获取表
		b := tx.Bucket([]byte(core.BlocksBucket))
		// 通过Hash获取区块字节数组
		blockBytes := b.Get(hashByte)

		block := core.DeserializeBlock(blockBytes)

		fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
		fmt.Printf("Timestamp：%s \n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x \n", block.Hash)
		fmt.Printf("Nonce：%d \n", block.Nonce)
		for _, tx := range block.Transaction {
			fmt.Printf("Transaction id: %x\n", tx.ID)

			// for _, in := range tx.Vin {
			// 	fmt.Printf("Transaction Input:\ntxid: %x\nout index: %d\nscriptSig:%s\n", in.Txid, in.Vout, in.ScriptSig)
			// }
			// fmt.Println()
			// for _, out := range tx.Vout {
			// 	fmt.Printf("Transaction Output:\nvalue: %d\nscriptPubKey: %s\n", out.Value, out.ScriptPubKey)
			// }
		}

		fmt.Println()

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

func (cli *CLI) sendMany(from string, to string, amount int) {
	fmt.Printf("from:%s\n", from)
	fmt.Printf("to:%s\n", to)
	fmt.Printf("amount:%d\n", amount)

	tx, err := core.CreateTransaction(from, to, amount, cli.Chain, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	cli.Chain.AddBlock([]*core.Transaction{tx})
}

func (cli *CLI) getBalance(address string) int {
	txs := cli.Chain.FindUnspentTX(address)
	pubKey := core.GetPublickey(address)
	balance := 0

	for _, tx := range txs {
		for _, out := range tx.Vout {
			if out.IsLockWithKey(core.HashPubKey(pubKey)) {
				balance += out.Value
			}
		}
	}

	return balance
}

func (cli *CLI) listaddress() {
	wallets, _ := core.NewWallets()
	for address, _ := range wallets.WalletsMap {
		fmt.Println(address)
	}
}

func (cli *CLI) createAddress() {
	wallets, _ := core.NewWallets()
	wallets.CreateNewWallet()
}

//默认的coinbase，setCoinbase
func (cli *CLI) Run() {
	cli.validateArgs()
	//getblock 2
	//gettransaction
	//addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addresslistCmd := flag.NewFlagSet("listaddress", flag.ExitOnError)
	getBlockCmd := flag.NewFlagSet("getblock", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendManyCmd := flag.NewFlagSet("sendmany", flag.ExitOnError)
	createAddrCmd := flag.NewFlagSet("createaddress", flag.ExitOnError)

	//addBlockData := addBlockCmd.String("data", "", "Block data")
	getBlockData := getBlockCmd.String("hash", "", "Block hash")
	getBalanceData := getBalanceCmd.String("address", "", "check balances of address")
	sendFrom := sendManyCmd.String("from", "", "from address")
	sendTo := sendManyCmd.String("to", "", "to address")
	sendAmount := sendManyCmd.Int("amount", 0, "the amount you want to send")

	//sendMemo := sendManyCmd.String("memo", "", "")
	//fmt.Println("CLI Run \n")
	switch os.Args[1] {
	// case "addblock":
	// 	err := addBlockCmd.Parse(os.Args[2:])
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

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
	case "listaddress":
		err := addresslistCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createaddress":
		err := createAddrCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	// if addBlockCmd.Parsed() {
	// 	if *addBlockData == "" {
	// 		addBlockCmd.Usage()
	// 		os.Exit(1)
	// 	}
	// 	cli.addBlock(*addBlockData)
	// }

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getBlockCmd.Parsed() {
		cli.getBlock(*getBlockData)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			cli.printUsage()
			os.Exit(1)
		}

		fmt.Printf("%s balance\n%d\n", *getBalanceData, cli.getBalance(*getBalanceData))
	}

	if sendManyCmd.Parsed() {
		cli.sendMany(*sendFrom, *sendTo, *sendAmount)
	}

	if addresslistCmd.Parsed() {
		cli.listaddress()
	}

	if createAddrCmd.Parsed() {
		cli.createAddress()
	}
}
