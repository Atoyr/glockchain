package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"os"

	"github.com/atoyr/glockchain/util"
	"golang.org/x/crypto/ripemd160"
)

// Wallet is Glockchain Wallet private and publick keyValue pair
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet is Create new wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionPayload := append([]byte{WalletVersion}, pubKeyHash...)
	checksum := checksum(versionPayload)
	fullPayload := append(versionPayload, checksum...)
	address := util.Base58Encode(fullPayload)
	return address
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

func HashPubKey(pubKey []byte) []byte {
	pubSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(pubSHA256[:])
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	pubRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return pubRIPEMD160
}
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLength]
}
