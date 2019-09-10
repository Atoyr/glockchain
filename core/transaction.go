package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/cbergoon/merkletree"
	"github.com/pkg/errors"
)

// Transaction TX Data
type Transaction struct {
	Version byte       `json:"version"`
	ID      []byte     `json:"id"`
	Input   []TXInput  `json:"input"`
	Output  []TXOutput `json:"output"`
}

// IsCoinbase is return this transaction is coinbase?
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Input) == 1 && len(tx.Input[0].PrevTXHash) == 0 && tx.Input[0].PrevTXIndex == -1
}

// Bytes Get Transaction Bytes
func (tx *Transaction) Bytes() []byte {
	var b []byte
	txCopy := *tx
	b = append(b, txCopy.Version)
	for _, in := range txCopy.Input {
		b = append(b, in.Hash()...)
	}
	for _, out := range txCopy.Output {
		b = append(b, out.Hash()...)
	}
	return b
}

// Hash Hash to transaction
func (tx *Transaction) Hash() []byte {
	b := tx.Bytes()
	var hash, hash2 [32]byte
	hash = sha256.Sum256(b)
	hash2 = sha256.Sum256(hash[:])
	return hash2[:]
}

// CalculateHash is Hash for Merkletree
func (tx Transaction) CalculateHash() ([]byte, error) {
	return tx.Hash(), nil
}

func (tx Transaction) Equals(other merkletree.Content) (bool, error) {
	if temptx, ok := other.(Transaction); ok {
		return fmt.Sprintf("%x", tx.ID) == fmt.Sprintf("%x", temptx.ID), nil
	}
	return false, NewGlockchainError(10001)
}

// Sign sign TX
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey) error {
	if tx.IsCoinbase() {
		return nil
	}
	txCopy := tx.TrimmedCopy()
	for vindex := range txCopy.Input {
		txCopy.Input[vindex].Signature = []byte{}
		dataToSign := fmt.Sprintf("%x\n", txCopy)
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			return errors.Wrap(err, getErrorMessage(94001))
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Input[vindex].Signature = signature
		txCopy.Input[vindex].Signature = nil
	}
	return nil
}

// Verify verify TX
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	var getXY func(b []byte) (x, y big.Int)
	getXY = func(b []byte) (x, y big.Int) {
		sigLen := len(b)
		x.SetBytes(b[:(sigLen / 2)])
		y.SetBytes(b[(sigLen / 2):])
		return
	}
	for vindex, vin := range tx.Input {
		prevTX := prevTXs[hex.EncodeToString(vin.PrevTXHash)]
		hashPubKey := prevTX.Output[vin.PrevTXIndex].PubKeyHash
		txCopy.Input[vindex].PubKey = hashPubKey
		dataToVerify := fmt.Sprintf("%x\n", txCopy)
		r, s := getXY(vin.Signature)
		x, y := getXY(vin.PubKey)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
	}
	return true
}

// TrimmedCopy copy TX (not in TXI pubkey and Signature)
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Input {
		inputs = append(inputs, TXInput{vin.PrevTXHash, vin.PrevTXIndex, nil, nil})
	}
	for _, vout := range tx.Output {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}
	txCopy := Transaction{tx.Version, tx.ID, inputs, outputs}
	return txCopy
}

func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Transaction  : %x", tx.Hash()))
	lines = append(lines, fmt.Sprintf("  version    : %x", tx.Version))
	lines = append(lines, fmt.Sprintf("  ID         : %x", tx.ID))
	lines = append(lines, fmt.Sprintf("  Inputs     : %d", len(tx.Input)))
	for i, in := range tx.Input {
		lines = append(lines, fmt.Sprintf("    Input %d", i))
		lines = append(lines, fmt.Sprintf("      PrevTX         : %x", in.PrevTXHash))
		lines = append(lines, fmt.Sprintf("      PrevTXIndex    : %d", in.PrevTXIndex))
		lines = append(lines, fmt.Sprintf("      Signature      : %x", in.Signature))
		lines = append(lines, fmt.Sprintf("      PubKey         : %x", in.PubKey))
	}
	lines = append(lines, fmt.Sprintf("  Outputs    : %d", len(tx.Output)))
	for i, out := range tx.Output {
		lines = append(lines, fmt.Sprintf("    Output %d", i))
		lines = append(lines, fmt.Sprintf("      Value          : %d", out.Value))
		lines = append(lines, fmt.Sprintf("      PubKeyHash     : %x", out.PubKeyHash))
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
func NewTransaction(wallet *Wallet, to []byte, amount int) (*Transaction, error) {
	var inputs []TXInput
	var outputs []TXOutput
	pubKeyHash := HashPubKey(wallet.PublicKey)
	utxopool, err := GetUTXOPool()
	if err != nil {
		return nil, err
	}

	acc, utxos := utxopool.FindSpendableOutputs(pubKeyHash, amount)
	if acc < amount {
		return nil, NewGlockchainError(93003)
	}

	inputs = make([]TXInput, len(utxos))
	index := 0
	for _, utxo := range utxos {
		var txin TXInput
		txin.PrevTXHash = utxo.TX.ID
		txin.PrevTXIndex = utxo.Index
		txin.PubKey = wallet.PublicKey
		inputs[index] = txin
		index++
	}
	diffamount := acc - amount

	outputs = make([]TXOutput, 1)
	outputs[0] = *NewTXOutput(amount, to)
	if diffamount > 0 {
		outputs = append(outputs, *NewTXOutput(diffamount, wallet.GetAddress()))
	}

	tx := Transaction{Version, []byte{}, inputs, outputs}
	tx.ID = tx.Hash()

	err = tx.Sign(wallet.PrivateKey)
	if err != nil {
		return nil, err
	}
	txp, err := NewTransactionPool()
	if err != nil {
		return nil, err
	}
	txp.AddTransaction(&tx)
	utxopool.AddUTXO(&tx)
	return &tx, nil
}

// NewCoinbaseTX Create New Coinbase TX
func NewCoinbaseTX(value int, wallet *Wallet) (*Transaction, error) {
	txi := &TXInput{[]byte{}, -1, []byte{}, []byte{}}
	txo := NewTXOutput(value, wallet.GetAddress())
	var tx Transaction
	tx.Version = 0x00
	tx.Input = []TXInput{*txi}
	tx.Output = []TXOutput{*txo}
	tx.ID = tx.Hash()
	err := tx.Sign(wallet.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
