package eosapi

import (
	"encoding/json"
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
	WAST        string      `json:"wast"` // TODO: decode into Go ast, see https://github.com/go-interpreter/wagon
	ABI         ABI         `json:"abi"`  // TODO: decode ABI into structs too
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

type BlockResp struct {
	Previous              string           `json:"previous"`                // : "0000007a9dde66f1666089891e316ac4cb0c47af427ae97f93f36a4f1159a194",
	Timestamp             JSONTime         `json:"timestamp"`               // : "2017-12-04T17:12:08",
	TransactionMerkleRoot string           `json:"transaction_merkle_root"` // : "0000000000000000000000000000000000000000000000000000000000000000",
	Producer              AccountName      `json:"producer"`                // : "initj",
	ProducerChanges       []ProducerChange `json:"producer_changes"`        // : [],
	ProducerSignature     string           `json:"producer_signature"`      // : "203dbf00b0968bfc47a8b749bbfdb91f8362b27c3e148a8a3c2e92f42ec55e9baa45d526412c8a2fc0dd35b484e4262e734bea49000c6f9c8dbac3d8861c1386c0",
	Cycles                []Cycle          `json:"cycles"`                  // : [],
	ID                    string           `json:"id"`                      // : "0000007b677719bdd76d729c3ac36bed5790d5548aadc26804489e5e179f4a5b",
	BlockNum              uint64           `json:"block_num"`               // : 123,
	RefBlockPrefix        uint64           `json:"ref_block_prefix"`        // : 2624744919

}

type ProducerChange struct {
}

type Cycle struct {
}

type GetTableRowsRequest struct {
	JSON       bool   `json:"json"`
	Scope      string `json:"scope"`
	Code       string `json:"code"`
	Table      string `json:"table"`
	TableKey   string `json:"table_key"`
	LowerBound string `json:"lower_bound"`
	UpperBound string `json:"upper_bount"`
	Limit      uint32 `json:"limit,omitempty"` // defaults to 10 => chain_plugin.hpp:struct get_table_rows_params
}

type GetTableRowsResp struct {
	More bool              `json:"more"`
	Rows []json.RawMessage `json:"rows"` // defer loading, as it depends on `JSON` being true/false.
}
