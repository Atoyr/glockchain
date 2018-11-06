package core

import (
	"fmt"
	"log"

	"github.com/atoyr/glockchain/util"
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
	utxopool := GetUTXOPool()
	pubKeyHash := util.Base58Decode([]byte(address))
	balance, _ := utxopool.FindUTXOs(pubKeyHash)
	fmt.Printf("Balance of %s : %d \n", address, balance)
}
