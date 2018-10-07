package core

import "github.com/atoyr/glockchain/util"

type TXOutput struct {
	Value     int
	PubSigKey []byte
}

func (out *TXOutput) Lock(address Address) {
	pubKeyHash := util.Base58Decode(address.Bytes())
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubSigKey = pubKeyHash
}
func NewTXOutput(value int, address Address) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock(address)
	return txo
}
