package BLC

import "bytes"

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

<<<<<<< HEAD
func (in *TXInput) UsesKey(pubKeyHash []byte) {
=======
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
<<<<<<< HEAD
=======

// func (in *TXInput) CanUnlockOutput(address string) bool {
// 	return in.ScriptSig == address
// }
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
