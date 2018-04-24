package eos

import "github.com/eoscanada/eos-go/ecc"

type P2PMessage interface {
}

type HandshakeMessage struct {
	// net_plugin/protocol.hpp handshake_message
	NetworkVersion           int16         `json:"network_version"`
	ChainID                  HexBytes      `json:"chain_id"`
	NodeID                   HexBytes      `json:"node_id"` // sha256
	Key                      ecc.PublicKey `json:"key"`     // can be empty, producer key, or peer key
	Time                     Tstamp        `json:"time"`    // time?!
	Token                    HexBytes      `json:"token"`   // digest of time to prove we own the private `key`
	Signature                ecc.Signature `json:"sig"`     // can be empty if no key, signature of the digest above
	P2PAddress               string        `json:"p2p_address"`
	LastIrreversibleBlockNum uint32        `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  HexBytes      `json:"last_irreversible_block_id"`
	HeadNum                  uint32        `json:"head_num"`
	HeadID                   HexBytes      `json:"head_id"`
	OS                       string        `json:"os"`
	Agent                    string        `json:"agent"`
	Generation               int16         `json:"generation"`
}

type GoAwayReason uint8

const (
	GoAwayNoReason = uint8(iota)
	GoAwaySelfConnect
	GoAwayDuplicate
	GoAwayWrongChain
	GoAwayWrongVersion
	GoAwayForked
	GoAwayUnlinkable
	GoAwayBadTransaction
	GoAwayValidation
	GoAwayAuthentication
	GoAwayFatalOther
	GoAwayBenignOther
	GoAwayCrazy
)

type GoAwayMessage struct {
	GoAwayReason
}

type TimeMessage struct {
	Origin      Tstamp `json:"org"`
	Receive     Tstamp `json:"rec"`
	Transmit    Tstamp `json:"xmt"`
	Destination Tstamp `json:"dst"`
}

type TransactionStatus uint8

const (
	TransactionStatusExecuted TransactionStatus = iota ///< succeed, no error handler executed
	TransactionStatusSoftFail                          ///< objectively failed (not executed), error handler executed
	TransactionStatusHardFail                          ///< objectively failed and error handler objectively failed thus no state change
	TransactionStatusSelayed                           ///< transaction delayed
)

type TransactionId SHA256Bytes

type TransactionReceipt struct {
	Status       TransactionStatus `json:"status"`
	KCPUUsage    uint32            `json:"kcpu_usage"`
	NeUsageWords uint32            `json:"net_usage_words"`
	Id           TransactionId     `json:"id"`
}

type ShardLock struct {
	AccountName AccountName `json:"account_name"`
	ScopeName   ScopeName   `json:"scope_name"`
}

type ShardSummary struct {
	ReadLocks    []ShardLock          `json:"read_locks"`
	WriteLocks   []ShardLock          `json:"write_locks"`
	Transactions []TransactionReceipt `json:"transactions"`
}

type Cycles []ShardSummary
type RegionSummary struct {
	Region        uint16
	CyclesSummary []Cycles `json:"cycles_summary"`
}

type ProducerKey struct {
	AccountName     AccountName `json:"account_name"`
	BlockSigningKey SHA256Bytes `json:"block_signing_key"` //todo: Surely not good
}

type ProducerScheduleType struct {
	Version   uint32        `json:"version"`
	Producers []ProducerKey `json:"producers"`
}

type BlockHeader struct {
	Digest           SHA256Bytes   `json:"digest"`
	BlockNumber      uint32        `json:"block_number"`
	NumFromId        uint32        `json:"num_from_id"`
	Previous         SHA256Bytes   `json:"previous"`
	Timestamp        Tstamp        `json:"timestamp"`
	TransactionMRoot SHA256Bytes   `json:"transaction_mroot"`
	ActionMRoot      SHA256Bytes   `json:"action_mroot"`
	BlockMRoot       SHA256Bytes   `json:"block_mroot"`
	Producer         AccountName   `json:"producer"`
	ScheduleVersion  uint32        `json:"schedule_version"`
	NewProducers     []ProducerKey `json:"new_producers"`
}

type SignedBlockHeader struct {
	BlockHeader
	ProducerSignature SHA256Bytes `json:"producer_signature"` //todo: Surely not good
}

type SignedBlockSummaryMessage struct {
	SignedBlockHeader
	Regions []RegionSummary `json:"regions"`
}

type SignedBlockMessage struct {
	SignedBlockSummaryMessage
	InputTransactions []PackedTransaction `json:"input_transactions"`
}
