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
	ID        []byte
	BlockHash []byte
	Input     []TXInput
	Output    []TXOutput
}

func (tx *Transaction) Bytes() []byte {
	var b []byte
	txCopy := *tx
	b = append(b, txCopy.Version)
	b = append(b, txCopy.BlockHash...)
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
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Input {
		inputs = append(inputs, TXInput{vin.PrevTXHash, vin.PrevTXIndex, nil, nil})
	}
	for _, vout := range tx.Output {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}
	txCopy := Transaction{tx.Version, tx.ID, []byte{}, inputs, outputs}
	return txCopy
}

func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Transaction  : %x", tx.Hash()))
	lines = append(lines, fmt.Sprintf("  version    : %x", tx.Version))
	lines = append(lines, fmt.Sprintf("  BlockHash  : %x", tx.BlockHash))
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
func NewTransaction(wallet *Wallet, to []byte, amount int) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	pubKeyHash := HashPubKey(wallet.PublicKey)
	utxopool := GetUTXOPool()
	acc, utxos := utxopool.FindSpendableOutputs(pubKeyHash, amount)
	if acc < amount {
		log.Panic("ERROR : Not enough funds")
	}

	inputs = make([]TXInput, len(utxos))
	index := 0
	for _, utxo := range utxos {
		var txin TXInput
		txin.PrevTXHash = utxo.TX.Hash()
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
	tx := Transaction{Version, []byte{}, []byte{}, inputs, outputs}
	tx.Sign(wallet.PrivateKey)
	txp := GetTransactionPool()
	txp.AddTransaction(&tx)
	utxopool.AddUTXO(&tx)
	return &tx
}

// NewCoinbaseTX Create New Coinbase TX
func NewCoinbaseTX(value int, to []byte) *Transaction {
	txi := &TXInput{[]byte{}, -1, []byte{}, []byte{}}
	txo := NewTXOutput(value, to)
	var tx Transaction
	tx.Version = 0x00
	tx.Input = []TXInput{*txi}
	tx.Output = []TXOutput{*txo}
	wallets := NewWallets()
	wallet := wallets.GetWallet(to)
	tx.Sign(wallet.PrivateKey)
	return &tx
}
