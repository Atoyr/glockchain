package gochain

type Transaction struct {
	Version   int
	Sender    []byte
	Recipient []byte
	Amount    int
}

func (t *Transaction) ToByte() []byte {
	bytes := make([]byte, 100)
	bytes = append(bytes, []byte(t.Version))
	bytes = append(bytes, t.Sender)
	bytes = append(bytes, t.Recipient)
	bytes = append(bytes, []byte(t.Amount))
	return bytes
}
