package core

import (
	"fmt"
)

func (cli *CLI) createTransaction(from, to string, amount int) {
	wallets := NewWallets()
	wallet := wallets.Wallets[from]
	NewTransaction(wallet, []byte(to), amount)
}

func (cli *CLI) printTransactionPool() {
	txp := GetTransactionPool()
	for _, tx := range txp.Pool {
		fmt.Println(tx.String())
		fmt.Println()
	}
}
