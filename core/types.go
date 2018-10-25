package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

//type Hash [HashLength]byte
type Address [AddressLength]byte
type TXOutputs []TXOutput

// func BytesToHash(b []byte) Hash {
// 	var h Hash
// 	h.SetBytes(b)
// 	return h
// }
// func (h Hash) Bytes() []byte  { return h[:] }
// func (h Hash) String() string { return fmt.Sprintf("%s", h) }
// func (h Hash) SetBytes(b []byte) {
// 	if len(b) > len(h) {
// 		b = b[len(b)-HashLength:]
// 	}
// 	copy(h[HashLength-len(b):], b)
// }

func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	log.Printf("a %x", a)
	return a
}
func (a Address) Bytes() []byte { return a[:] }
func (a *Address) SetBytes(b []byte) {
	if len(b) > AddressLength {
		b = b[len(b)-AddressLength:]
	}
	copy(a[:], b[:])
	log.Printf("a %x", a)
	log.Printf("b %x", b)
}

func (outs TXOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(buffer)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func DeserializeTXOutputs(data []byte) TXOutputs {
	var outputs TXOutputs
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	return outputs
}
