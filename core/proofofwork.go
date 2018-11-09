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
	pow := ProofOfWork{b, 0}

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
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	var b []byte
	b = append(b, pow.block.PreviousHash...)
	b = append(b, pow.block.HashTransactions()...)
	b = append(b, util.Int642bytes(pow.block.Timestamp)...)
	b = append(b, util.Int642bytes(nonce)...)
	return b
}
