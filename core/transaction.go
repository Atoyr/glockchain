package core

import "github.com/atoyr/glockchain/util"

type Transaction struct {
	Version int
	ID      []byte
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
