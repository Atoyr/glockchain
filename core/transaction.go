package core

import "github.com/atoyr/gochain/util"

type Transaction struct {
	Version int
	ID      []byte
	Input   []TXInput
	Output  []TXOutput
}

func (t *Transaction) ToByte() []byte {
	bytes := make([]byte, 100)
	bytes = append(bytes, util.Interface2bytes(t.Version)...)
	return bytes
}

type TXOutput struct {
	value        int
	ScripuPubKey string
}

type TXInput struct {
	TXID []byte
}
