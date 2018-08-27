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
	Name   string  `json:"name"`
	Base   string  `json:"base"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *Struct) BaseFields() []Field {

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

func (a *ABI) ActionForName(name string) *Action {
	for _, a := range a.Actions {
		if a.Name == name {
			return &a
		}
	}
	return nil
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

	action := abi.ActionForName(actionName)
	if action == nil {
		return fmt.Errorf("action %s not found in abi", actionName)
	}

	//toto append base and fields from action
	//todo get type,
	//search in type and and change type name if found.
	//search in struct and loop on fields
	//read value from data

	return nil
}
