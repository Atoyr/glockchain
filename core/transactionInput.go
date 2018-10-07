package core

import (
	"github.com/atoyr/glockchain/util"
)

type TXInput struct {
	PrevTXHash  Hash
	PrevTXIndex int
	Address     Address
	Signature   []byte
	PubKey      []byte
}

func NewTXData(txindex int, address Address, value int) *TXInput {
	var txi TXInput
	txi.Address = address
	return &txi
}

func (txd *TXInput) Lock(address Address) {
	pubKeyHash := util.Base58Decode(address.Bytes())
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	txd.PubKey = pubKeyHash
}
