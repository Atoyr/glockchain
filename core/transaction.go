package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
)

// Transaction TX Data
type Transaction struct {
	Version   byte
	BlockHash []byte
	Input     []*TXInput
	Output    []*TXOutput
}

// Hash Hash to transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// Sign sign TX
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey) {
	txCopy := tx.TrimmedCopy()
	for vindex, _ := range txCopy.Input {
		txCopy.Input[vindex].Signature = []byte{}
		txCopyHash := txCopy.Hash()
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopyHash)
		errorHandle(err)
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Input[vindex].Signature = signature
		txCopy.Input[vindex].Signature = nil
	}
}

// Verify verify TX
func (tx *Transaction) Verify() bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for vindex, vin := range txCopy.Input {
		hashPubKey := HashPubKey(vin.PubKey)
		txCopy.Input[vindex].Signature = hashPubKey
		txHash := txCopy.Hash()

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, txHash, &r, &s) == false {
			return false
		}
	}
	return true
}

// TrimmedCopy copy TX (not in TXI pubkey and Signature)
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

func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.Hash()))
	for i, in := range tx.Input {
		lines = append(lines, fmt.Sprintf("  Input %d:", i))
		lines = append(lines, fmt.Sprintf("    PrevTX      %x", in.PrevTXHash))
		lines = append(lines, fmt.Sprintf("    PrevTXIndex %d", in.PrevTXIndex))
		lines = append(lines, fmt.Sprintf("    Signature   %x", in.Signature))
		lines = append(lines, fmt.Sprintf("    PubKey      %x", in.PubKey))
	}

	for i, out := range tx.Output {
		lines = append(lines, fmt.Sprintf("  Output %d:", i))
		lines = append(lines, fmt.Sprintf("    Value       %d", out.Value))
		lines = append(lines, fmt.Sprintf("    PubKeyHash  %x", out.PubKeyHash))
	}
	return strings.Join(lines, "\n")
}

// Serialize serialize tx
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

// DeserializeTransaction deserialize tx
func DeserializeTransaction(data []byte) Transaction {
	var tx Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}
	return tx
}

// NewTransaction create New TX
func NewTransaction(utxos []*UTXO, from, to Address, value int, returnValue int) *Transaction {
	var tx Transaction
	txp := GetTransactionPool()
	tx.Input = make([]*TXInput, len(utxos))
	sumValue := 0
	for _, utxo := range utxos {
		sumValue += utxo.TX.Output[utxo.Index].Value
		var txin *TXInput
		txin.PrevTXHash = utxo.TX.Hash()
		txin.PrevTXIndex = utxo.Index
		tx.Input = append(tx.Input, txin)
	}

	diffValue := sumValue - value - returnValue
	if diffValue < 0 {
		return nil
	}
	tx.Version = Version
	tx.BlockHash = []byte{}

	outputs := make([]*TXOutput, 2)
	outputs = append(outputs, []*TXOutput{NewTXOutput(value, to)}...)
	if diffValue > 0 {
		outputs = append(outputs, []*TXOutput{NewTXOutput(diffValue, from)}...)
	}
	tx.Output = outputs
	txp.AddTransaction(&tx)
	return &tx
}

// NewCoinbaseTX Create New Coinbase TX
func NewCoinbaseTX(value int, to Address) *Transaction {
	txi := &TXInput{[]byte{}, -1, []byte{}, []byte{}}
	txo := NewTXOutput(value, to)
	var tx Transaction
	tx.Version = 0x00
	tx.Input = []*TXInput{txi}
	tx.Output = []*TXOutput{txo}
	wallets := NewWallets()
	log.Println(to)
	wallet := wallets.GetWallet(to)
	tx.Sign(wallet.PrivateKey)
	return &tx
}
