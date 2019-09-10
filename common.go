package glockchain

import "log"

// Version Glockchain version
const Version = byte(0x00)

// WalletVersion Wallet Version
const WalletVersion = byte(0x00)

const addressChecksumLength = 4

const walletFile = "wallet.dat"

const dbFile = "glockchains.db"

const blocksBucket = "glockchainbucket"
const utxoBucket = "utxobucket"
const txpoolBucket = "txpoolbucket"

func errorHandle(err error) {
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
}
