package rpcapi

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"

	"blockchain/core"
)


type CLI struct {
	rpcclient *rpc.Client
	apiName string
	reply string
}

func (cli *CLI)printUsage() {
	fmt.Println(`Usage:
		getbalance -address ADDRESS - Get balance of ADDRESS 
		createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS
		printchain - Print all the blocks of the blockchain
		sendmany -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO
		listaddress list all address from wallet
		createaddress create a new address
	`)
}

func InitClient(address string) (*CLI, error){
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}

	cli := CLI{
		apiName: "",
		rpcclient: client,
	}

	return &cli, nil
}

func (cli *CLI)validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI)printChain() {
	cli.apiName = "printChain"
	if core.DbExists() == false {
		fmt.Println("Database not exist")
		os.Exit(1)
	}

	var args = ""
	
	err := cli.rpcclient.Call("Rpc"+ cli.apiName, &args, &cli.reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply);
}

func (cli *CLI)getBlock(hash string) {
	cli.apiName = "getBlock"
	
	err := cli.rpcclient.Call("Rpc" + cli.apiName, hash, &cli.reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply);
}

func (cli *CLI) sendMany(from string, to string, amount int) {
	fmt.Printf("from:%s\n", from)
	fmt.Printf("to:%s\n", to)
	fmt.Printf("amount:%d\n", amount)

	tx, err := core.CreateTransaction(from, to, amount, Chain, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	Chain.AddBlock([]*core.Transaction{tx})
}

func (cli *CLI) getBalance(address string) int {
	txs := Chain.FindUnspentTX(address)
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
func (cli *CLI)ParseCmdAndCall() {
	cli.validateArgs()

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addresslistCmd := flag.NewFlagSet("listaddress", flag.ExitOnError)
	getBlockCmd := flag.NewFlagSet("getblock", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendManyCmd := flag.NewFlagSet("sendmany", flag.ExitOnError)
	createAddrCmd := flag.NewFlagSet("createaddress", flag.ExitOnError)

	//defines a flag with specified name
	getBlockData := getBlockCmd.String("hash", "", "Block hash")
	getBalanceData := getBalanceCmd.String("address", "", "check balances of address")
	sendFrom := sendManyCmd.String("from", "", "from address")
	sendTo := sendManyCmd.String("to", "", "to address")
	sendAmount := sendManyCmd.Int("amount", 0, "the amount you want to send")

	switch os.Args[1] {
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
			cli.apiName = "Help"
			cli.printUsage()
			os.Exit(1)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getBlockCmd.Parsed() {
		getBlock(*getBlockData)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			printUsage()
			os.Exit(1)
		}

		fmt.Printf("%s balance\n%d\n", *getBalanceData, getBalance(*getBalanceData))
	}

	if sendManyCmd.Parsed() {
		sendMany(*sendFrom, *sendTo, *sendAmount)
	}

	if addresslistCmd.Parsed() {
		listaddress()
	}

	if createAddrCmd.Parsed() {
		createAddress()
	}
}
