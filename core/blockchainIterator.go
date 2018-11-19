package core

import (
	"log"

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
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bci.currentHash = block.PreviousHash
	return block
}
