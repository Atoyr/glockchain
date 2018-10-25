package core

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	Bc *Blockchain
}

func (cli *CLI) printExecute() {
	fmt.Println("  /$$$$$$  /$$                     /$$  ")
	fmt.Println(" /$$__  $$| $$                    | $$  ")
	fmt.Println("| $$  \\__/| $$  /$$$$$$   /$$$$$$$| $$   /$$")
	fmt.Println("| $$ /$$$$| $$ /$$__  $$ /$$_____/| $$  /$$/")
	fmt.Println("| $$|_  $$| $$| $$  \\ $$| $$      | $$$$$$/ ")
	fmt.Println("| $$  \\ $$| $$| $$  | $$| $$      | $$_  $$ ")
	fmt.Println("|  $$$$$$/| $$|  $$$$$$/|  $$$$$$$| $$ \\  $$")
	fmt.Println(" \\______/ |__/ \\______/  \\_______/|__/  \\__/")
	fmt.Println("")
	fmt.Println("      /$$$$$$  /$$                 /$$      ")
	fmt.Println("     /$$__  $$| $$                |__/")
	fmt.Println("    | $$  \\__/| $$$$$$$   /$$$$$$  /$$ /$$$$$$$")
	fmt.Println("    | $$      | $$__  $$ |____  $$| $$| $$__  $$")
	fmt.Println("    | $$      | $$  \\ $$  /$$$$$$$| $$| $$  \\ $$")
	fmt.Println("    | $$    $$| $$  | $$ /$$__  $$| $$| $$  | $$")
	fmt.Println("    |  $$$$$$/| $$  | $$|  $$$$$$$| $$| $$  | $$")
	fmt.Println("     \\______/ |__/  |__/ \\_______/|__/|__/  |__/")

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
	//cli.validateArgs()
	cli.printExecute()
	createBlockchainCmd := flag.NewFlagSet("create", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("pc", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("cw", flag.ExitOnError)
	printWalletCmd := flag.NewFlagSet("pw", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address tosend genesis block reward to")
	var err error
	switch os.Args[1] {
	case "create":
		err = createBlockchainCmd.Parse(os.Args[2:])
	case "addblock":
		err = addBlockCmd.Parse(os.Args[2:])
	case "pc":
		err = printChainCmd.Parse(os.Args[2:])
	case "cw":
		err = createWalletCmd.Parse(os.Args[2:])
	case "pw":
		err = printWalletCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if err != nil {
		os.Exit(1)
	}
	if createBlockchainCmd.Parsed() {
		cli.createBlockchain(*createBlockchainAddress)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	if printWalletCmd.Parsed() {
		cli.printWallets()
	}
}

func (cli *CLI) createBlockchain(address string) {
	log.Println(address)
	a := []byte(address)
	CreateBlockchain(BytesToAddress(a))
}

func (cli *CLI) printChain() {
	cli.Bc = GetBlockchain()
	bci := cli.Bc.Iterator()
	for {
		log.Println("hoge")
		block := bci.Next()
		fmt.Printf("Hash: %x \n", block.Hash)
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()
		if len(block.PreviousHash) == 0 {
			break
		}
	}
}
func (cli *CLI) createWallet() {
	wallets := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("address: %s\n", address)
}

func (cli *CLI) printWallets() {
	wallets := NewWallets()
	for address := range wallets.Wallets {
		fmt.Printf("address: %s\n", address)
	}
}
