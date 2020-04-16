package core

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

func NewWallets() (*Wallets, error) {

	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
		return wallets,err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256( ))
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets, nil
}

func (w *Wallets) CreateNewWallet() string{
	wallet := newWallet()

	address := wallet.GetAddress()
	w.WalletsMap[string(address)] = wallet
	w.saveWallets()
	return string(address)
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

func getCoinbase() string {
	wallets,_ := NewWallets()

	for address, _ := range wallets.WalletsMap{
		return address
	}
	return ""
}
