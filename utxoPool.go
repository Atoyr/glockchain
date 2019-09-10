package glockchain

import (
	"encoding/hex"
	"strings"
	"sync"

	"github.com/atoyr/glockchain/util"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// UTXOPool is singleton
type UTXOPool struct {
	Pool map[string]*UTXO
}

var utxopool *UTXOPool
var once sync.Once

func newUTXOPool() (*UTXOPool, error) {
	var up UTXOPool
	up.Pool = make(map[string]*UTXO)
	if dbExists(dbFile) == false {
		return nil, NewGlockchainError(91001)
	}

	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(utxoBucket))
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := hex.EncodeToString(k)
			utxo := DeserializeUtxo(v)
			up.Pool[key] = &utxo
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, getErrorMessage(91003))
	}
	return &up, nil
}

// GetUTXOPool is create singleton instance
func GetUTXOPool() (*UTXOPool, error) {
	var err error
	once.Do(func() {
		utxopool, err = newUTXOPool()
	})
	if err != nil {
		return nil, err
	}

	return utxopool, nil
}

// AddUTXO is create UTXO and Pooling UTXP
func (up *UTXOPool) AddUTXO(t *Transaction) error {
	if dbExists(dbFile) == false {
		return NewGlockchainError(91001)
	}

	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(utxoBucket))
		if err != nil {
			return err
		}
		for _, input := range t.Input {
			key := getUTXOPoolKey(input.PrevTXHash, input.PrevTXIndex)
			b.Delete(key)
			delete(up.Pool, hex.EncodeToString(key))
		}
		for i := range t.Output {
			utxo := UTXO{t, i}
			err = b.Put(utxo.Key(), utxo.Serialize())
			if err != nil {
				return err
			}
			up.Pool[hex.EncodeToString(utxo.Key())] = &utxo
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, getErrorMessage(91003))
	}
	return nil
}

// FindSpendableOutputs is finc tx from pool with pubkeyhash ando amount
func (up *UTXOPool) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string]UTXO) {
	utxos := make(map[string]UTXO)
	acc := 0
	db := getBlockchainDatabase()
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := hex.EncodeToString(k)
			utxo := DeserializeUtxo(v)
			if utxo.TX.Output[utxo.Index].IsLockedWithKey(pubKeyHash) && acc < amount {
				acc = acc + utxo.TX.Output[utxo.Index].Value
				utxos[key] = utxo
			}
			if amount <= acc {
				return nil
			}
		}
		return nil
	})
	return acc, utxos
}

// FindUTXOs is find tx from pool with pubkeyhash
func (up *UTXOPool) FindUTXOs(pubKeyHash []byte) (int, map[string]UTXO) {
	utxos := make(map[string]UTXO)
	balance := 0
	db := getBlockchainDatabase()
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := hex.EncodeToString(k)
			utxo := DeserializeUtxo(v)
			if utxo.TX.Output[utxo.Index].IsLockedWithKey(pubKeyHash) {
				balance = balance + utxo.TX.Output[utxo.Index].Value
				utxos[key] = utxo
			}
		}
		return nil
	})
	return balance, utxos
}

// String is convert UTXOPool to string
func (up *UTXOPool) String() string {
	var lines []string
	for _, utxo := range up.Pool {
		lines = append(lines, utxo.String())
	}
	return strings.Join(lines, "\n")
}
func getUTXOPoolKey(txid []byte, index int) []byte {
	return append(txid, util.Int2bytes(index, 8)...)
}
