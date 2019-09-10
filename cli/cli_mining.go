package blockchain

import (
	"fmt"
	"log"
	"time"
	"github.com/atoyr/glockchain/blockchain"
)

func (cli *CLI) mining(address string) {
	wallets := blockchain.NewWallets()
	wallet, err := wallets.GetWallet(address)
	if err != nil {
		log.Fatal(err)
	}
	t := time.Now()
	fmt.Printf("Execute mining at %d-%d-%d %d:%d:%d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	txpool, err := blockchain.NewTransactionPool()
	if err != nil {
		log.Fatal(err)
	}
	if len(txpool.Pool) == 0 {
		return
	}
	bc, tip, err := blockchain.GetBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	block, err := blockchain.NewBlock(txpool.Pool, bc, tip)
	if err != nil {
		log.Fatal(err)
	}
	bc.AddBlock(block)
	txpool.ClearTransactionPool()
	t = time.Now()
	fmt.Printf("\n\nFinished mining at %d-%d-%d %d:%d:%d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Println("Add Block ")
	fmt.Println(block)
	tx, _ := blockchain.NewCoinbaseTX(100, wallet)
	txp, _ := blockchain.NewTransactionPool()
	up, _ := blockchain.GetUTXOPool()
	txp.AddTransaction(tx)
	up.AddUTXO(tx)
}
