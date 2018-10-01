package core

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "glockchain_%s.db"
const blocksBucket = "gobucket"

type Blockchain struct {
	Blocks []*Block
	pool   []*Transaction
}

func (bc *Blockchain) AddBlock() {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(bc.pool[:], prevBlock.Hash)
	bc.pool = make([]*Transaction, 0, 0)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlockchain() *Blockchain {
	dbFile := fmt.Sprint(dbFile, "0")
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			tip = genesis.Hash
		}
		return nil
	})
	return &Blockchain{[]*Block{NewGenesisBlock()}, []*Transaction{}}
}
