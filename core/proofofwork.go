package core

import (
	"crypto/sha256"
	"math"
	"math/big"

	"github.com/atoyr/glockchain/util"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

var maxnonce = math.MaxInt64

const targetbits = 4

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetbits))
	pow := ProofOfWork{b, target}

	return &pow
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < maxnonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	var b []byte
	b = append(b, pow.block.PreviousHash...)
	b = append(b, pow.block.HashTransactions()...)
	b = append(b, util.Int642bytes(pow.block.Timestamp)...)
	b = append(b, util.Int2bytes(nonce, 8)...)
	return b
}
