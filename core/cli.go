package core

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct {
	Bc *Blockchain
}

func (cli *CLI) printExecute() {
	fmt.Println("   ____")
	fmt.Println("  / __ \\")
	fmt.Println(" / / __")
	fmt.Println("/ /_/  |")
	fmt.Println("\\__/|_|")
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
}
func (cli *CLI) Run() {
	cli.validateArgs()
	cli.printExecute()
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			os.Exit(1)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

}
