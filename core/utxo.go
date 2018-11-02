package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atoyr/glockchain/util"
)

// UTXO UTXO Data (tryal)
type UTXO struct {
	TX    *Transaction
	Index int
}

func (utxo *UTXO) Hash() []byte {
	var hash [32]byte
	utxoCopy := *utxo
	hash = sha256.Sum256(utxoCopy.Serialize())
	return hash[:]
}

func (utxo *UTXO) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("--- UTXO %x:", utxo.Hash()))
	lines = append(lines, fmt.Sprintf("  Output "))
	lines = append(lines, fmt.Sprintf("  TX    %x", utxo.TX.Hash()))
	lines = append(lines, fmt.Sprintf("  Index %d", utxo.Index))
	lines = append(lines, fmt.Sprintf("    Value       %d", utxo.TX.Output[utxo.Index].Value))
	lines = append(lines, fmt.Sprintf("    PubKeyHash  %x", utxo.TX.Output[utxo.Index].PubKeyHash))
	return strings.Join(lines, "\n")
}

func (utxo *UTXO) Key() []byte {
	return append(utxo.TX.Hash(), util.Int2bytes(utxo.Index, 8)...)
}

func (utxo *UTXO) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(utxo)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	return encoded.Bytes()
}

// DeserializeUtxo deserialize tx
func DeserializeUtxo(data []byte) UTXO {
	var utxo UTXO
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&utxo)
	if err != nil {
		log.Panic(err)
	}
	return utxo
}