package main

import (
	"github.com/atoyr/glockchain/core"
)

func main() {
	bc := core.NewBlockchain()
	cli := core.CLI{bc}
	cli.Run()
}
