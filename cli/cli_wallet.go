package blockchain

import (
	"fmt"
	"log"
	"github.com/atoyr/glockchain/blockchain"
)

func (cli *CLI) createWallet() {
	wallets := blockchain.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("address: %s\n", address)
}

func (cli *CLI) printWallets() {
	wallets := blockchain.NewWallets()
	for address := range wallets.Wallets {
		fmt.Printf("address: %s\n", address)
	}
}

func (cli *CLI) getBalance(address string) {
	if !blockchain.ValidateAddress([]byte(address)) {
		log.Panic("ERROR: Address is not valid")
	}
	utxopool, err := blockchain.GetUTXOPool()
	if err != nil {
		log.Fatal(err)
	}
	pubKeyHash := blockchain.AddressToPubKeyHash([]byte(address))
	balance, _ := utxopool.FindUTXOs(pubKeyHash)
	fmt.Printf("Balance of %s : %d \n", address, balance)
}

func (cli *CLI) getAllBalance() {
	wallets := blockchain.NewWallets()
	utxopool, err := blockchain.GetUTXOPool()
	if err != nil {
		log.Fatal(err)
	}
	for address := range wallets.Wallets {
		pubKeyHash := blockchain.AddressToPubKeyHash([]byte(address))
		balance, _ := utxopool.FindUTXOs(pubKeyHash)
		fmt.Printf("Balance of %s : %d \n", address, balance)
	}
}
