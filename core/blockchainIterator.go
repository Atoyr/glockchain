package core

import (
	"log"
	"os"

	"github.com/tidwall/buntdb"
)

type BlockchainIterator struct {
	currentHash Hash
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip}
	return bci
}

func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	db := getBlockchainDatabase()
	err := db.View(func(tx *buntdb.Tx) error {
		encodedBlock, err := tx.Get(bci.currentHash.String())
		errorHandle(err)
		block = DeserializeBlock([]byte(encodedBlock))
		return nil
	})
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	bci.currentHash = block.PreviousHash
	return block
}
