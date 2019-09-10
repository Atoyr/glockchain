package glockChain

import (
	"bytes"
	"crypto/sha256"

	"github.com/atoyr/glockchain/util"
)

// TXOutput Transaction output
type TXOutput struct {
	Value      int    `json:"value"`
	PubKeyHash []byte `json:"pub_key_hash"`
}

// NewTXOutput TXOutput constructor
func NewTXOutput(value int, address []byte) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock(address)
	return txo
}

// Hash Hash transaction output
func (txo *TXOutput) Hash() []byte {
	var b []byte
	b = append(b, util.Int2bytes(txo.Value, 8)...)
	b = append(b, txo.PubKeyHash...)
	hash := sha256.Sum256(b)
	return hash[:]
}

// Lock TransactionOutput Locked
func (txo *TXOutput) Lock(address []byte) {
	pubKeyHash := AddressToPubKeyHash(address)
	txo.PubKeyHash = pubKeyHash
}

// IsLockedWithKey Is Locked with args keyhash?
func (txo *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(txo.PubKeyHash, pubKeyHash) == 0
}
