package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

//(temp) miner
const miner = "temp-miner"

type Transaction struct {
	ID   []byte
	Memo string
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) isCoinbase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1 && len(tx.Vin[0].Txid) == 0
}

func createCoinbaseTx(pubKeyHash []byte, memo string, params Chainparams) *Transaction {
	//create transaction input
	txIn := TXInput{[]byte{}, -1, nil, nil}

	//create transaction output
	txOut := TXOutput{params.Subsidy, pubKeyHash}

	tx := Transaction{nil, memo, []TXInput{txIn}, []TXOutput{txOut}}
	tx.SetID()

	return &tx
}

func createTransaction(from string, to string, value int, bc *Blockchain, memo string) *Transaction {

	var feeOutput TXOutput
	var inputs []TXInput
	needValue := value + bc.Params.Fee
	var frompubKey []byte
	//查找账户可用的UTXO
	findAmount, unspentOut := bc.findUnspentOutput(from, needValue)

	for txid, outs := range unspentOut {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		frompubKey = getPublickey(from)
		if frompubKey == nil {
			log.Panic("Can't find this address")
		}
		for _, out := range outs {
			input := TXInput{txID, out, nil, frompubKey}

			inputs = append(inputs, input)
		}
	}
	//目的地址的公钥，怎么从目的地址的地址推到公钥？如果本地没有这个密钥文件的话
	topubKey := getPublickey(to)
	if topubKey == nil {
		log.Panic("Can't find this address")
	}

	minerPubkey := getPublickey(bc.Params.Miner)
	if minerPubkey == nil {
		log.Panic("Can't find this address")
	}

	//新建交易结构
	tx := Transaction{nil, memo, inputs, []TXOutput{TXOutput{value, topubKey}}}

	//处理找零
	change := findAmount - value
	if change > bc.Params.Fee {
		changeOutput := TXOutput{change - bc.Params.Fee, frompubKey}
		feeOutput = TXOutput{bc.Params.Fee, minerPubkey}
		tx.Vout = append(tx.Vout, changeOutput)
		tx.Vout = append(tx.Vout, feeOutput)

	} else if change == bc.Params.Fee {
		feeOutput = TXOutput{bc.Params.Fee, minerPubkey}
		tx.Vout = append(tx.Vout, feeOutput)
	} else {
		fmt.Println("Transaction fee is not enough!\n")
		os.Exit(1)
	}

	tx.SetID()

	return &tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil {
		log.Panic(err)
	}
	// 将序列化以后的字节数组生成256hash
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}
