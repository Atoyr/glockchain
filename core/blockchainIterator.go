package core

import (
	"log"

	"github.com/tidwall/buntdb"
)

type BlockchainIterator struct {
	currentHash []byte
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip}
	return bci
}

func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	db := getBlockchainDatabase()
	err := db.View(func(tx *buntdb.Tx) error {
		encodedBlock, err := tx.Get(string(bci.currentHash))
		errorHandle(err)
		block = DeserializeBlock([]byte(encodedBlock))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bci.currentHash = block.PreviousHash
	return block
}
