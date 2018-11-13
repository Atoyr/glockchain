package util

import (
	"bytes"
	"math/big"
)

var base58Char = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(base58Char)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, base58Char[mod.Int64()])
	}
	if input[0] == 0x00 {
		result = append(result, base58Char[0])
	}
	ReverseBytes(result)
	return result
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	for _, b := range input {
		charIndex := bytes.IndexByte(base58Char, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decode := result.Bytes()
	if input[0] == base58Char[0] {
		decode = append([]byte{0x00}, decode...)
	}
	return decode
}

func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
