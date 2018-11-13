package util

import (
	"crypto/sha256"
)

func P2PKHToAddress(pkscript []byte) ([]byte, error) {
	p := make([]byte, 1)
	p[0] = 0x00
	pub := pkscript[3 : len(pkscript)-2]
	pf := append(p[:], pub[:]...)
	h1 := sha256.Sum256(pf)
	h2 := sha256.Sum256(h1[:])
	b := append(pf[:], h2[:4]...)
	address := Base58Encode(b)
	return address, nil
}
