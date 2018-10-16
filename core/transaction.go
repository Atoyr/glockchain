package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"os"
)

// Transaction Data
type Transaction struct {
	Version   byte
	BlockHash []byte
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
	txCopy := tx.TrimmedCopy()
	for vindex, _ := range txCopy.Input {
		txCopy.Input[vindex].Signature = []byte{}
		txCopyHash := tx.Hash()
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopyHash)
		errorHandle(err)
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Input[vindex].Signature = signature
		txCopy.Input[vindex].Signature = nil
	}
}

func (tx *Transaction) Verify() bool {
	//txCopy := tx.TrimmedCopy()
	// curve := elliptic.P256()
	//for vindex, vin := range txCopy.Input {
	//}
	return true
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Input {
		inputs = append(inputs, &TXInput{vin.PrevTXHash, vin.PrevTXIndex, nil, nil})
	}
	for _, vout := range tx.Output {
		outputs = append(outputs, &TXOutput{vout.Value, vout.PubKeyHash})
	}
	txCopy := Transaction{tx.Version, []byte{}, inputs, outputs}
	return txCopy
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
	tx.Version = Version
	tx.BlockHash = []byte{}
	// tx.Input = inputs

	outputs := make([]*TXOutput, 2)
	outputs = append(outputs, []*TXOutput{NewTXOutput(value, to)}...)
	//if diffValue > 0 {
	//outputs = append(outputs, []*TXOutput{NewTXOutput(from, diffValue)}...)
	//}
	tx.Output = outputs
	return &tx
}

func NewCoinbaseTX(value int, to Address) *Transaction {
	txi := &TXInput{[]byte{}, -1, []byte{}, []byte{}}
	txo := NewTXOutput(value, to)
	var tx *Transaction
	tx.Version = Version
	tx.Input = []*TXInput{txi}
	tx.Output = []*TXOutput{txo}
	return tx
}
