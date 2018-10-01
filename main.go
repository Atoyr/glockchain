package main

import (
	"github.com/atoyr/gochain/core"
)

func main() {
	bc := core.NewBlockchain()
	cli := core.CLI{bc}
	cli.Run()
}
