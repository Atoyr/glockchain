package core

import "github.com/boltdb/bolt"

type Blockchain struct {
	Blocks []*Block
	pool   []*Transaction
}

func (bc *Blockchain) AddBlock() {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(bc.pool[:], prevBlock.Hash)
	bc.pool = make([]*Transaction, 0, 0)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open("test.db", 0600, nil)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(1))
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte())
			err = b.Put(genesis.Hash, genesis.Serialize())
			tip = genesis.Hash
			return nil
		}
	})
	return &Blockchain{[]*Block{NewGenesisBlock()}, []*Transaction{}}
}
