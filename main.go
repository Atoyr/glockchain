package main

import (
	"fmt"

	"github.com/atoyr/gochain/core"
)

func main() {
	bc := core.NewBlockchain()
	for _, block := range bc.Blocks {
		fmt.Printf("data: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
	}
}
