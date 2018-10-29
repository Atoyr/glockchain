package core

import (
	"log"
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

type UTXOPool struct {
	Pool []*UTXO
}

func GetUTXOPool() *UTXOPool {
	var up UTXOPool
	if dbExists(dbFile) == false {
		log.Println("Not exists db file")
		os.Exit(1)
	}
	db := getBlockchainDatabase()
	defer db.Close()
	err := db.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(utxoBucket))
		errorHandle(err)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			utxo := DeserializeUtxo(v)
			up.Pool = append(up.Pool, &utxo)
		}
		return nil
	})
	errorHandle(err)
	return &up
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
		err = b.Put(utxo.Hash(), utxo.Serialize())
		errorHandle(err)
		return nil
	})
	errorHandle(err)
	up.Pool = append(up.Pool, utxo)
}

func (up *UTXOPool) String() string {
	var lines []string
	for _, utxo := range up.Pool {
		lines = append(lines, utxo.String())
	}
	return strings.Join(lines, "\n")
}
