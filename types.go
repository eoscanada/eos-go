package eosapi

import (
	"bytes"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/lunixbochs/struc"
)

var base32Encoding = base32.NewEncoding(".12345abcdefghijklmnopqrstuvwxyz").WithPadding(base32.NoPadding)

type AccountName string

func (acct *AccountName) Pack(p []byte, opt *struc.Options) (int, error) {
	_, err := base32Encoding.Decode(p[:8], []byte(*acct))
	if err != nil {
		return 8, err
	}
	return 8, nil
}

func (acct *AccountName) Unpack(r io.Reader, length int, opt *struc.Options) error {
	data := make([]byte, 0, 8)
	_, err := r.Read(data[:8])
	if err != nil {
		return err
	}

	var out []byte
	base32Encoding.Encode(out, data)

	*acct = AccountName(out)
	return nil
}

func (acct *AccountName) Size(opt *struc.Options) int {
	return 8
}

func (acct AccountName) String() string {
	return string(acct)
}

type Asset struct {
	Precision int    `struc:"uint8"`
	Symbol    string `struc:"[7]byte"`
} // decode "1000.0000 EOS" as `Asset{Amount: 10000000, Symbol: "EOS", Precision: 4}`

func (a *Asset) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return nil
}

// NOT RIGHT SIGNATURE:
func (a *Asset) MarshalJSON() (data []byte, err error) {
	return nil, nil
}

type AccountResp struct {
	AccountName AccountName  `json:"account"`
	Permissions []Permission `json:"permissions"`
}

type CurrencyBalanceResp struct {
	EOSBalance        Asset    `json:"eos_balance"`
	StakedBalance     Asset    `json:"staked_balance"`
	UnstakingBalance  Asset    `json:"unstaking_balance"`
	LastUnstakingTime JSONTime `json:"last_unstaking_time"`
}

type PublicKey string

type Permission struct {
	PermName     string `json:"perm_name"`
	Parent       string `json:"parent"`
	RequiredAuth Auth   `json:"required_auth"`
}

type PermissionLevel struct {
	Account    AccountName `json:"account"`
	Permission string      `json:"permission"`
}

type Auth struct {
	Threshold int           `json:"threshold"`
	Keys      []WeightedKey `json:"keys"`
	Accounts  []AccountName `json:"accounts"`
}

type WeightedKey struct {
	Key    PublicKey `json:"key"`
	Weight int       `json:"weight"`
}

type Code struct {
	AccountName AccountName `json:"account_name"`
	CodeHash    string      `json:"code_hash"`
	WAST        string      `json:"wast"` // TODO: decode into Go ast, see https://github.com/go-interpreter/wagon
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
	Account       AccountName       `json:"account"`
	Name          string            `json:"name"`
	Authorization []PermissionLevel `json:"authorization"`
	Data          string            `json:"data"` // as HEX when we receive it.. FIXME: decode from hex directly.. and encode back plz!

	Type       string        `json:"type"`       // dawn-2
	Code       string        `json:"code"`       // dawn-2
	Recipients []AccountName `json:"recipients"` // dawn-2 ?
}

type Table struct {
	Name      string   `json:"name"`
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
	ServerVersion            string      `json:"server_version"`              // "2cc40a4e"
	HeadBlockNum             uint32      `json:"head_block_num"`              // 2465669,
	LastIrreversibleBlockNum uint32      `json:"last_irreversible_block_num"` // 2465655
	HeadBlockID              string      `json:"head_block_id"`               // "00259f856bfa142d1d60aff77e70f0c4f3eab30789e9539d2684f9f8758f1b88",
	HeadBlockTime            JSONTime    `json:"head_block_time"`             //  "2018-02-02T04:19:32"
	HeadBlockProducer        AccountName `json:"head_block_producer"`         // "inita"
	RecentSlots              string      `json:"recent_slots"`                //  "1111111111111111111111111111111111111111111111111111111111111111"
	ParticipationRate        string      `json:"participation_rate"`          // "1.00000000000000000" // this should be a `double`, or a decimal of some sort..

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
	More bool            `json:"more"`
	Rows json.RawMessage `json:"rows"` // defer loading, as it depends on `JSON` being true/false.
}

func (resp *GetTableRowsResp) JSONToStructs(v interface{}) error {
	return json.Unmarshal(resp.Rows, v)
}

func (resp *GetTableRowsResp) BinaryToStructs(v interface{}) error {
	var rows []string

	err := json.Unmarshal(resp.Rows, &rows)
	if err != nil {
		return err
	}

	for _, row := range rows {
		bin, err := hex.DecodeString(row)
		if err != nil {
			return err
		}

		fmt.Println("MAMA", bin)

		ourstruct := &MyStruct{}
		if err := struc.Unpack(bytes.NewReader(bin), ourstruct); err != nil {
			return err
		}

		spew.Dump(ourstruct)
	}

	return nil
}

type MyStruct struct {
	Key      string `struc:"[8]int8,little"`
	Balance  uint64 `struc:"uint64,little"`
	Currency string `struc:"[8]int8,little"`
}

type GetRequiredKeysResp struct {
	RequiredKeys []PublicKey `json:"required_keys"`
}

type Transaction struct { // WARN: is a `variant` in C++
	RefBlockNum    string   `json:"ref_block_num"`
	RefBlockPrefix string   `json:"ref_block_prefix"`
	Expiration     JSONTime `json:"expiration"`
	Scope          []string `json:"scope"`
	Actions        []Action `json:"actions"`
	Signatures     []string `json:"signatures"`
	Authorizations []string `json:"authorizations"`
}

type DeferredTransaction struct {
	Transaction

	SenderID   uint32      `json:"sender_id"`
	Sender     AccountName `json:"sender"`
	DelayUntil JSONTime    `json:"delay_until"`
}

type PushTransactionResp struct {
	TransactionID string `json:"transaction_id"`
	Processed     bool   `json:"processed"` // WARN: is an `fc::variant` in server..
}
