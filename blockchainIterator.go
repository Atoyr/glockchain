package glockChain

import (
	"github.com/boltdb/bolt"
)

// BlockchainIterator is a blockchain iterator
type BlockchainIterator struct {
	currentHash []byte
}

// Next is moved index for next
func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	db := getBlockchainDatabase()
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	bci.currentHash = block.PreviousHash
	return block
}
