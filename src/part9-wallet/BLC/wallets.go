package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"fmt"
)

const walletFile = "wallet.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

func newWallets() (*Wallets, error) {
<<<<<<< HEAD
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.WalletMap = make(map[string]*Wallet)
=======

	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
		return wallets,err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256( ))
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
<<<<<<< HEAD
	err = decoder.Deode(&wallets)
=======
	err = decoder.Decode(&wallets)
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	if err != nil {
		log.Panic(err)
	}

	return &wallets, nil
}

func (w *Wallets) createNewWallet() {
<<<<<<< HEAD
	wallet := newWallets()
=======
	wallet := newWallet()
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
	fmt.Printf("Address:%s\n", wallet.GetAddress())
	w.WalletsMap[string(wallet.GetAddress())] = wallet
	w.saveWallets()
}

func (w *Wallets) saveWallets() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err :=  encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
