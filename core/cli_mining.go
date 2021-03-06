package core

import (
	"fmt"
	"os"
	"time"
)

func (cli *CLI) mining() {
	t := time.Now()
	fmt.Printf("Execute mining at %d-%d-%d %d:%d:%d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	txpool := NewTransactionPool()
	if len(txpool.Pool) == 0 {
		fmt.Println("Not have data into TXPool")
		os.Exit(0)
	}
	bc, tip := GetBlockchain()
	block := NewBlock(txpool.Pool, tip)
	bc.AddBlock(block)
	txpool.ClearTransactionPool()
	t = time.Now()
	fmt.Printf("\n\nFinished mining at %d-%d-%d %d:%d:%d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Println("Add Block ")
	fmt.Println(block)
}
