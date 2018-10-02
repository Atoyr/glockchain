package core

type TXInput struct {
	PrevHash  []byte
	TXOIndex  []byte
	VI        []byte
	Signature []byte
	PubKey    []byte
}
