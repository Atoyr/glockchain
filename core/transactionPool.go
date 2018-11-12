package core

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type TransactionPool struct {
	Pool []*Transaction
}

func GetTransactionPool() *TransactionPool {
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
