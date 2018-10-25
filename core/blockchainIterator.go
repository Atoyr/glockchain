package core

import (
	"log"

	"github.com/boltdb/bolt"
)

type BlockchainIterator struct {
	currentHash []byte
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := BlockchainIterator{bc.tip}
	return &bci
}

func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	db := getBlockchainDatabase()
	defer db.Close()
	log.Println(bci.currentHash)
	err := db.View(func(tx *bolt.Tx) error {
		log.Println(bci.currentHash)
		b := tx.Bucket([]byte(blocksBucket))
		log.Println(bci.currentHash)
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
