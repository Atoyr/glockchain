package core

type Blockchain struct {
	blocks []*Block
	pool   []*Transaction
}

func (bc *Blockchain) AddBlock() {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(bc.pool[:], prevBlock.Hash)
	bc.pool = make([]*Transaction, 0, 0)
	bc.blocks = append(bc.blocks, newBlock)
}
