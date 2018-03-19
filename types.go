package eosapi

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/eosioca/eosapi/ecc"
)

// For reference:
// https://github.com/mithrilcoin-io/EosCommander/blob/master/app/src/main/java/io/mithrilcoin/eoscommander/data/remote/model/types/EosByteWriter.java

type Name string
type AccountName Name
type PermissionName Name
type ActionName Name
type TableName Name

func (acct AccountName) MarshalBinary() ([]byte, error)    { return Name(acct).MarshalBinary() }
func (acct PermissionName) MarshalBinary() ([]byte, error) { return Name(acct).MarshalBinary() }
func (acct ActionName) MarshalBinary() ([]byte, error)     { return Name(acct).MarshalBinary() }
func (acct TableName) MarshalBinary() ([]byte, error)      { return Name(acct).MarshalBinary() }
func (acct Name) MarshalBinary() ([]byte, error) {
	val, err := StringToName(string(acct))
	if err != nil {
		return nil, err
	}
	var out [8]byte
	binary.LittleEndian.PutUint64(out[:8], val)
	return out[:], nil
}

func (n *AccountName) UnmarshalBinary(data []byte) error {
	*n = AccountName(NameToString(binary.LittleEndian.Uint64(data)))
	return nil
}
func (n *Name) UnmarshalBinary(data []byte) error {
	*n = Name(NameToString(binary.LittleEndian.Uint64(data)))
	return nil
}
func (n *PermissionName) UnmarshalBinary(data []byte) error {
	*n = PermissionName(NameToString(binary.LittleEndian.Uint64(data)))
	return nil
}
func (n *ActionName) UnmarshalBinary(data []byte) error {
	*n = ActionName(NameToString(binary.LittleEndian.Uint64(data)))
	return nil
}
func (n *TableName) UnmarshalBinary(data []byte) error {
	*n = TableName(NameToString(binary.LittleEndian.Uint64(data)))
	return nil
}

func (AccountName) UnmarshalBinarySize() int    { return 8 }
func (PermissionName) UnmarshalBinarySize() int { return 8 }
func (ActionName) UnmarshalBinarySize() int     { return 8 }
func (TableName) UnmarshalBinarySize() int      { return 8 }
func (Name) UnmarshalBinarySize() int           { return 8 }

// OTHER TYPES: eosjs/src/structs.js

// CurrencyName

type CurrencyName string

func (c CurrencyName) MarshalBinary() ([]byte, error) {
	out := make([]byte, 7, 7)
	fmt.Println("AMAM", out)
	copy(out, []byte(c))
	return out, nil
}

func (c *CurrencyName) UnmarshalBinary(data []byte) error {
	*c = CurrencyName(strings.TrimRight(string(data), "\x00"))
	return nil
}
func (CurrencyName) UnmarshalBinarySize() int { return 7 }

// Asset

// NOTE: there's also ExtendedAsset which is a quantity with the attached contract (AccountName)
type Asset struct {
	Amount int64
	Symbol
}

// NOTE: there's also a new ExtendedSymbol (which includes the contract (as AccountName) on which it is)
type Symbol struct {
	Precision int
	Symbol    string
}

func (a *Asset) UnmarshalBinary(data []byte) error {
	// pick up uint64 for amount
	// then one byte for Precision, and another 7 bytes as string for currency
	return nil
}
func (a *Asset) UnmarshalJSON(data []byte) error {
	// decode "1000.0000 EOS" as `Asset{Amount: 10000000, Symbol: {Precision: 4, Symbol: "EOS"}`
	// deal with the underlying `Symbol`
	return nil
}

type Permission struct {
	PermName     string    `json:"perm_name"`
	Parent       string    `json:"parent"`
	RequiredAuth Authority `json:"required_auth"`
}

type PermissionLevel struct {
	Actor      AccountName    `json:"actor"`
	Permission PermissionName `json:"permission"`
}

type PermissionLevelWeight struct {
	Permission PermissionLevel `json:"permission"`
	Weight     uint16          `json:"weight"`
}

type Authority struct {
	Threshold uint32                  `json:"threshold"`
	Keys      []KeyWeight             `json:"keys"`
	Accounts  []PermissionLevelWeight `json:"accounts"`
}

type KeyWeight struct {
	PublicKey *ecc.PublicKey `json:"public_key"`
	Weight    uint16         `json:"weight"`
}

type Code struct {
	AccountName AccountName `json:"account_name"`
	CodeHash    string      `json:"code_hash"`
	WAST        string      `json:"wast"` // TODO: decode into Go ast, see https://github.com/go-interpreter/wagon
	ABI         ABI         `json:"abi"`
}

// JSONTime

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

func (t JSONTime) UnmarshalBinary(data []byte) error {
	t.Time = time.Unix(int64(binary.LittleEndian.Uint32(data)), 0)
	return nil
}

func (t JSONTime) MarshalBinary() ([]byte, error) {
	out := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(out, uint32(t.Unix()))
	return out, nil
}

func (t JSONTime) MarshalBinarySize() int { return 4 }

// HexBytes

type HexBytes []byte

func (t HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(t))
}

func (t *HexBytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = hex.DecodeString(s)
	return
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

type Transaction struct { // WARN: is a `variant` in C++, can be a SignedTransaction or a Transaction.
	Expiration     JSONTime `json:"expiration,omitempty"`
	Region         uint16   `json:"region"`
	RefBlockNum    uint16   `json:"ref_block_num,omitempty"`
	RefBlockPrefix uint32   `json:"ref_block_prefix,omitempty"`
	// number of 8 byte words this transaction can compress into
	PackedBandwidthWords    uint16    `json:"packed_bandwidth_words,omitempty"`
	ContextFreeCPUBandwidth uint16    `json:"context_free_cpu_bandwidth,omitempty"`
	ContextFreeActions      []*Action `json:"context_free_actions,omitempty"`
	Actions                 []*Action `json:"actions,omitempty"`
}

func (tx *Transaction) Fill(api *EOSAPI) ([]byte, error) {
	info, err := api.GetInfo()
	if err != nil {
		return nil, err
	}

	blockID, err := hex.DecodeString(info.HeadBlockID)
	if err != nil {
		return nil, fmt.Errorf("decode hex: %s", err)
	}

	tx.setRefBlock(blockID)

	fmt.Println("refblockprefix:", tx.RefBlockPrefix)
	/// TODO: configure somewhere the default time for transactions,
	/// etc.. add a `.Timeout` with that duration, default to 30
	/// seconds ?
	tx.Expiration = JSONTime{info.HeadBlockTime.Add(30 * time.Second)}

	return blockID, nil
}

func (tx *Transaction) setRefBlock(blockID []byte) {
	tx.RefBlockNum = uint16(binary.BigEndian.Uint16(blockID[2:4]))
	tx.RefBlockPrefix = binary.LittleEndian.Uint32(blockID[12:16])
}

type SignedTransaction struct {
	*Transaction

	Signatures      []string   `json:"signatures"`
	ContextFreeData []HexBytes `json:"context_free_data"`
}

func NewSignedTransaction(tx *Transaction) *SignedTransaction {
	return &SignedTransaction{
		Transaction:     tx,
		Signatures:      make([]string, 0),
		ContextFreeData: make([]HexBytes, 0),
	}
}

type DeferredTransaction struct {
	*Transaction

	SenderID   uint32      `json:"sender_id"`
	Sender     AccountName `json:"sender"`
	DelayUntil JSONTime    `json:"delay_until"`
}
