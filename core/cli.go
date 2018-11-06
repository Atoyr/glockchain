package core

import (
	"fmt"
	"log"

	urfaveCli "github.com/urfave/cli"
)

type CLI struct {
	App *urfaveCli.App
	Bc  *Blockchain
}

func NewCLI() *CLI {
	var c CLI
	app := urfaveCli.NewApp()
	app.Name = "GlockChain"
	app.Usage = "A golang blockchain application"
	c.App = app

	c.Initialize()

	return &c
}

func (cli *CLI) Initialize() {
	cli.App.Before = func(c *urfaveCli.Context) error {
		cli.printExecute()
		return nil
	}
	cli.App.Author = "atoyr"

	// blockchain interactive interface
	//cli.App.Action = func(c *urfaveCli.Context) error {
	//return nil
	//}
	cli.App.Commands = []urfaveCli.Command{
		{
			Name:    "initialize",
			Aliases: []string{"i", "init"},
			Usage:   "Execute createwallet and create a blockchain and send genesis block",
			Action: func(c *urfaveCli.Context) error {
				cli.initializeBlockchain()
				return nil
			},
		},
		{
			Name:    "wallet",
			Aliases: []string{"w"},
			Usage:   "wallet action",
			Subcommands: []urfaveCli.Command{
				{
					Name:  "create",
					Usage: "Generate a new key-pair and saves it into the wallet file",
					Action: func(c *urfaveCli.Context) error {
						cli.createWallet()
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "Lists all address from the wallet file",
					Action: func(c *urfaveCli.Context) error {
						cli.printWallets()
						return nil
					},
				},
				{
					Name:  "balance",
					Usage: "Get balance",
					Action: func(c *urfaveCli.Context) error {
						address := c.String("a")
						if address == "" {
							log.Println("ERROR: Address not fount")
							return nil
						}
						cli.getBalance(address)
						return nil
					},
					Flags: []urfaveCli.Flag{
						urfaveCli.StringFlag{
							Name: "address, a",
						},
					},
				},
			},
		},
		{
			Name:    "block",
			Aliases: []string{"b"},
			Usage:   "block action",
			Subcommands: []urfaveCli.Command{
				{
					Name:  "create",
					Usage: "create block",
					Action: func(c *urfaveCli.Context) error {
						fmt.Println("Not implements")
						return nil
					},
				},
			},
		},
		{
			Name:    "blockchain",
			Aliases: []string{"bc"},
			Usage:   "blockchain action",
			Subcommands: []urfaveCli.Command{
				{
					Name:  "print",
					Usage: "print blockchain",
					Action: func(c *urfaveCli.Context) error {
						cli.printChain()
						return nil
					},
				},
			},
		},
		{
			Name:    "transaction",
			Aliases: []string{"tran", "t"},
			Usage:   "transaction action",
			Subcommands: []urfaveCli.Command{
				{
					Name:  "create",
					Usage: "create transaction",
					Action: func(c *urfaveCli.Context) error {
						fmt.Println("Not implements")
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "show transaction pool",
					Action: func(c *urfaveCli.Context) error {
						fmt.Println("Not implements")
						return nil
					},
				},
			},
		},
	}
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
	fmt.Println("")
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

func (cli *CLI) printUtxo() {
	up := GetUTXOPool()
	fmt.Println(up)
}
