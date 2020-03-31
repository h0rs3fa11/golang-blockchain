package core

import (
	"bytes"
	"log"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"

	"math/big"
	"crypto/elliptic"
	"fmt"
	"os"
)

//(temp) miner
const miner = "temp-miner"

type Transaction struct {
	ID   []byte
	Memo string
	Fee  int
	Vin  []*TXInput
	Vout []*TXOutput
}

func (tx *Transaction) isCoinbase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1 && len(tx.Vin[0].Txid) == 0
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.isCoinbase() {
		return
	}

	for _,vin := range tx.Vin {
		if prevTxs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTxs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		//当前交易输入中使用的vout
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.SetID()
		txCopy.Vin[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}
}

//拷贝一份交易用于签名
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, &TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, &TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, tx.Memo, 0, inputs, outputs}
	return txCopy
}

func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if tx.isCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTxs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTxs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.SetID()
		txCopy.Vin[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}

func createCoinbaseTx(pubKeyHash []byte, memo string, params Chainparams) *Transaction {
	//create transaction input
	txIn := &TXInput{[]byte{}, -1, nil, nil}

	//create transaction output
	txOut := &TXOutput{params.Subsidy, pubKeyHash}

	tx := Transaction{nil, memo, 0, []*TXInput{txIn}, []*TXOutput{txOut}}
	tx.SetID()

	return &tx
}

func CreateTransaction(from string, to string, value int, bc *Blockchain, memo string) *Transaction {

	//var feeOutput &TXOutput{}
	var inputs []*TXInput
	needValue := value + bc.Params.Fee

	wallets, _ := NewWallets()
	wallet := wallets.WalletsMap[from]

	//查找账户可用的UTXO
	findAmount, unspentOut := bc.findUnspentOutput(from, needValue)

	for txid, outs := range unspentOut {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := &TXInput{txID, out, nil, wallet.PublicKey}

			inputs = append(inputs, input)
		}
	}
	
	topubKey := GetPublickey(to)
	if topubKey == nil {
		log.Panic("Can't find this address")
	}

	minerPubkey := GetPublickey(getCoinbase())
	if minerPubkey == nil {
		log.Panic("Can't find this address")
	}

	//新建交易结构
	tx := Transaction{nil, memo, 0, inputs, []*TXOutput{&TXOutput{value, HashPubKey(topubKey)}}}

	//处理找零
	change := findAmount - value
	if change > bc.Params.Fee {
		changeOutput := &TXOutput{change - bc.Params.Fee, HashPubKey(wallet.PublicKey)}
		feeOutput := &TXOutput{bc.Params.Fee, minerPubkey}
		tx.Vout = append(tx.Vout, changeOutput)
		tx.Vout = append(tx.Vout, feeOutput)

	} else if change == bc.Params.Fee {
		feeOutput := &TXOutput{bc.Params.Fee, minerPubkey}
		tx.Vout = append(tx.Vout, feeOutput)
	} else {
		fmt.Println("Transaction fee is not enough!\n")
		os.Exit(1)
	}

	tx.SetID()


	bc.SignTransaction(&tx, wallet.PrivateKey)
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
