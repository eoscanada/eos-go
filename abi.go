package eosapi

// see: libraries/chain/contracts/abi_serializer.cpp:53...
type ABI struct {
	Types   []ABIType   `json:"types"`
	Structs []StructDef `json:"structs"`
	Actions []Action    `json:"actions"`
	Tables  []Table     `json:"tables"`
}

type ABIType struct {
	NewTypeName string `json:"new_type_name"`
	Type        string `json:"type"`
}

type StructDef struct {
	Name   string            `json:"name"`
	Base   string            `json:"base"`
	Fields map[string]string `json:"fields"` // WARN: UNORDERED!!! Should use `https://github.com/virtuald/go-ordered-json/blob/master/example_test.go`
}

type Table struct {
	Name      string   `json:"name"`
	IndexType string   `json:"index_type"`
	KeyNames  []string `json:"key_names"`
	KeyTypes  []string `json:"key_types"`
	Type      string   `json:"type"`
}
