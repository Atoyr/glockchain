package core

import (
	"fmt"
)

func (cli *CLI) createTransaction(from, to string, amount int) {
	wallets := NewWallets()
	wallet := wallets.Wallets[from]
	tx := NewTransaction(wallet, []byte(to), amount)
	fmt.Println(tx.String())
}

func (cli *CLI) printTransactionPool() {
	txp := NewTransactionPool()
	for _, tx := range txp.Pool {
		fmt.Println(tx.String())
		fmt.Println()
	}
}
