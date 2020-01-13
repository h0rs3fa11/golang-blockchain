package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func IntToHex(num uint64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// 将json转array
func JSONToArray(jsonStr string) []string {

	var array []string
	if err := json.Unmarshal([]byte(jsonStr), &array); err != nil {
		fmt.Println("传入参数不是表示的JSON数组格式...")
		os.Exit(1)
	}

	return array
}
