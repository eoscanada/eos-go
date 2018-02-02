package eosapi

import (
	"fmt"
	"time"
)

type AccountName string
type Asset string // make it a struct

type AssetNG struct {
	Amount    uint64
	Symbol    string
	Precision int
} // decode "1000.0000 EOS" as `Asset{Amount: 10000000, Symbol: "EOS", Precision: 4}`

type AccountResp struct {
	AccountName       AccountName  `json:"account_name"`
	EOSBalance        Asset        `json:"eos_balance"`
	StakedBalance     Asset        `json:"staked_balance"`
	UnstakingBalance  Asset        `json:"unstaking_balance"`
	LastUnstakingTime JSONTime     `json:"last_unstaking_time"`
	Permissions       []Permission `json:"permissions"`
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

type JSONTime struct {
	time.Time
}

const JSONTimeFormat = "2006-01-02T15:04:05"

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", t.Format(JSONTimeFormat))), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}

	t.Time, err = time.Parse(`"`+JSONTimeFormat+`"`, string(data))
	return err
}

type InfoResp struct {
	ServerVersion            string   `json:"server_version"`              // "2cc40a4e"
	HeadBlockRun             int64    `json:"head_block_num"`              // 2465669,
	LastIrreversibleBlockNum int64    `json:"last_irreversible_block_num"` // 2465655
	HeadBlockID              string   `json:"head_block_id"`               // "00259f856bfa142d1d60aff77e70f0c4f3eab30789e9539d2684f9f8758f1b88",
	HeadBlockTime            JSONTime `json:"head_block_time"`             //  "2018-02-02T04:19:32"
	HeadBlockProducer        string   `json:"head_block_producer"`         // "inita"
	RecentSlots              string   `json:"recent_slots"`                //  "1111111111111111111111111111111111111111111111111111111111111111"
	ParticipationRate        string   `json:"participation_rate"`          // "1.00000000000000000"

}
