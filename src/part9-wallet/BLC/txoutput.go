package BLC

import (
	"bytes"
)

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockWithKey(pubKeyHash []byte) {
	return bytes.Compare(out.PubKeyHash, pubKeyHash)
}
