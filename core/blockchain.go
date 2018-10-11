package core

import (
	"log"
	"os"

	"github.com/tidwall/buntdb"
)

const blocksBucket = "gobucket"

type Blockchain struct {
	tip Hash
}

func (bc *Blockchain) AddBlock(tx []*Transaction) {
	var lastHash Hash
	db := getBlockchainDatabase()
	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("l")
		if err != nil {
			return err
		}
		lastHash = BytesToHash([]byte(val))
		return nil
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	newBlock := NewBlock(tx, lastHash)

	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(newBlock.Hash.String(), string(newBlock.Serialize()), nil)
		_, _, err = tx.Set("l", newBlock.Hash.String(), nil)
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
	db := getBlockchainDatabase()
	var findData bool
	findData = false
	err := db.Update(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(k, v string) bool {
			findData = true
			return false
		})
		if !findData {
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

func getBlockchainDatabase() *buntdb.DB {
	db, err := buntdb.Open(dbFile)
	errorHandle(err)
	return db
}
