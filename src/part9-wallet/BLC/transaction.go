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
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) isCoinbase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1 && len(tx.Vin[0].Txid) == 0
}

func createCoinbaseTx(to string, data string, params Chainparams) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	//create transaction input
	txIn := TXInput{[]byte{}, -1, data}

	//create transaction output
	txOut := TXOutput{params.Subsidy, to}

	tx := Transaction{nil, []TXInput{txIn}, []TXOutput{txOut}}
	tx.SetID()

	return &tx
}

func createTransaction(from string, to string, value int, bc *Blockchain) *Transaction {

	var feeOutput TXOutput
	var inputs []TXInput
	needValue := value + bc.Params.Fee
	//查找账户可用的UTXO
	findAmount, unspentOut := bc.findUnspentOutput(from, needValue)

	for txid, outs := range unspentOut {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}

			inputs = append(inputs, input)
		}
	}
	//新建交易结构
	tx := Transaction{nil, inputs, []TXOutput{TXOutput{value, to}}}

	//处理找零
	change := findAmount - value
	if change > bc.Params.Fee {
		changeOutput := TXOutput{change - bc.Params.Fee, from}
		feeOutput = TXOutput{bc.Params.Fee, miner}
		tx.Vout = append(tx.Vout, changeOutput)
		tx.Vout = append(tx.Vout, feeOutput)

	} else if change == bc.Params.Fee {
		feeOutput = TXOutput{bc.Params.Fee, miner}
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

func (in *TXInput) CanUnlockOutput(address string) bool {
	return in.ScriptSig == address
}

func (out *TXOutput) CanUnlock(address string) bool {
	return out.ScriptPubKey == address
}
