package gochain

import "crypto/sha256"

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
	for i, v := range b.Transactions {
		bytes = append(bytes, []byte(v.ToByte()))
	}
	return bytes
}

func (b *Block) ToByte() []byte {
	bytes := make([]byte, 0, 10)
	bytes = append(bytes, []byte(b.Index))
	bytes = append(bytes, []byte(b.PreviousHash))
	bytes = append(bytes, []byte(b.Timestamp))
	bytes = append(bytes, []byte(b.Hash))
	bytes = append(bytes, []byte(b.Nonce))
	bytes = append(bytes, []byte(b.ToTransactionByte()))
	return bytes
}

func (b *Block) ToHashString() string {
	return sha256.Sum256([]byte(b.ToByte()))
}

var (
	blockchain = []*Block{}
)
var genesisBlock = &Block{
	Index:        0,
	PreviousHash: "0",
	Timestamp:    0,
}
