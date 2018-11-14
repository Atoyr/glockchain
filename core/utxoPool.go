package core

import (
	"encoding/hex"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/atoyr/glockchain/util"
	"github.com/boltdb/bolt"
)

// UTXOPool is singleton
type UTXOPool struct {
	Pool map[string]*UTXO
}

var utxopool *UTXOPool
var once sync.Once

func newUTXOPool() *UTXOPool {
	var up UTXOPool
	up.Pool = make(map[string]*UTXO)
	if dbExists(dbFile) == false {
		log.Println("Not exists db file")
		os.Exit(1)
	}

	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(utxoBucket))
		errorHandle(err)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := hex.EncodeToString(k)
			utxo := DeserializeUtxo(v)
			up.Pool[key] = &utxo
		}
		return nil
	})
	errorHandle(err)
	return &up
}

// GetUTXOPool is create singleton instance
func GetUTXOPool() *UTXOPool {
	once.Do(func() {
		utxopool = newUTXOPool()
	})
	return utxopool
}

// AddUTXO is create UTXO and Pooling UTXP
func (up *UTXOPool) AddUTXO(t *Transaction) {
	if dbExists(dbFile) == false {
		log.Println("Not exists db file")
		os.Exit(1)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(utxoBucket))
		errorHandle(err)
		for _, input := range t.Input {
			key := getUTXOPoolKey(input.PrevTXHash, input.PrevTXIndex)
			b.Delete(key)
			delete(up.Pool, hex.EncodeToString(key))
		}
		for i := range t.Output {
			utxo := UTXO{t, i}
			err = b.Put(utxo.Key(), utxo.Serialize())
			errorHandle(err)
			up.Pool[hex.EncodeToString(utxo.Key())] = &utxo
		}
		errorHandle(err)
		return nil
	})
	errorHandle(err)
}

// FindSpendableOutputs is finc tx from pool with pubkeyhash ando amount
func (up *UTXOPool) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string]UTXO) {
	utxos := make(map[string]UTXO)
	acc := 0
	db := getBlockchainDatabase()
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
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
	errorHandle(err)
	return acc, utxos
}

// FindUTXOs is find tx from pool with pubkeyhash
func (up *UTXOPool) FindUTXOs(pubKeyHash []byte) (int, map[string]UTXO) {
	utxos := make(map[string]UTXO)
	balance := 0
	db := getBlockchainDatabase()
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
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
	errorHandle(err)
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
func getUTXOPoolKey(txhash []byte, index int) []byte {
	return append(txhash, util.Int2bytes(index, 8)...)
}
