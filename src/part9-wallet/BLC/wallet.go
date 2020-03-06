package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func newWallet() *Wallet {
	privateKey, publicKey := newKeyPair()

	return &Wallet{privateKey, publicKey}
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, public
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hash := ripemd160.New()
	_, err := RIPEMD160Hash.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hash.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

<<<<<<< HEAD
//疑问 创世区块的第一个地址怎么产生 密钥怎么存储
=======
func getPublickey(address string) []byte {
	wallets,_ := newWallets()
	return wallets.WalletsMap[address].PublicKey
}
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
