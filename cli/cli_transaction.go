package blockchain

import (
	"fmt"
	"github.com/atoyr/glockchain/blockchain"
	"log"
)

func (cli *CLI) createTransaction(from, to string, amount int) {
	wallets := blockchain.NewWallets()
	wallet := wallets.Wallets[from]
	if wallet == nil {
		log.Fatal(blockchain.NewGlockchainError(94001))
	}
	tx, err := blockchain.NewTransaction(wallet, []byte(to), amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx.String())
}

func (cli *CLI) printTransactionPool() {
	txp, err := blockchain.NewTransactionPool()
	if err != nil {
		log.Fatal(err)
	}
	for _, tx := range txp.Pool {
		fmt.Println(tx.String())
		fmt.Println()
	}
}
