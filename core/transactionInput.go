package core

import (
	"github.com/atoyr/glockchain/util"
)

type TXInput struct {
	PrevTXHash  []byte
	PrevTXIndex int
	Signature   []byte
	PubKey      []byte
}

func NewTXInput(prevTX Transaction, prevTXIndex int) *TXInput {
	var txi TXInput
	return &txi
}

func (txd *TXInput) Lock(address []byte) {
	pubKeyHash := util.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
}
