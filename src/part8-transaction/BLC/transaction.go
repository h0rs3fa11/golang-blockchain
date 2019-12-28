package BLC

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/sha256"
	"log"
	"os"
	"encoding/hex"
)

const subsidy = 10
const fee = 1
//(temp) miner
const miner = "temp-miner"

type Transaction struct {
	ID []byte
	Vin []TXInput
	Vout []TXOutput
}

func (tx *Transaction) isCoinbase() bool {
	return len(tx.Vin) == 1	&& tx.Vin[0].Vout == -1 && len(tx.Vin[0].Txid) == 0
}

func createCoinbaseTx(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	//create transaction input
	txIn := TXInput{[]byte{}, -1, data}

	//create transaction output
	txOut := TXOutput{subsidy, to}

	tx := Transaction{nil, []TXInput{txIn}, []TXOutput{txOut}}
	tx.SetID()

	return &tx
}

func createTransaction(from string, to string, value int, bc *Blockchain) *Transaction {

	var feeOutput TXOutput
	var inputs []TXInput
	//查找账户可用的UTXO
	findAmount, unspentOut := bc.findUnspentOutput(from, value)

	for txid, outs := range unspentOut {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _,out := range outs {
			input := TXInput{txID, out, from}

			inputs = append(inputs, input)
		}
	}
	//新建交易结构
	tx := Transaction{nil, inputs, []TXOutput{TXOutput{value, to}}}

	//处理找零
	change := findAmount - value
	if change > fee {
		changeOutput := TXOutput{change - fee, from}
		feeOutput = TXOutput{fee, miner}
		tx.Vout = append(tx.Vout, changeOutput)
		tx.Vout = append(tx.Vout, feeOutput)

	} else if change == fee {
		feeOutput = TXOutput{fee, miner}
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

type TXInput struct {
	Txid []byte
	Vout int
	ScriptSig string //name
}

func (in *TXInput) CanUnlockOutput(address string) bool {
	return in.ScriptSig == address
}

type TXOutput struct {
	Value int
	ScriptPubKey string //name
}

func (out *TXOutput) CanUnlock(address string) bool {
	return out.ScriptPubKey == address
}