package core

import "log"

// Version is Glockchain version
const Version = byte(0x00)

// HashLength is the expected length of the hash
const HashLength = 32

// AddressLength is the expected length of the address
const AddressLength = 20

const WalletVersion = byte(0x00)

const addressChecksumLength = 4

const walletFile = "wallet.dat"

const dbFile = "glockchains.db"

func errorHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
