package cli

import (
	"fmt"
	"log"
)

func (cli *CLI) createTransaction(from, to string, amount int) {
	wallets := NewWallets()
	wallet := wallets.Wallets[from]
	if wallet == nil {
		log.Fatal(NewGlockchainError(94001))
	}
	tx, err := NewTransaction(wallet, []byte(to), amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx.String())
}

func (cli *CLI) printTransactionPool() {
	txp, err := NewTransactionPool()
	if err != nil {
		log.Fatal(err)
	}
	for _, tx := range txp.Pool {
		fmt.Println(tx.String())
		fmt.Println()
	}
}
