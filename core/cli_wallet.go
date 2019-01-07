package core

import (
	"fmt"
	"log"
)

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

func (cli *CLI) getBalance(address string) {
	if !ValidateAddress([]byte(address)) {
		log.Panic("ERROR: Address is not valid")
	}
	utxopool, err := GetUTXOPool()
	if err != nil {
		log.Fatal(err)
	}
	pubKeyHash := AddressToPubKeyHash([]byte(address))
	balance, _ := utxopool.FindUTXOs(pubKeyHash)
	fmt.Printf("Balance of %s : %d \n", address, balance)
}

func (cli *CLI) getAllBalance() {
	wallets := NewWallets()
	utxopool, err := GetUTXOPool()
	if err != nil {
		log.Fatal(err)
	}
	for address := range wallets.Wallets {
		pubKeyHash := AddressToPubKeyHash([]byte(address))
		balance, _ := utxopool.FindUTXOs(pubKeyHash)
		fmt.Printf("Balance of %s : %d \n", address, balance)
	}
}
