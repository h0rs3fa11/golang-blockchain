package BLC

import (
	"math/big"
	"bytes"
	"fmt"
	"math"
	"crypto/sha256"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 10

type  ProofOfWork struct {
	block *Block
	target *big.Int
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining block...\n")
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x ", hash)
		hashInt.SetBytes(hash[:])
		//target compare with hash
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n")

	return nonce,hash[:]
}

func (pow *ProofOfWork) prepareData(nonce int) ([]byte) {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.hashTransaction(),
			IntToHex(uint64(pow.block.Timestamp)),
			IntToHex(uint64(targetBits)),
			IntToHex(uint64(nonce)),
		},
		[]byte{},
	)
	return data
}

//validate proof of work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	//fmt.Printf("\nhashInt:%d, pow.target:%d", hashInt, pow.target)
	//fmt.Println(hash)

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))

	pow := &ProofOfWork{block, target}
	//fmt.Printf("%x\n", target)
	return pow
}