package glockchain

import (
	"crypto/sha256"

	"github.com/atoyr/glockchain/util"
)

// TXInput Transaction Input
type TXInput struct {
	PrevTXHash  []byte `json:"prev_tx_hash"`
	PrevTXIndex int    `json:"prev_tx_index"`
	Signature   []byte `json:"signature"`
	PubKey      []byte `json:"pub_key"`
}

// NewTXInput TXInput constructor
func NewTXInput(prevTX Transaction, prevTXIndex int) *TXInput {
	var txi TXInput
	return &txi
}

// Hash hash Transaction input
func (txi *TXInput) Hash() []byte {
	var b []byte
	b = append(b, txi.PrevTXHash...)
	b = append(b, util.Int2bytes(txi.PrevTXIndex, 4)...)
	b = append(b, txi.Signature...)
	b = append(b, txi.PubKey...)
	hash := sha256.Sum256(b)
	return hash[:]
}
