package core

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// Blockchain is chain with block
// tip is last block hash
type Blockchain struct {
	tip []byte
}

// CreateBlockchain is create blockchain
// if blockchain exists, error it
func CreateBlockchain(address []byte) (*Blockchain, error) {
	if dbExists(dbFile) {
		return nil, NewGlockchainError(91002)
	}
	var tip []byte
	cbtx, err := NewCoinbaseTX(100, address)
	if err != nil {
		return nil, errors.Wrap(err, getErrorMessage(93002))
	}

	genesis := NewGenesisBlock(cbtx)
	db := getBlockchainDatabase()
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			return err
		}
		tip = genesis.Hash
		return nil
	})
	db.Close()
	if err != nil {
		return nil, errors.Wrap(err, getErrorMessage(91003))
	}

	up, err := GetUTXOPool()
	if err != nil {
		return nil, err
	}
	up.AddUTXO(genesis.Transactions[0])

	bc := Blockchain{tip}
	return &bc, nil
}

// GetBlockchain is Getting exist blockchain
func GetBlockchain() (*Blockchain, []byte, error) {
	var tip []byte
	// TODO:(atoyr):return error
	if dbExists(dbFile) == false {
		return nil, nil, NewGlockchainError(91001)
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
	return &bc, tip2, nil
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
