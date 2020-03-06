package BLC

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
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
	fmt.Println("\tsendmany -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
<<<<<<< HEAD

=======
	fmt.Println("\tlistaddress list all address from wallet")
<<<<<<< HEAD
	fmt.Println("\tcreateaddress create a new address")
	fmt.Println("\tlistaddress list all adress")
	fmt.Println("\tcleanblockchain clean the blockchain database file and wallet file")
=======
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
>>>>>>> 39d6bc9b07579ef92c09d2d4ce3ede171e6048e2
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

			fmt.Printf("Height:%d\n", block.Height)
			fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
			fmt.Printf("Timestamp：%s \n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
			fmt.Printf("Hash：%x \n", block.Hash)
			fmt.Printf("Nonce：%d \n", block.Nonce)
			for _, tx := range block.Transaction {
				fmt.Printf("Transaction id: %x\n", tx.ID)
				//遍历vin
				fmt.Println("----------transaction input----------")
				for _, txin := range tx.Vin {
					fmt.Printf("Vin transaction ID: %x\n", txin.Txid)
					fmt.Printf("Vin Vout: %d\n", txin.Vout)
					//fmt.Printf("Script Sig: %s\n", txin.ScriptSig)
				}
				fmt.Println("----------transaction output----------")
				//遍历vout
<<<<<<< HEAD
				for _, txout := range tx.Vout {
					fmt.Printf("Vout value: %d\n", txout.Value)
					//fmt.Printf("Vout ScriptPubKey: %s\n", txout.ScriptPubKey)
=======
				fmt.Println("Vouts:")
				for _, txout := range tx.Vout {
					fmt.Println(txout.Value)
					fmt.Printf("%s",GetAddressFromPubkey(txout.PubKeyHash))
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
				}
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

<<<<<<< HEAD
func (cli *CLI) addBlock(data string) {
	tx := createTransaction("system", "mqx", 5, cli.Chain)
	cli.Chain.AddBlock([]*Transaction{tx})
}
=======
// func (cli *CLI) addBlock(data string, address string) {
// 	tx := createTransaction("system", "mqx", 5, cli.Chain)
// 	cli.Chain.AddBlock([]*Transaction{tx}, address)
// }
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c

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

<<<<<<< HEAD
func (cli *CLI) sendMany(from string, to string, amount int) {
=======
func (cli *CLI) sendMany(from string, to string, amount int, miner string) {
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	fmt.Printf("from:%s\n", from)
	fmt.Printf("to:%s\n", to)
	fmt.Printf("amount:%d\n", amount)

<<<<<<< HEAD
	tx := createTransaction(from, to, amount, cli.Chain)
	cli.Chain.AddBlock([]*Transaction{tx})
=======
	tx := createTransaction(from, to, amount, cli.Chain, "")
	cli.Chain.AddBlock([]*Transaction{tx}, miner)
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
}

func (cli *CLI) getBalance(address string) int {
	txs := cli.Chain.findUnspentTX(address)
<<<<<<< HEAD

=======
	pubKey := getPublickey(address)
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	balance := 0

	for _, tx := range txs {
		for _, out := range tx.Vout {
<<<<<<< HEAD
			if out.CanUnlock(address) {
=======
			if out.IsLockWithKey(HashPubKey(pubKey)) {
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
				balance += out.Value
			}
		}
	}

	return balance
}

func (cli *CLI) listaddress() {
<<<<<<< HEAD
	wallet,_ := NewWallets()
=======
	wallets,_ := newWallets()
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	for address, _ := range wallets.WalletsMap {
		fmt.Println(address)
	}
}

<<<<<<< HEAD
=======
func (cli *CLI) createAddress() {
	wallets,_ := newWallets()
	wallets.createNewWallet()
}

<<<<<<< HEAD


=======
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
>>>>>>> 39d6bc9b07579ef92c09d2d4ce3ede171e6048e2
func (cli *CLI) Run() {
	cli.validateArgs()
	//getblock 2
	//gettransaction
<<<<<<< HEAD
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
=======
	//addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addresslistCmd := flag.NewFlagSet("listaddress", flag.ExitOnError)
	getBlockCmd := flag.NewFlagSet("getblock", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendManyCmd := flag.NewFlagSet("sendmany", flag.ExitOnError)
<<<<<<< HEAD

	addBlockData := addBlockCmd.String("data", "", "Block data")
=======
	createAddrCmd := flag.NewFlagSet("createaddress", flag.ExitOnError)
	//cleandbCmd := flag.NewFlagSet("cleanblockchain", flag.ExitOnError)

	//addBlockData := addBlockCmd.String("data", "", "Block data")
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	getBlockData := getBlockCmd.String("hash", "", "Block hash")
	getBalanceData := getBalanceCmd.String("address", "", "check balances of address")
	sendFrom := sendManyCmd.String("from", "", "from address")
	sendTo := sendManyCmd.String("to", "", "to address")
	sendAmount := sendManyCmd.Int("amount", 0, "the amount you want to send")
<<<<<<< HEAD

	//fmt.Println("CLI Run \n")
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
=======
	sendMiner := sendManyCmd.String("miner", "", "the miner of this block")

	//sendMemo := sendManyCmd.String("memo", "", "")
	//fmt.Println("CLI Run \n")
	switch os.Args[1] {
	// case "addblock":
	// 	err := addBlockCmd.Parse(os.Args[2:])
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c

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
<<<<<<< HEAD

=======
	case "createaddress":
		err := createAddrCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
<<<<<<< HEAD
	case "cleanblockchain":
		err := cleandbCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
=======
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
>>>>>>> 39d6bc9b07579ef92c09d2d4ce3ede171e6048e2
	default:
		cli.printUsage()
		os.Exit(1)
	}

<<<<<<< HEAD
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}
=======
	// if addBlockCmd.Parsed() {
	// 	if *addBlockData == "" {
	// 		addBlockCmd.Usage()
	// 		os.Exit(1)
	// 	}
	// 	cli.addBlock(*addBlockData)
	// }
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c

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
<<<<<<< HEAD
		cli.sendMany(*sendFrom, *sendTo, *sendAmount)
=======
		var miner string
		if *sendMiner == "" {
			miner = cli.Chain.Params.Miner
		} else {
			miner = *sendMiner 
		}
		cli.sendMany(*sendFrom, *sendTo, *sendAmount, miner)
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	}

	if addresslistCmd.Parsed() {
		cli.listaddress()
	}
<<<<<<< HEAD
=======

	if createAddrCmd.Parsed() {
		cli.createAddress()
	}
<<<<<<< HEAD

	if cleandbCmd.Parsed() {
		cli.cleanBlockchain()
	}
=======
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
>>>>>>> 39d6bc9b07579ef92c09d2d4ce3ede171e6048e2
}
