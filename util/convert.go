package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Interface2bytes(i interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}
