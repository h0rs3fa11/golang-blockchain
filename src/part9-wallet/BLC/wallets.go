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

func newWallets() (*wallets, error) {
	if _, err := os.Stat(walletFie); os.IsNotExist(err) {
		wallets := &Walets{}
		wallets.WalletMap = make(map[tring]*Wallet)
		return wallets,err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256( ))
	decoder := gob.NewDeoder(bytes.NewReader(fileContet))
	err = decoder.Deode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets, nil
}

func (w *Wallets) createNewWallet() {
	wallet := newWallets()
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
