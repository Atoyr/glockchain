package core

import (
	"bytes"
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

func (txo *TXOutput) Lock(address []byte) {
	pubKeyHash := AddressToPubKeyHash(address)
	txo.PubKeyHash = pubKeyHash
}
func (txo *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(txo.PubKeyHash, pubKeyHash) == 0
}
func NewTXOutput(value int, address []byte) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock(address)
	return txo
}
