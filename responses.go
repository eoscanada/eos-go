package eos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/eoscanada/eos-go/ecc"
)

/*
{
  "server_version": "f537bc50",
  "head_block_num": 9,
  "last_irreversible_block_num": 8,
  "last_irreversible_block_id": "00000008f98f0580d7efe7abc60abaaf8a865c9428a4267df30ff7d1937a1084",
  "head_block_id": "00000009ecd0e9fb5719431f4b86f5c9ca1887f6b6f73e5a301aaff740fd6bd3",
  "head_block_time": "2018-05-19T07:47:31",
  "head_block_producer": "eosio",
  "virtual_block_cpu_limit": 100800,
  "virtual_block_net_limit": 1056996,
  "block_cpu_limit": 99900,
  "block_net_limit": 1048576
}

*/

type InfoResp struct {
	ServerVersion            string         `json:"server_version"` // "2cc40a4e"
	ChainID                  Checksum256    `json:"chain_id"`
	HeadBlockNum             uint32         `json:"head_block_num"`              // 2465669,
	LastIrreversibleBlockNum uint32         `json:"last_irreversible_block_num"` // 2465655
	LastIrreversibleBlockID  Checksum256    `json:"last_irreversible_block_id"`  // "00000008f98f0580d7efe7abc60abaaf8a865c9428a4267df30ff7d1937a1084"
	HeadBlockID              Checksum256    `json:"head_block_id"`               // "00259f856bfa142d1d60aff77e70f0c4f3eab30789e9539d2684f9f8758f1b88",
	HeadBlockTime            BlockTimestamp `json:"head_block_time"`             //  "2018-02-02T04:19:32"
	HeadBlockProducer        AccountName    `json:"head_block_producer"`         // "inita"

	VirtualBlockCPULimit Int64  `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit Int64  `json:"virtual_block_net_limit"`
	BlockCPULimit        Int64  `json:"block_cpu_limit"`
	BlockNetLimit        Int64  `json:"block_net_limit"`
	ServerVersionString  string `json:"server_version_string"`
}

type BlockResp struct {
	SignedBlock
	ID             Checksum256 `json:"id"`
	BlockNum       uint32      `json:"block_num"`
	RefBlockPrefix uint32      `json:"ref_block_prefix"`
}

type ScheduledTransactionsResp struct {
	Transactions []ScheduledTransaction `json:"transactions"`
	More         string                 `json:"more"`
}

type DBSizeResp struct {
	FreeBytes Int64 `json:"free_bytes"`
	UsedBytes Int64 `json:"used_bytes"`
	Size      Int64 `json:"size"`
	Indices   []struct {
		Index    string `json:"index"`
		RowCount Int64  `json:"row_count"`
	} `json:"indices"`
}

type TransactionResp struct {
	ID      Checksum256 `json:"id"`
	Receipt struct {
		Status            TransactionStatus `json:"status"`
		CPUUsageMicrosec  int               `json:"cpu_usage_us"`
		NetUsageWords     int               `json:"net_usage_words"`
		PackedTransaction TransactionWithID `json:"trx"`
	} `json:"receipt"`
	Transaction           ProcessedTransaction `json:"trx"`
	BlockTime             BlockTimestamp       `json:"block_time"`
	BlockNum              uint32               `json:"block_num"`
	LastIrreversibleBlock uint32               `json:"last_irreversible_block"`
	Traces                []ActionTrace        `json:"traces"`
}

type ProcessedTransaction struct {
	Transaction SignedTransaction `json:"trx"`
}

type ActionTraceReceipt struct {
	Receiver        AccountName                    `json:"receiver"`
	ActionDigest    Checksum256                    `json:"act_digest"`
	GlobalSequence  Uint64                         `json:"global_sequence"`
	ReceiveSequence Uint64                         `json:"recv_sequence"`
	AuthSequence    []TransactionTraceAuthSequence `json:"auth_sequence"` // [["account", sequence], ["account", sequence]]
	CodeSequence    Varuint32                      `json:"code_sequence"`
	ABISequence     Varuint32                      `json:"abi_sequence"`
}

type ActionTrace struct {
	ActionOrdinal                          Varuint32           `json:"action_ordinal"`
	CreatorActionOrdinal                   Varuint32           `json:"creator_action_ordinal"`
	ClosestUnnotifiedAncestorActionOrdinal Varuint32           `json:"closest_unnotified_ancestor_action_ordinal"`
	Receipt                                *ActionTraceReceipt `json:"receipt,omitempty" eos:"optional"`
	Receiver                               AccountName         `json:"receiver"`
	Action                                 *Action             `json:"act"`
	ContextFree                            bool                `json:"context_free"`
	Elapsed                                Int64               `json:"elapsed"`
	Console                                ConsoleLog          `json:"console"`
	TransactionID                          Checksum256         `json:"trx_id"`
	BlockNum                               uint32              `json:"block_num"`
	BlockTime                              BlockTimestamp      `json:"block_time"`
	ProducerBlockID                        Checksum256         `json:"producer_block_id" eos:"optional"`
	AccountRAMDeltas                       []*AccountRAMDelta  `json:"account_ram_deltas"`
	Except                                 *Except             `json:"except,omitempty" eos:"optional"`
	ErrorCode                              *Uint64             `json:"error_code,omitempty" eos:"optional"`

	// Not present in EOSIO >= 1.8.x
	InlineTraces []ActionTrace `json:"inline_traces,omitempty" eos:"-"`
}

type AccountRAMDelta struct {
	Account AccountName `json:"account"`
	Delta   Int64       `json:"delta"`
}

type TransactionTraceAuthSequence struct {
	Account  AccountName
	Sequence Uint64
}

// [ ["account", 123123], ["account2", 345] ]
func (auth *TransactionTraceAuthSequence) UnmarshalJSON(data []byte) error {
	var ins []interface{}
	if err := json.Unmarshal(data, &ins); err != nil {
		return err
	}

	if len(ins) != 2 {
		return fmt.Errorf("expected 2 items, received %d", len(ins))
	}

	account, ok := ins[0].(string)
	if !ok {
		return fmt.Errorf("expected 1st item to be a string (account name)")
	}

	var seq Uint64
	switch el := ins[1].(type) {
	case float64:
		seq = Uint64(el)
	case string:
		seqInt, err := strconv.ParseUint(el, 10, 64)
		if err != nil {
			return fmt.Errorf("decoding auth_sequence as string: %s", err)
		}

		seq = Uint64(seqInt)
	default:

		return fmt.Errorf("expected 2nd item of auth_sequence to be a sequence number (float64 or string)")
	}

	*auth = TransactionTraceAuthSequence{AccountName(account), seq}

	return nil
}

func (auth TransactionTraceAuthSequence) MarshalJSON() (data []byte, err error) {
	return json.Marshal([]interface{}{auth.Account, auth.Sequence})
}

type SequencedTransactionResp struct {
	SeqNum int `json:"seq_num"`
	TransactionResp
}

type TransactionsResp struct {
	Transactions []SequencedTransactionResp
}

type ProducerChange struct {
}

type AccountResp struct {
	AccountName            AccountName          `json:"account_name"`
	Privileged             bool                 `json:"privileged"`
	LastCodeUpdate         JSONTime             `json:"last_code_update"`
	Created                JSONTime             `json:"created"`
	CoreLiquidBalance      Asset                `json:"core_liquid_balance"`
	RAMQuota               Int64                `json:"ram_quota"`
	RAMUsage               Int64                `json:"ram_usage"`
	NetWeight              Int64                `json:"net_weight"`
	CPUWeight              Int64                `json:"cpu_weight"`
	NetLimit               AccountResourceLimit `json:"net_limit"`
	CPULimit               AccountResourceLimit `json:"cpu_limit"`
	Permissions            []Permission         `json:"permissions"`
	TotalResources         TotalResources       `json:"total_resources"`
	SelfDelegatedBandwidth DelegatedBandwidth   `json:"self_delegated_bandwidth"`
	RefundRequest          *RefundRequest       `json:"refund_request"`
	VoterInfo              VoterInfo            `json:"voter_info"`
}

type CurrencyBalanceResp struct {
	EOSBalance        Asset    `json:"eos_balance"`
	StakedBalance     Asset    `json:"staked_balance"`
	UnstakingBalance  Asset    `json:"unstaking_balance"`
	LastUnstakingTime JSONTime `json:"last_unstaking_time"`
}

type GetTableByScopeRequest struct {
	Code       string `json:"code"`
	Table      string `json:"table"`
	LowerBound string `json:"lower_bound,omitempty"`
	UpperBound string `json:"upper_bound,omitempty"`
	Limit      uint32 `json:"limit,omitempty"`
}

type GetTableByScopeResp struct {
	More string          `json:"more"`
	Rows json.RawMessage `json:"rows"`
}

type GetTableRowsRequest struct {
	Code       string `json:"code"` // Contract "code" account where table lives
	Scope      string `json:"scope"`
	Table      string `json:"table"`
	LowerBound string `json:"lower_bound,omitempty"`
	UpperBound string `json:"upper_bound,omitempty"`
	Limit      uint32 `json:"limit,omitempty"`          // defaults to 10 => chain_plugin.hpp:struct get_table_rows_params
	KeyType    string `json:"key_type,omitempty"`       // The key type of --index, primary only supports (i64), all others support (i64, i128, i256, float64, float128, ripemd160, sha256). Special type 'name' indicates an account name.
	Index      string `json:"index_position,omitempty"` // Index number, 1 - primary (first), 2 - secondary index (in order defined by multi_index), 3 - third index, etc. Number or name of index can be specified, e.g. 'secondary' or '2'.
	EncodeType string `json:"encode_type,omitempty"`    // The encoding type of key_type (i64 , i128 , float64, float128) only support decimal encoding e.g. 'dec'" "i256 - supports both 'dec' and 'hex', ripemd160 and sha256 is 'hex' only
	Reverse    bool   `json:"reverse,omitempty"`        // Get rows in reverse of the index
	JSON       bool   `json:"json"`                     // JSON expectOutput if true, binary if false
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

		decoder := NewDecoder(bin)
		if err := decoder.Decode(newStruct.Interface()); err != nil {
			return err
		}

		outSlice = reflect.Append(outSlice, reflect.Indirect(newStruct))
	}

	reflect.ValueOf(v).Elem().Set(outSlice)

	return nil
}

type CreateSnapshotResp struct {
	SnapshotName string `json:"snapshot_name"`
	HeadBlockID  string `json:"head_block_id"`
}

type GetIntegrityHashResp struct {
	HeadBlockID  string `json:"head_block_id"`
	SnapshotName string `json:"integrity_hash"`
}

type Currency struct {
	Precision uint8
	Name      CurrencyName
}

type GetRawABIRequest struct {
	AccountName string      `json:"account_name"`
	ABIHash     Checksum256 `json:"abi_hash,omitempty"`
}

type GetRawABIResp struct {
	AccountName string      `json:"account_name"`
	CodeHash    Checksum256 `json:"code_hash"`
	ABIHash     Checksum256 `json:"abi_hash"`
	ABI         Blob        `json:"abi"`
}

type GetRequiredKeysResp struct {
	RequiredKeys []ecc.PublicKey `json:"required_keys"`
}

// PushTransactionFullResp unwraps the responses from a successful `push_transaction`.
// FIXME: REVIEW the actual expectOutput, things have moved here.
type PushTransactionFullResp struct {
	StatusCode    string
	TransactionID string               `json:"transaction_id"`
	Processed     TransactionProcessed `json:"processed"` // WARN: is an `fc::variant` in server..
	BlockID       string               `json:"block_id"`
	BlockNum      uint32               `json:"block_num"`
}

type TransactionProcessed struct {
	Status               string      `json:"status"`
	ID                   Checksum256 `json:"id"`
	ActionTraces         []Trace     `json:"action_traces"`
	DeferredTransactions []string    `json:"deferred_transactions"` // that's not right... dig to find what's there..
}

type Trace struct {
	Receiver AccountName `json:"receiver"`
	// Action     Action       `json:"act"` // FIXME: how do we unpack that ? what's on the other side anyway?
	Console    ConsoleLog   `json:"console"`
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

	Signatures []ecc.Signature `json:"signatures"`
}

type MyStruct struct {
	Currency
	Balance Uint64
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

type Global struct {
	MaxBlockNetUsage               int     `json:"max_block_net_usage"`
	TargetBlockNetUsagePct         int     `json:"target_block_net_usage_pct"`
	MaxTransactionNetUsage         int     `json:"max_transaction_net_usage"`
	BasePerTransactionNetUsage     int     `json:"base_per_transaction_net_usage"`
	NetUsageLeeway                 int     `json:"net_usage_leeway"`
	ContextFreeDiscountNetUsageNum int     `json:"context_free_discount_net_usage_num"`
	ContextFreeDiscountNetUsageDen int     `json:"context_free_discount_net_usage_den"`
	MaxBlockCPUUsage               int     `json:"max_block_cpu_usage"`
	TargetBlockCPUUsagePct         int     `json:"target_block_cpu_usage_pct"`
	MaxTransactionCPUUsage         int     `json:"max_transaction_cpu_usage"`
	MinTransactionCPUUsage         int     `json:"min_transaction_cpu_usage"`
	MaxTransactionLifetime         int     `json:"max_transaction_lifetime"`
	DeferredTrxExpirationWindow    int     `json:"deferred_trx_expiration_window"`
	MaxTransactionDelay            int     `json:"max_transaction_delay"`
	MaxInlineActionSize            int     `json:"max_inline_action_size"`
	MaxInlineActionDepth           int     `json:"max_inline_action_depth"`
	MaxAuthorityDepth              int     `json:"max_authority_depth"`
	MaxRAMSize                     string  `json:"max_ram_size"`
	TotalRAMBytesReserved          Int64   `json:"total_ram_bytes_reserved"`
	TotalRAMStake                  Int64   `json:"total_ram_stake"`
	LastProducerScheduleUpdate     string  `json:"last_producer_schedule_update"`
	LastPervoteBucketFill          Int64   `json:"last_pervote_bucket_fill,string"`
	PervoteBucket                  int     `json:"pervote_bucket"`
	PerblockBucket                 int     `json:"perblock_bucket"`
	TotalUnpaidBlocks              int     `json:"total_unpaid_blocks"`
	TotalActivatedStake            float64 `json:"total_activated_stake,string"`
	ThreshActivatedStakeTime       Int64   `json:"thresh_activated_stake_time,string"`
	LastProducerScheduleSize       int     `json:"last_producer_schedule_size"`
	TotalProducerVoteWeight        float64 `json:"total_producer_vote_weight,string"`
	LastNameClose                  string  `json:"last_name_close"`
}

type Producer struct {
	Owner         string      `json:"owner"`
	TotalVotes    float64     `json:"total_votes,string"`
	ProducerKey   string      `json:"producer_key"`
	IsActive      int         `json:"is_active"`
	URL           string      `json:"url"`
	UnpaidBlocks  int         `json:"unpaid_blocks"`
	LastClaimTime JSONFloat64 `json:"last_claim_time"`
	Location      int         `json:"location"`
}
type ProducersResp struct {
	Producers []Producer `json:"producers"`
}
type GetActionsRequest struct {
	AccountName AccountName `json:"account_name"`
	Pos         Int64       `json:"pos"`
	Offset      Int64       `json:"offset"`
}
type ActionResp struct {
	GlobalSeq  JSONInt64      `json:"global_action_seq"`
	AccountSeq JSONInt64      `json:"account_action_seq"`
	BlockNum   uint32         `json:"block_num"`
	BlockTime  BlockTimestamp `json:"block_time"`
	Trace      ActionTrace    `json:"action_trace"`
}
type ActionsResp struct {
	Actions               []ActionResp `json:"actions"`
	LastIrreversibleBlock uint32       `json:"last_irreversible_block"`
}
type KeyAccountsResp struct {
	AccountNames []string `json:"account_names"`
}

type ControlledAccountsResp struct {
	ControlledAccounts []string `json:"controlled_accounts"`
}

type GetCurrencyStatsResp struct {
	Supply    Asset       `json:"supply"`
	MaxSupply Asset       `json:"max_supply"`
	Issuer    AccountName `json:"issuer"`
}

type Except struct {
	Code    Int64               `json:"code"`
	Name    string              `json:"name"`
	Message string              `json:"message"`
	Stack   []*ExceptLogMessage `json:"stack"`
}

// LogMessage is a line of message in an exception.
type ExceptLogMessage struct {
	Context ExceptLogContext `json:"context"`
	Format  string           `json:"format"`
	Data    json.RawMessage  `json:"data"`
}

var exceptLogMessageTypes = map[string]fcVariantType{
	"context": fcVariantObjectType,
	"format":  fcVariantStringType,
}

func (m *ExceptLogMessage) UnmarshalBinary(decoder *Decoder) error {
	variant := fcVariant{}
	err := decoder.Decode(&variant)
	if err != nil {
		return fmt.Errorf("unable to decode except log message: %s", err)
	}

	if variant.TypeID != fcVariantObjectType {
		return fmt.Errorf("invalid log message, expected type %s, got %s", fcVariantObjectType, variant.TypeID)
	}

	object := variant.MustAsObject()

	if err := object.validateFields(exceptLogMessageTypes); err != nil {
		return fmt.Errorf("invalid log message object: %s", err)
	}

	if err = m.Context.fromObject(object["context"].MustAsObject()); err != nil {
		return fmt.Errorf("unable to assign context: %s", err)
	}

	m.Format = object["format"].MustAsString()

	if dataVariant := object["data"]; !dataVariant.IsNil() {
		if m.Data, err = json.Marshal(dataVariant.ToNative()); err != nil {
			return fmt.Errorf("unable to assign data: %s", err)
		}
	}

	return nil
}

type ExceptLogContext struct {
	Level      ExceptLogLevel    `json:"level"`
	File       string            `json:"file"`
	Line       uint64            `json:"line"`
	Method     string            `json:"method"`
	Hostname   string            `json:"hostname"`
	ThreadName string            `json:"thread_name"`
	Timestamp  JSONTime          `json:"timestamp"`
	Context    *ExceptLogContext `json:"context,omitempty"`
}

var exceptLogContextTypes = map[string]fcVariantType{
	"level":       fcVariantStringType,
	"file":        fcVariantStringType,
	"line":        fcVariantUint64Type,
	"method":      fcVariantStringType,
	"hostname":    fcVariantStringType,
	"thread_name": fcVariantStringType,
	"timestamp":   fcVariantStringType,
	"?context":    fcVariantObjectType,
}

func (c *ExceptLogContext) fromObject(object fcVariantObject) error {
	if err := object.validateFields(exceptLogContextTypes); err != nil {
		return fmt.Errorf("invalid log context: %s", err)
	}

	c.Level.FromString(object["level"].MustAsString())
	c.File = object["file"].MustAsString()
	c.Line = object["line"].MustAsUint64()
	c.Method = object["method"].MustAsString()
	c.Hostname = object["hostname"].MustAsString()
	c.ThreadName = object["thread_name"].MustAsString()

	var err error
	if c.Timestamp, err = ParseJSONTime(object["timestamp"].MustAsString()); err != nil {
		return fmt.Errorf("invalid log context timestamp: %s", err)
	}

	contextVariant := object["context"]
	if contextVariant.TypeID != fcVariantNullType {
		c.Context = new(ExceptLogContext)
		if err := c.Context.fromObject(contextVariant.MustAsObject()); err != nil {
			return fmt.Errorf("unable to assign nested context: %s", err)
		}
	}

	return nil
}

type ExceptLogLevel uint8

const (
	ExceptLogLevelAll ExceptLogLevel = iota
	ExceptLogLevelDebug
	ExceptLogLevelInfo
	ExceptLogLevelWarn
	ExceptLogLevelError
	ExceptLogLevelOff
)

func (s *ExceptLogLevel) FromString(input string) {
	switch input {
	case "all":
		*s = ExceptLogLevelAll
	case "debug":
		*s = ExceptLogLevelDebug
	case "info":
		*s = ExceptLogLevelInfo
	case "warn":
		*s = ExceptLogLevelWarn
	case "error":
		*s = ExceptLogLevelError
	case "off":
		*s = ExceptLogLevelOff
	default:
		*s = ExceptLogLevelOff
	}
}

func (s *ExceptLogLevel) UnmarshalJSON(data []byte) error {
	var decoded string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	s.FromString(decoded)
	return nil
}

func (s ExceptLogLevel) MarshalJSON() (data []byte, err error) {
	out := "off"
	switch s {
	case ExceptLogLevelAll:
		out = "all"
	case ExceptLogLevelDebug:
		out = "debug"
	case ExceptLogLevelInfo:
		out = "info"
	case ExceptLogLevelWarn:
		out = "warn"
	case ExceptLogLevelError:
		out = "error"
	case ExceptLogLevelOff:
		out = "off"
	}
	return json.Marshal(out)
}

func (s ExceptLogLevel) String() string {
	switch s {
	case ExceptLogLevelAll:
		return "all"
	case ExceptLogLevelDebug:
		return "debug"
	case ExceptLogLevelInfo:
		return "info"
	case ExceptLogLevelWarn:
		return "warn"
	case ExceptLogLevelError:
		return "error"
	case ExceptLogLevelOff:
		return "off"
	}

	return "off"
}
