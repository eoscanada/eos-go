package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	eos "github.com/eoscanada/eos-go"
)

func main() {
	abi, err := eos.NewABI(strings.NewReader(abiJSON()))
	if err != nil {
		panic(fmt.Errorf("get ABI: %s", err))
	}

	tableDef := abi.TableForName(eos.TableName("activebets"))
	if tableDef == nil {
		panic(fmt.Errorf("table be should be present"))
	}

	bytes, err := abi.DecodeTableRowTyped(tableDef.Type, data())
	if err != nil {
		panic(fmt.Errorf("decode row: %s", err))
	}

	fmt.Println(string(bytes))
}

func data() []byte {
	bytes, err := hex.DecodeString(`1358285f09db6dc0`)
	if err != nil {
		panic(fmt.Errorf("decode data: %s", err))
	}

	return bytes
}

func abiJSON() string {
	return `{
			"structs": [
				{
					"name": "bet",
					"fields": [
						{ "name": "id", "type": "uint64" }
					]
				}
			],
			"actions": [],
			"tables": [
				{
					"name": "activebets",
					"type": "bet"
				}
			]
	}`
}
