package eos_test

import (
	"encoding/hex"
	"fmt"
	"strings"

	eos "github.com/eoscanada/eos-go"
)

func ExampleABI_DecodeTableRowTyped() {
	abi, err := eos.NewABI(strings.NewReader(abiJSON()))
	if err != nil {
		panic(fmt.Errorf("get ABI: %s", err))
	}

	tableDef := abi.TableForName(eos.TableName("votes"))
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
	bytes, err := hex.DecodeString(`424e544441505000`)
	if err != nil {
		panic(fmt.Errorf("decode data: %s", err))
	}

	return bytes
}

func abiJSON() string {
	return `{
			"structs": [
				{
					"name": "vote_t",
					"fields": [
						{ "name": "symbl", "type": "symbol_code" }
					]
				}
			],
			"actions": [],
			"tables": [
				{
					"name": "votes",
					"type": "vote_t"
				}
			]
	}`
}
