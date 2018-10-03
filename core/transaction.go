package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"os"
)

// Transaction Data
type Transaction struct {
	Version   int
	BlockHash Hash
	Input     []*TXData
	Output    []*TXData
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func DeserializeTransaction(data []byte) Transaction {
	var tx Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}
	return tx
}

func NewTransaction(inputs []*TXData, to Address, value int) *Transaction {
	sumValue := 0
	from := inputs[0].Address
	for _, txd := range inputs {
		sumValue += txd.Value
	}
	diffValue := sumValue - value
	if diffValue < 0 {
		return nil
	}
	var tx Transaction
	tx.Version = int(version)
	tx.BlockHash = BytesToHash([]byte{})
	tx.Input = inputs

	outputs := make([]*TXData, 2)
	outputs = append(outputs, []*TXData{NewTXData(0, to, value)}...)
	if diffValue > 0 {
		outputs = append(outputs, []*TXData{NewTXData(1, from, diffValue)}...)
	}
	tx.Output = outputs
	return &tx
}
