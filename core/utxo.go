package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
)

// UTXO UTXO Data (tryal)
type UTXO struct {
	TX    *Transaction
	Index int
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
