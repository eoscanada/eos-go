package eosapi

type AccountName string
type Asset string

type AccountResp struct {
	AccountName      AccountName `json:"account_name"`
	EOSBalance       Asset       `json:"eos_balance"`
	StakedBalance    Asset       `json:"staked_balance"`
	UnstakingBalance Asset       `json:"unstaking_balance"`
	//LastUnstakingTime time.Time    `json:"last_unstaking_time"`
	// use a wrapping time, always UTC..
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	PermName     string `json:"perm_name"`
	Parent       string `json:"parent"`
	RequiredAuth Auth   `json:"required_auth"`
}

type Auth struct {
	Threshold int           `json:"threshold"`
	Keys      []WeightedKey `json:"keys"`
	Accounts  []AccountName `json:"accounts"`
}

type WeightedKey struct {
	Key    string `json:"key"`
	Weight int    `json:"weight"`
}

type Contract struct {
	AccountName AccountName `json:"account_name"`
	CodeHash    string      `json:"code_hash"`
	WAST        string      `json:"wast"`
	ABI         ABI         `json:"abi"`
}

type ABI struct {
	Types   []ABIType   `json:"types"`
	Structs []ABIStruct `json:"structs"`
	Actions []Action    `json:"actions"`
	Tables  []Table     `json:"tables"`
}

type ABIType struct {
	NewTypeName string `json:"new_type_name"`
	Type        string `json:"type"`
}

type ABIStruct struct {
	Name   string            `json:"name"`
	Base   string            `json:"base"`
	Fields map[string]string `json:"fields"` // WARN: UNORDERED!!! Should use `https://github.com/virtuald/go-ordered-json/blob/master/example_test.go`
}

type Action struct {
	ActionName string `json:"action_name"`
	Type       string `json:"type"`
}

type Table struct {
	TableName string   `json:"table_name"`
	IndexType string   `json:"index_type"`
	KeyNames  []string `json:"key_names"`
	KeyTypes  []string `json:"key_types"`
	Type      string   `json:"type"`
}
