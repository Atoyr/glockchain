package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"os"

	"github.com/atoyr/glockchain/util"
)

type Transaction struct {
	version     int
	ID          []byte
	blockHash   []byte
	blockNumber []byte
	Input       []TXInput
	Output      []TXOutput
}

func (t *Transaction) ToByte() []byte {
	bytes := make([]byte, 100)
	bytes = append(bytes, util.Interface2bytes(t.version)...)
	return bytes
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
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}
