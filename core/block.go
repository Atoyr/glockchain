package core

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/atoyr/gochain/util"
)

type Block struct {
	Index        int64
	PreviousHash string
	Timestamp    int64
	Hash         string
	Nonce        int
	Transactions []Transaction
}

func (b *Block) ToTransactionByte() []byte {
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
	bytes = append(bytes, b.ToTransactionByte()...)
	return bytes
}

func (b *Block) ToHashString() string {
	converted := sha256.Sum256(b.ToByte())
	return hex.EncodeToString(converted[:])
}

var (
	blockchain = []*Block{}
)
var genesisBlock = &Block{
	Index:        0,
	PreviousHash: "0",
	Timestamp:    0,
}
