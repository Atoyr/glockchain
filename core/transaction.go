package core

import "github.com/atoyr/gochain/util"

type Transaction struct {
	Version   int
	Sender    []byte
	Recipient []byte
	Amount    int
}

func (t *Transaction) ToByte() []byte {
	bytes := make([]byte, 100)
	bytes = append(bytes, util.Interface2bytes(t.Version)...)
	bytes = append(bytes, util.Interface2bytes(t.Sender)...)
	bytes = append(bytes, util.Interface2bytes(t.Recipient)...)
	bytes = append(bytes, util.Interface2bytes(t.Amount)...)
	return bytes
}
