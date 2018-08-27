package abi

import (
	"encoding/json"
	"fmt"
	"io"
)

type Type struct {
	NewTypeName string `json:"new_type_name"`
	Type        string `json:"type"`
}

type Struct struct {
	Name   string `json:"name"`
	Base   string `json:"base"`
	Fields []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"fields"`
}

type Action struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	RicardianContract string `json:"ricardian_contract"`
}

type Table struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	IndexType string   `json:"index_type"`
	KeyNames  []string `json:"key_names"`
	KeyTypes  []string `json:"key_types"`
}

type ABI struct {
	Version          string        `json:"version"`
	Types            []Type        `json:"types"`
	Structs          []Struct      `json:"structs"`
	Actions          []Action      `json:"actions"`
	Tables           []Table       `json:"tables"`
	RicardianClauses []interface{} `json:"ricardian_clauses"`
	AbiExtensions    []interface{} `json:"abi_extensions"`
}

type ABIDecoder struct {
	data      []byte
	abiReader io.Reader
	pos       int
}

func NewABIDecoder(data []byte, abiReader io.Reader) *ABIDecoder {

	return &ABIDecoder{
		data:      data,
		abiReader: abiReader,
	}
}

func (d *ABIDecoder) Decode(result map[string]interface{}, actionName string) error {

	abi := &ABI{}
	abiDecoder := json.NewDecoder(d.abiReader)
	err := abiDecoder.Decode(abi)
	if err != nil {
		return fmt.Errorf("read abi: %s", err)
	}

	return nil
}
