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

type UTXOPool struct {
	Pool map[string]*UTXO
}

var utxopool *UTXOPool
var once sync.Once

func NewUTXOPool() *UTXOPool {
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

func GetUTXOPool() *UTXOPool {
	once.Do(func() {
		utxopool = NewUTXOPool()
	})
	return utxopool
}

func (up *UTXOPool) AddUTXO(utxo *UTXO) {
	if dbExists(dbFile) == false {
		log.Println("Not exists db file")
		os.Exit(1)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(utxoBucket))
		errorHandle(err)
		err = b.Put(append(utxo.TX.Hash(), util.Int2bytes(utxo.Index, 8)...), utxo.Serialize())
		errorHandle(err)
		return nil
	})
	errorHandle(err)
	up.Pool[hex.EncodeToString(utxo.Key())] = utxo
}

func (up *UTXOPool) FindUTXOs(pubKeyHash []byte, amount int) (int, map[string]UTXO) {
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

func (up *UTXOPool) GetUTXO(txhash []byte, index int) *UTXO {
	if dbExists(dbFile) == false {
		log.Println("Not exists db file")
		os.Exit(1)
	}
	var utxo UTXO
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		encodeUtxo := b.Get(append(txhash, util.Int2bytes(index, 8)...))
		utxo = DeserializeUtxo(encodeUtxo)
		return nil
	})
	errorHandle(err)
	return &utxo
}

func (up *UTXOPool) String() string {
	var lines []string
	for _, utxo := range up.Pool {
		lines = append(lines, utxo.String())
		lines = append(lines, utxo.TX.String())
	}
	return strings.Join(lines, "\n")
}
