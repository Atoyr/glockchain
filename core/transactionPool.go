package core

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
)

// TransactionPool Not take up transaction into block
type TransactionPool struct {
	Pool []*Transaction
}

// NewTransactionPool TransactionPool constructor
func NewTransactionPool() *TransactionPool {
	var txp TransactionPool
	if dbExists(dbFile) == false {
		log.Println("Not exists db file")
		os.Exit(1)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(txpoolBucket))
		errorHandle(err)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tran := DeserializeTransaction(v)
			txp.Pool = append(txp.Pool, &tran)
		}
		return nil
	})
	errorHandle(err)
	return &txp
}

// AddTransaction Pool Transaciton
func (txp *TransactionPool) AddTransaction(transaction *Transaction) {
	if dbExists(dbFile) == false {
		log.Println("Not exists db file")
		os.Exit(1)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(txpoolBucket))
		errorHandle(err)
		err = b.Put(transaction.Hash(), transaction.Serialize())
		errorHandle(err)
		return nil
	})
	errorHandle(err)
	txp.Pool = append(txp.Pool, transaction)
}

// ClearTransactionPool clear pool
func (txp *TransactionPool) ClearTransactionPool() {
	if dbExists(dbFile) == false {
		log.Println("Not exists db file")
		os.Exit(1)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(txpoolBucket))
		_, err := tx.CreateBucket([]byte(txpoolBucket))
		errorHandle(err)
		return nil
	})
	errorHandle(err)
	txp.Pool = make([]*Transaction, 0, 1)
}
