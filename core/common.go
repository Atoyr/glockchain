package core

import "log"

// Version is Glockchain version
const Version = byte(0x00)

const WalletVersion = byte(0x00)

const addressChecksumLength = 4

const walletFile = "wallet.dat"

const dbFile = "glockchains.db"

func errorHandle(err error) {
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
}
