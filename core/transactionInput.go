package core

import (
	"github.com/atoyr/glockchain/util"
)

type TXInput struct {
	PrevTXHash  Hash
	PrevTXIndex int
	Signature   []byte
}

func NewTXInput(prevTX Transaction, prevTXIndex int) *TXInput {
	var txi TXInput
	return &txi
}

func (txd *TXInput) Lock(address Address) {
	pubKeyHash := util.Base58Decode(address.Bytes())
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
}
