package core

import (
	"log"
	"os"

	"github.com/tidwall/buntdb"
)

const blocksBucket = "gobucket"

type Blockchain struct {
	tip []byte
}

func (bc *Blockchain) AddBlock(tx []*Transaction) {
	var lastHash []byte
	db := getBlockchainDatabase()
	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("l")
		if err != nil {
			return err
		}
		lastHash = []byte(val)
		return nil
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	newBlock := NewBlock(tx, lastHash)

	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(string(newBlock.Hash), string(newBlock.Serialize()), nil)
		_, _, err = tx.Set("l", string(newBlock.Hash), nil)
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
	var tip []byte
	err := db.Update(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(k, v string) bool {
			findData = true
			return false
		})
		log.Println("hogehoge")
		log.Println(findData)
		if !findData {
			var addr Address
			genesis := NewGenesisBlock(NewCoinbaseTX(100, addr))
			_, _, err := tx.Set(string(genesis.Hash), string(genesis.Serialize()), nil)
			_, _, err = tx.Set("l", string(genesis.Hash), nil)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			log.Println(genesis.Hash)
			tip = genesis.Hash
		}
		return nil
	})
	errorHandle(err)
	bc := Blockchain{tip}
	return &bc
}

func getBlockchainDatabase() *buntdb.DB {
	db, err := buntdb.Open(dbFile)
	errorHandle(err)
	return db
}
