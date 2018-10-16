package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"os"
	"time"

	"github.com/atoyr/glockchain/util"
)

type Block struct {
	Index        int64
	PreviousHash []byte
	Timestamp    int64
	Hash         []byte
	Nonce        int64
	Transactions []*Transaction
}

func (b *Block) ToTransactionsByte() []byte {
	bytes := make([]byte, 0, 10)
	for _, v := range b.Transactions {
		bytes = append(bytes, v.Hash()...)
	}
	return bytes
}

func (b *Block) ToByte() []byte {
	bytes := make([]byte, 0, 10)
	bytes = append(bytes, util.Interface2bytes(b.Index)...)
	bytes = append(bytes, util.Interface2bytes(b.PreviousHash)...)
	bytes = append(bytes, util.Interface2bytes(b.Timestamp)...)
	bytes = append(bytes, util.Interface2bytes(b.Hash)...)
	bytes = append(bytes, util.Interface2bytes(b.Nonce)...)
	bytes = append(bytes, b.ToTransactionsByte()...)
	return bytes
}

func (b *Block) SetHash() {
	hash := sha256.Sum256(b.ToByte())
	b.Hash = hash[:]
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), Transactions: transactions, PreviousHash: prevBlockHash}
	block.SetHash()
	return block
}
func NewGenesisBlock(tx *Transaction) *Block {
	return NewBlock([]*Transaction{tx}, []byte{})
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err == nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return &block
}
