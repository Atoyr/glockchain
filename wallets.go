package glockchain

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Wallets is Wallet collection
type Wallets struct {
	Wallets map[string]*Wallet
}

// NewWallets is Create Wallets
func NewWallets() *Wallets {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	wallets.LoadFromFile()
	return &wallets
}

// CreateWallet is created wallet and added Wallets collection
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	ws.Wallets[address] = wallet
	return address
}

// GetAddresses is calling wallet function with Wallets collection
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

// GetWallet is Wallet from Wallets collection
func (ws Wallets) GetWallet(address string) (*Wallet, error) {
	wallet, ok := ws.Wallets[address]
	if ok == false {
		return nil, NewGlockchainError(94001)
	}
	return wallet, nil
}

// SaveToFile is Save wallets data to file
func (ws Wallets) SaveToFile() {
	var content bytes.Buffer
	walletFile := walletFile
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// LoadFromFile is Loading wallets from file
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = wallets.Wallets
	return nil
}
