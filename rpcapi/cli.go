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

func InitClient(address string) (*CLI, error){
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}

	cli := CLI{
		apiName: ".",
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

func (cli *CLI)printUsage() {
	cli.apiName += "Help"
	err := cli.rpcclient.Call("Rpc"+ cli.apiName, "", &cli.reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cli.reply);
}

func (cli *CLI)printChain() {
	cli.apiName += "PrintChain"
	if core.DbExists() == false {
		fmt.Println("Database not exist")
		os.Exit(1)
	}

	var args = ""
	
	err := cli.rpcclient.Call("Rpc"+ cli.apiName, &args, &cli.reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cli.reply);
}

func (cli *CLI)getBlock(hash string) {
	cli.apiName += "GetBlock"
	
	err := cli.rpcclient.Call("Rpc" + cli.apiName, &hash, &cli.reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cli.reply);
}

func (cli *CLI) sendMany(from string, to string, amount int) {
	
	fmt.Printf("from:%s\n", from)
	fmt.Printf("to:%s\n", to)
	fmt.Printf("amount:%d\n", amount)

	tx := NewTransfer(from, to, amount)

	cli.apiName += "Sendmany"
	
	err := cli.rpcclient.Call("Rpc" + cli.apiName, &tx, &cli.reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cli.reply);
}

func (cli *CLI) getBalance(address string) {
	cli.apiName += "Getbalance"
	err := cli.rpcclient.Call("Rpc" + cli.apiName, address, &cli.reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cli.reply);
}

func (cli *CLI) listaddress() {
	cli.apiName += "Listaddress"
	err := cli.rpcclient.Call("Rpc" + cli.apiName, "", &cli.reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cli.reply);
}

func (cli *CLI) createAddress() {
	cli.apiName += "Createaddress"
	err := cli.rpcclient.Call("Rpc" + cli.apiName, "", &cli.reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cli.reply);
}

//默认的coinbase，setCoinbase
func (cli *CLI)ParseCmdAndCall(address string) {
	cli, err := InitClient(address)
	if err != nil {
		fmt.Println(err)
	}

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
			cli.printUsage()
			os.Exit(1)
	}

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
		cli.getBalance(*getBalanceData)
		//fmt.Printf("%s balance\n%d\n", *getBalanceData, cli.getBalance(*getBalanceData))
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
