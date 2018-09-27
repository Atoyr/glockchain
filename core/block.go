package core

import (
	"crypto/sha256"
	"time"

	"github.com/atoyr/gochain/util"
)

type Block struct {
	Index        int64
	PreviousHash []byte
	Timestamp    int64
	Hash         []byte
	Nonce        int
	Transactions []*Transaction
}

func (b *Block) ToTransactionsByte() []byte {
	bytes := make([]byte, 0, 10)
	for _, v := range b.Transactions {
		bytes = append(bytes, util.Interface2bytes(v.ToByte())...)
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
func NewGenesisBlock() *Block {
	return NewBlock([]*Transaction{}, []byte{})
}
