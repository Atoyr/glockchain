package core

import (
	"fmt"
	"os"
)

func (cli *CLI) mining() {
	txpool := GetTransactionPool()
	if len(txpool.Pool) == 0 {
		fmt.Println("Not have data into TXPool")
		os.Exit(0)
	}
	bc, tip := GetBlockchain()
	block := NewBlock(txpool.Pool, tip)
	bc.AddBlock(block)
	txpool.ClearTransactionPool()
	fmt.Println("Add Block ")
	fmt.Println(block)
}
