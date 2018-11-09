package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
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

func (block *Block) ToHash() []byte {
	var b []byte
	b = append(b, block.PreviousHash...)
	b = append(b, block.HashTransactions()...)
	b = append(b, util.Int642bytes(block.Timestamp)...)
	b = append(b, util.Int642bytes(block.Nonce)...)

	hash := sha256.Sum256(b)
	return hash[:]
}

func (b *Block) SetHash() {
	hash := sha256.Sum256(b.Serialize())
	b.Hash = hash[:]
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), Transactions: transactions, PreviousHash: prevBlockHash}
	block.SetHash()
	for _, tx := range block.Transactions {
		tx.BlockHash = block.Hash
	}
	return block
}
func NewGenesisBlock(tx *Transaction) *Block {
	block := NewBlock([]*Transaction{tx}, []byte{})
	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	errorHandle(err)
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
func (b *Block) HashTransactions() []byte {
	var txbytes [][]byte

	for _, tx := range b.Transactions {
		txbytes = append(txbytes, tx.Bytes())
	}
	mtree := NewMerkleTree(txbytes)
	return mtree.RootNode.Data
}

func (b *Block) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Block : %x", b.Hash))
	lines = append(lines, fmt.Sprintf("index : %d", b.Index))
	lines = append(lines, fmt.Sprintf("prev : %x", b.PreviousHash))
	lines = append(lines, fmt.Sprintf("timestamp : %x", b.Timestamp))
	lines = append(lines, fmt.Sprintf("nonce : %d", b.Nonce))
	for _, tx := range b.Transactions {
		lines = append(lines, tx.String())
	}
	return strings.Join(lines, "\n")
}
