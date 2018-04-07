package eos

import (
	"encoding/hex"
	"encoding/json"
	"reflect"

	"github.com/eoscanada/eos-go/ecc"
)

type InfoResp struct {
	ServerVersion            string      `json:"server_version"`              // "2cc40a4e"
	HeadBlockNum             uint32      `json:"head_block_num"`              // 2465669,
	LastIrreversibleBlockNum uint32      `json:"last_irreversible_block_num"` // 2465655
	HeadBlockID              string      `json:"head_block_id"`               // "00259f856bfa142d1d60aff77e70f0c4f3eab30789e9539d2684f9f8758f1b88",
	HeadBlockTime            JSONTime    `json:"head_block_time"`             //  "2018-02-02T04:19:32"
	HeadBlockProducer        AccountName `json:"head_block_producer"`         // "inita"
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

	outSlice := reflect.ValueOf(v).Elem()
	structType := reflect.TypeOf(v).Elem().Elem()

	for _, row := range rows {
		bin, err := hex.DecodeString(row)
		if err != nil {
			return err
		}

		// access the type of the `Slice`, create a bunch of them..
		newStruct := reflect.New(structType)
		if err := UnmarshalBinary(bin, newStruct.Interface()); err != nil {
			return err
		}

		outSlice = reflect.Append(outSlice, reflect.Indirect(newStruct))
	}

	reflect.ValueOf(v).Elem().Set(outSlice)

	return nil
}

type Currency struct {
	Precision uint8
	Name      CurrencyName
}

type GetRequiredKeysResp struct {
	RequiredKeys []ecc.PublicKey `json:"required_keys"`
}

// PushTransactionFullResp unwraps the responses from a successful `push_transaction`.
// FIXME: REVIEW the actual output, things have moved here.
type PushTransactionFullResp struct {
	TransactionID string               `json:"transaction_id"`
	Processed     TransactionProcessed `json:"processed"` // WARN: is an `fc::variant` in server..
}

type TransactionProcessed struct {
	Status               string        `json:"status"`
	ID                   SHA256Bytes   `json:"id"`
	ActionTraces         []ActionTrace `json:"action_traces"`
	DeferredTransactions []string      `json:"deferred_transactions"` // that's not right... dig to find what's there..
}

type ActionTrace struct {
	Receiver AccountName `json:"receiver"`
	// Action     Action       `json:"act"` // FIXME: how do we unpack that ? what's on the other side anyway?
	Console    string       `json:"console"`
	RegionID   uint16       `json:"region_id"`
	CycleIndex int          `json:"cycle_index"`
	DataAccess []DataAccess `json:"data_access"`
}

type DataAccess struct {
	Type     string      `json:"type"` // "write", "read"?
	Code     AccountName `json:"code"`
	Scope    AccountName `json:"scope"`
	Sequence int         `json:"sequence"`
}

type PushTransactionShortResp struct {
	TransactionID string `json:"transaction_id"`
	Processed     bool   `json:"processed"` // WARN: is an `fc::variant` in server..
}

//

type WalletSignTransactionResp struct {
	// Ignore the rest of the transaction, so the wallet server
	// doesn't forge some transactions on your behalf, and you send it
	// to the network..  ... although.. it's better if you can trust
	// your wallet !

	Signatures []string `json:"signatures"`
}

type MyStruct struct {
	Currency
	Balance uint64
}

// NetConnectionResp
type NetConnectionsResp struct {
	Peer          string           `json:"peer"`
	Connecting    bool             `json:"connecting"`
	Syncing       bool             `json:"syncing"`
	LastHandshake HandshakeMessage `json:"last_handshake"`
}

type NetStatusResp struct {
}

type NetConnectResp string

type NetDisconnectResp string
