package core

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

// Blockchain is chain with block
// tip is last block hash
type Blockchain struct {
	tip []byte
}

// CreateBlockchain is create blockchain
// if blockchain exists, error it
func CreateBlockchain(address []byte) *Blockchain {
	// TODO:(atoyr):return error
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

// GetBlockchain is Getting exist blockchain
func GetBlockchain() (*Blockchain, []byte) {
	var tip []byte
	// TODO:(atoyr):return error
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

// AddBlock is adding block into blockchain
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

// Iterator is created blockchain iterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := BlockchainIterator{bc.tip}
	return &bci
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
