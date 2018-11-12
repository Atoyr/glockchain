package core

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type Blockchain struct {
	tip []byte
}

func (bc *Blockchain) AddBlock(block *Block) {
	//var lastHash []byte
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b.Get(block.Hash) != nil {
			return nil
		}

		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		errorHandle(err)
		err = b.Put([]byte("l"), block.Hash)
		return nil

	})
	errorHandle(err)
}

func CreateBlockchain(address []byte) *Blockchain {
	if dbExists(dbFile) {
		log.Println("Exist db file")
		os.Exit(1)
	}
	var tip []byte
	cbtx := NewCoinbaseTX(100, address)
	genesis := NewGenesisBlock(cbtx)
	db := getBlockchainDatabase()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		errorHandle(err)
		err = b.Put(genesis.Hash, genesis.Serialize())
		errorHandle(err)
		err = b.Put([]byte("l"), genesis.Hash)
		errorHandle(err)
		tip = genesis.Hash
		return nil
	})
	errorHandle(err)
	db.Close()

	up := GetUTXOPool()
	up.AddUTXO(genesis.Transactions[0])

	bc := Blockchain{tip}
	return &bc
}

func GetBlockchain() (*Blockchain, []byte) {
	var tip []byte
	if dbExists(dbFile) == false {
		fmt.Println("Not exist db file")
		os.Exit(1)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})
	errorHandle(err)
	bc := Blockchain{tip}
	tip2 := make([]byte, len(tip))
	copy(tip2, tip)
	return &bc, tip2
}

func getBlockchainDatabase() *bolt.DB {
	db, err := bolt.Open(dbFile, 0600, nil)
	errorHandle(err)
	return db
}

func dbExists(dbfile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
