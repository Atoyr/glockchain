package core

import (
	"github.com/atoyr/glockchain/util"
)

type TXData struct {
	TXHash    Hash
	TXIndex   int
	Address   Address
	Value     int
	Signature []byte
	PubKey    []byte
}

func NewTXData(txindex int, address Address, value int) *TXData {
	var txdata TXData
	txdata.TXIndex = txindex
	txdata.Address = address
	txdata.Value = value
	return &txdata
}

func (txd *TXData) Lock(address Address) {
	pubKeyHash := util.Base58Decode(address.Bytes())
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	txd.PubKey = pubKeyHash
}
