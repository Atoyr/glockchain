package core

type TXData struct {
	TXHash    Hash
	TXIndex   int
	Address   Address
	Value     int
	Signature []byte
	PubKey    []byte
}

func NewTXData(txindex int, address Address, value int) *TXData {
	var txdata TXData
	txdata.TXIndex = txindex
	txdata.Address = address
	txdata.Value = value
	return &txdata
}
