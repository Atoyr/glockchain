package core

import (
	"fmt"
	"log"
	"time"
)

func (cli *CLI) mining() {
	t := time.Now()
	fmt.Printf("Execute mining at %d-%d-%d %d:%d:%d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	txpool, err := NewTransactionPool()
	if err != nil {
		log.Fatal(err)
	}
	if len(txpool.Pool) == 0 {
		fmt.Println(getErrorMessage(93001))
		return
	}
	bc, tip, err := GetBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	block, err := NewBlock(txpool.Pool, tip)
	if err != nil {
		log.Fatal(err)
	}
	bc.AddBlock(block)
	txpool.ClearTransactionPool()
	t = time.Now()
	fmt.Printf("\n\nFinished mining at %d-%d-%d %d:%d:%d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Println("Add Block ")
	fmt.Println(block)
}
