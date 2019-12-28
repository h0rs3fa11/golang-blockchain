package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(num uint64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
