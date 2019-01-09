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

// Block is transaction block
type Block struct {
	Index        int64
	PreviousHash []byte
	Timestamp    int64
	Hash         []byte
	Nonce        int
	Transactions []*Transaction
	TXHash       []byte
}

// NewBlock is block constructor
// transactions is the transaction you want to add to the block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) (block *Block, err error) {
	block = &Block{}
	block.Timestamp = time.Now().Unix()
	block.Transactions = transactions
	block.PreviousHash = prevBlockHash
	pow, err := NewProofOfWork(block)
	if err != nil {
		return nil, err
	}
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash
	block.TXHash = block.HashTransactions()
	return
}

// ToHash converts block to hash
func (block *Block) ToHash() []byte {
	var b []byte
	b = append(b, block.PreviousHash...)
	b = append(b, block.HashTransactions()...)
	b = append(b, util.Int642bytes(block.Timestamp)...)
	b = append(b, util.Int2bytes(block.Nonce, 8)...)

	hash := sha256.Sum256(b)
	return hash[:]
}

// NewGenesisBlock is created genesis block
// tx is the coinbase transaction
func NewGenesisBlock(tx *Transaction) *Block {
	block, _ := NewBlock([]*Transaction{tx}, []byte{})
	return block
}

// Serialize block
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	errorHandle(err)
	return result.Bytes()
}

// DeserializeBlock is deserializing block
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

// HashTransactions is hasing transacitons from block
func (block *Block) HashTransactions() []byte {
	var txbytes [][]byte

	for _, tx := range block.Transactions {
		txbytes = append(txbytes, tx.Bytes())
	}
	mtree := NewMerkleTree(txbytes)
	return mtree.RootNode.Data
}

// String is convert block to string
func (block *Block) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Block : %x", block.Hash))
	lines = append(lines, fmt.Sprintf("index : %d", block.Index))
	lines = append(lines, fmt.Sprintf("prev : %x", block.PreviousHash))
	lines = append(lines, fmt.Sprintf("timestamp : %x", block.Timestamp))
	lines = append(lines, fmt.Sprintf("nonce : %d", block.Nonce))
	for _, tx := range block.Transactions {
		lines = append(lines, tx.String())
	}
	return strings.Join(lines, "\n")
}
