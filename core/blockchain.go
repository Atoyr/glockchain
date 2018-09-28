package core

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
	return &Blockchain{[]*Block{NewGenesisBlock()}, []*Transaction{}}
}
