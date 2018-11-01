package core

import (
	"crypto/sha256"

	"github.com/atoyr/glockchain/util"
)

type TXInput struct {
	PrevTXHash  []byte
	PrevTXIndex int
	Signature   []byte
	PubKey      []byte
}

func (txi *TXInput) Hash() []byte {
	var b []byte
	b = append(b, txi.PrevTXHash...)
	b = append(b, util.Interface2bytes(txi.PrevTXHash)...)
	b = append(b, txi.Signature...)
	b = append(b, txi.PubKey...)
	hash := sha256.Sum256(b)
	return hash[:]
}

func NewTXInput(prevTX Transaction, prevTXIndex int) *TXInput {
	var txi TXInput
	return &txi
}

func (txd *TXInput) Lock(address []byte) {
	pubKeyHash := util.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
}
