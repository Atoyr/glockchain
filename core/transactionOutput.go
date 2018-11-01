package core

import (
	"crypto/sha256"

	"github.com/atoyr/glockchain/util"
)

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

func (txo *TXOutput) Hash() []byte {
	var b []byte
	b = append(b, util.Int2bytes(txo.Value, 8)...)
	b = append(b, txo.PubKeyHash...)
	hash := sha256.Sum256(b)
	return hash[:]
}

func (out *TXOutput) Lock(address Address) {
	pubKeyHash := util.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}
func NewTXOutput(value int, address Address) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock(address)
	return txo
}
