package eos

// see: libraries/chain/contracts/abi_serializer.cpp:53...
// see: libraries/chain/include/eosio/chain/contracts/types.hpp:100
type ABI struct {
	Version          string            `json:"version"`
	Types            []ABIType         `json:"types,omitempty"`
	Structs          []StructDef       `json:"structs,omitempty"`
	Actions          []ActionDef       `json:"actions,omitempty"`
	Tables           []TableDef        `json:"tables,omitempty"`
	RicardianClauses []ClausePair      `json:"ricardian_clauses,omitempty"`
	ErrorMessages    []ABIErrorMessage `json:"error_messages,omitempty"`
	Extensions       []*Extension      `json:"abi_extensions,omitempty"`
}

type ABIType struct {
	NewTypeName string `json:"new_type_name"`
	Type        string `json:"type"`
}

type StructDef struct {
	Name   string     `json:"name"`
	Base   string     `json:"base"`
	Fields []FieldDef `json:"fields,omitempty"`
}

type FieldDef struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ActionDef struct {
	Name              ActionName `json:"name"`
	Type              string     `json:"type"`
	RicardianContract string     `json:"ricardian_contract"`
}

// TableDef defines a table. See libraries/chain/include/eosio/chain/contracts/types.hpp:78
type TableDef struct {
	Name      TableName `json:"name"`
	IndexType string    `json:"index_type"`
	KeyNames  []string  `json:"key_names,omitempty"`
	KeyTypes  []string  `json:"key_types,omitempty"`
	Type      string    `json:"type"`
}

// ClausePair represents clauses, related to Ricardian Contracts.
type ClausePair struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type ABIErrorMessage struct {
	Code    uint64 `json:"error_code"`
	Message string `json:"error_msg"`
}
