package core

import (
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// TransactionPool Not take up transaction into block
type TransactionPool struct {
	Pool []*Transaction
}

// NewTransactionPool TransactionPool constructor
func NewTransactionPool() (*TransactionPool, error) {
	var txp TransactionPool
	if dbExists(dbFile) == false {
		return nil, NewGlockchainError(91001)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(txpoolBucket))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tran := DeserializeTransaction(v)
			txp.Pool = append(txp.Pool, &tran)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, getErrorMessage(91003))
	}
	return &txp, nil
}

// AddTransaction Pool Transaciton
func (txp *TransactionPool) AddTransaction(transaction *Transaction) error {
	if dbExists(dbFile) == false {
		return NewGlockchainError(91001)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(txpoolBucket))
		if err != nil {
			return err
		}
		err = b.Put(transaction.Hash(), transaction.Serialize())
		return err
	})
	if err != nil {
		return errors.Wrap(err, getErrorMessage(91003))
	}
	txp.Pool = append(txp.Pool, transaction)
	return nil
}

// ClearTransactionPool clear pool
func (txp *TransactionPool) ClearTransactionPool() error {
	if dbExists(dbFile) == false {
		return NewGlockchainError(91001)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(txpoolBucket))
		_, err := tx.CreateBucket([]byte(txpoolBucket))
		return err
	})
	if err != nil {
		return errors.Wrap(err, getErrorMessage(91003))
	}
	txp.Pool = make([]*Transaction, 0, 1)
	return nil
}
