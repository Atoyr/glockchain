package core

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "glockchains.db"
const blocksBucket = "gobucket"

type Blockchain struct {
	tip Hash
	db  *bolt.DB
}

func (bc *Blockchain) AddBlock(tx []*Transaction) {
	var lastHash Hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = BytesToHash(b.Get([]byte("l")))
		return nil
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	newBlock := NewBlock(tx, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash.Bytes(), newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash.Bytes())
		if err != nil {
			log.Panic(err)
			os.Exit(1)
		}
		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func NewBlockchain() *Blockchain {
	dbFile := fmt.Sprint(dbFile, "")
	var tip Hash
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			var t *Transaction
			genesis := NewGenesisBlock(t)
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash.Bytes(), genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash.Bytes())
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			tip = genesis.Hash
		}
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc
}
