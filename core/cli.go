package core

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

type CLI struct {
	App *cli.App
	Bc  *Blockchain
}

func NewCLI() *CLI {
	var c CLI
	app := cli.NewApp()
	app.Name = "GlockChain"

	app.Before = func(c *cli.Context) error {
		printExecute()
		return nil
	}

	c.App = app
	return &c
}

func printExecute() {
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
	fmt.Println("Usage:")
	fmt.Println("  init")
	fmt.Println("    Execute createwallet and create a blockchain and send genesis block")
	fmt.Println("  wallet -create")
	fmt.Println("    Generate a new key-pair and saves it into the wallet file")
	fmt.Println("  wallet -list")
	fmt.Println("    Lists all address from the wallet file")
	fmt.Println("  wallet -balance")
	fmt.Println("    Get balance of all address")
	fmt.Println("  wallet -balance -address ADDRESS")
	fmt.Println("    Get balance of ADDRESS")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
}
func (cli *CLI) Run() {
	//cli.validateArgs()
	initializeBlockchainCmd := flag.NewFlagSet("init", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("create", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("pc", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("cw", flag.ExitOnError)
	printWalletCmd := flag.NewFlagSet("pw", flag.ExitOnError)
	printUtxoCmd := flag.NewFlagSet("pu", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address tosend genesis block reward to")
	printUtxoTxhash := printUtxoCmd.String("txhash", "", "The address tosend genesis block reward to")
	index := printUtxoCmd.Int("index", 0, "The address tosend genesis block reward to")
	var err error
	switch os.Args[1] {
	case "init":
		err = initializeBlockchainCmd.Parse(os.Args[2:])
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
	case "pu":
		err = printUtxoCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if err != nil {
		os.Exit(1)
	}
	if initializeBlockchainCmd.Parsed() {
		cli.initializeBlockchain()
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
	if printUtxoCmd.Parsed() {
		if *printUtxoTxhash != "" {
			up := GetUTXOPool()
			b, _ := hex.DecodeString(*printUtxoTxhash)
			fmt.Println(up.GetUTXO(b, *index))
		} else {
			cli.printUtxo()
		}
	}
}

func (cli *CLI) initializeBlockchain() {
	wallets := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	a := []byte(address)
	CreateBlockchain(a)
	fmt.Printf("address: %s\n", address)
	cli.printChain()
}
func (cli *CLI) createBlockchain(address string) {
	a := []byte(address)
	CreateBlockchain(a)
}

func (cli *CLI) printChain() {
	cli.Bc = GetBlockchain()
	bci := cli.Bc.Iterator()
	for {
		block := bci.Next()
		fmt.Println(block)
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

func (cli *CLI) printUtxo() {
	up := GetUTXOPool()
	fmt.Println(up)
}
