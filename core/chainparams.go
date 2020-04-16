package core

import (
	"fmt"
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

var configpath = "chainparam.json"

// TODO 配置文件
type Chainparams struct {
	TargetBits int
	Subsidy    int
	Fee        int
	Miner	string
}

type JSONStruct struct {

}

func GetConfig() (*Chainparams) {
	params := Chainparams{}

	json := JSONStruct{}

	json.loadJSONFile(configpath, &params)

	return &params
}

func (params *Chainparams)Updateparams(key string, value string) error {
	var result map[string]string
	json := JSONStruct{}

	json.loadJSONFile(configpath, &result)

	if _, ok := result[key]; ok {
		result[key] = value
	} else {
		return &blockchainError{fmt.Sprintf("Config file hava not %s keyword", key)}
	}

	byteValue, err := jsoniter.Marshal(result)
	if err != nil {
		return &blockchainError{"Json marshal failed"}
	}

	err = ioutil.WriteFile(configpath, byteValue, 0644)
	params = GetConfig()
	
	return err
}

func (json *JSONStruct)loadJSONFile(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Open file error")
	}

	err = jsoniter.Unmarshal(data, v)
	if err != nil {
		fmt.Println("Parsing json file failed")
	}
}

func (params *Chainparams) setCoinbase() {
	wallets, _ := NewWallets()
	for address, _ := range wallets.WalletsMap {
		params.Miner = address
		break
	}
}