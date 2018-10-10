package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
	"os"
)

// Transaction Data
type Transaction struct {
	Version   int
	BlockHash Hash
	R         big.Int
	S         big.Int
	Input     []*TXInput
	Output    []*TXOutput
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

func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey) {

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

func NewTransaction(prevOutput []*TXOutput, to Address, value int) *Transaction {
	sumValue := 0
	for _, txo := range prevOutput {
		sumValue += txo.Value
	}
	diffValue := sumValue - value
	if diffValue < 0 {
		return nil
	}
	var tx Transaction
	tx.Version = int(Version)
	tx.BlockHash = BytesToHash([]byte{})
	// tx.Input = inputs

	outputs := make([]*TXOutput, 2)
	outputs = append(outputs, []*TXOutput{NewTXOutput(value, to)}...)
	//if diffValue > 0 {
	//outputs = append(outputs, []*TXOutput{NewTXOutput(from, diffValue)}...)
	//}
	tx.Output = outputs
	return &tx
}
