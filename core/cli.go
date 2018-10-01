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
	fmt.Println("\\___/|_|")
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
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	var err error
	switch os.Args[1] {
	case "addblock":
		err = addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		err = printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if err != nil {
		os.Exit(1)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func (cli *CLI) printChain() {
	bci := cli.Bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Hash: %x \n", block.Hash)
		fmt.Println()
		if len(block.PreviousHash) == 0 {
			break
		}
	}
}
