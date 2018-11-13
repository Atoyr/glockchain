package main

import (
	"os"

	"github.com/atoyr/glockchain/core"
)

func main() {
	cli := core.NewCLI()
	cli.App.Run(os.Args)
}
