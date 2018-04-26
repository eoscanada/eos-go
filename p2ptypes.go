package eos

import (
	"fmt"

	"github.com/eoscanada/eos-go/ecc"
)

type P2PMessage interface {
	GetType() P2PMessageType
}

type HandshakeMessage struct {
	// net_plugin/protocol.hpp handshake_message
	NetworkVersion           int16         `json:"network_version"`
	ChainID                  SHA256Bytes   `json:"chain_id"`
	NodeID                   SHA256Bytes   `json:"node_id"` // sha256
	Key                      ecc.PublicKey `json:"key"`     // can be empty, producer key, or peer key
	Time                     Tstamp        `json:"time"`    // time?!
	Token                    SHA256Bytes   `json:"token"`   // digest of time to prove we own the private `key`
	Signature                ecc.Signature `json:"sig"`     // can be empty if no key, signature of the digest above
	P2PAddress               string        `json:"p2p_address"`
	LastIrreversibleBlockNum uint32        `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  SHA256Bytes   `json:"last_irreversible_block_id"`
	HeadNum                  uint32        `json:"head_num"`
	HeadID                   SHA256Bytes   `json:"head_id"`
	OS                       string        `json:"os"`
	Agent                    string        `json:"agent"`
	Generation               int16         `json:"generation"`
}

func (m *HandshakeMessage) GetType() P2PMessageType {
	return HandshakeMessageType
}

func (m *HandshakeMessage) String() string {
	return fmt.Sprintf("Handshake: Head [%d] Last Irreversible [%d] Time [%s]", m.HeadNum, m.LastIrreversibleBlockNum, m.Time)
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
	Reason GoAwayReason `json:"reason"`
	NodeID SHA256Bytes  `json:"node_id"`
}

func (m *GoAwayMessage) GetType() P2PMessageType {
	return GoAwayMessageType
}

type TimeMessage struct {
	Origin      Tstamp `json:"org"`
	Receive     Tstamp `json:"rec"`
	Transmit    Tstamp `json:"xmt"`
	Destination Tstamp `json:"dst"`
}

func (m *TimeMessage) GetType() P2PMessageType {
	return TimeMessageType
}

func (t *TimeMessage) String() string {
	return fmt.Sprintf("Origin [%s], Receive [%s], Transmit [%s], Destination [%s]", t.Origin, t.Receive, t.Transmit, t.Destination)
}

type TransactionStatus uint8

const (
	TransactionStatusExecuted TransactionStatus = iota ///< succeed, no error handler executed
	TransactionStatusSoftFail                          ///< objectively failed (not executed), error handler executed
	TransactionStatusHardFail                          ///< objectively failed and error handler objectively failed thus no state change
	TransactionStatusDelayed                           ///< transaction delayed
)

//type TransactionID SHA256Bytes

type TransactionReceipt struct {
	Status        TransactionStatus `json:"status"`
	KCPUUsage     Varuint32         `json:"kcpu_usage"`
	NetUsageWords Varuint32         `json:"net_usage_words"`
	ID            SHA256Bytes       `json:"id"`
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
	AccountName     AccountName   `json:"account_name"`
	BlockSigningKey ecc.PublicKey `json:"block_signing_key"`
}

type ProducerSchedule struct {
	Version   uint32        `json:"version"`
	Producers []ProducerKey `json:"producers"`
}

type BlockHeader struct {
	Previous         SHA256Bytes              `json:"previous"`
	Timestamp        BlockTimestamp           `json:"timestamp"`
	TransactionMRoot SHA256Bytes              `json:"transaction_mroot"`
	ActionMRoot      SHA256Bytes              `json:"action_mroot"`
	BlockMRoot       SHA256Bytes              `json:"block_mroot"`
	Producer         AccountName              `json:"producer"`
	ScheduleVersion  uint32                   `json:"schedule_version"`
	NewProducers     OptionalProducerSchedule `json:"new_producers"`
}

type OptionalProducerSchedule struct {
	ProducerSchedule
}

func (a *OptionalProducerSchedule) OptionalBinaryMarshalerPresent() bool {
	return a == nil
}

type SignedBlockHeader struct {
	BlockHeader
	ProducerSignature ecc.Signature `json:"producer_signature"`
}

type SignedBlockSummaryMessage struct {
	SignedBlockHeader
	Regions []RegionSummary `json:"regions"`
}

func (m *SignedBlockSummaryMessage) GetType() P2PMessageType {
	return SignedBlockSummaryMessageType
}

type SignedBlockMessage struct {
	SignedBlockSummaryMessage
	InputTransactions []PackedTransaction `json:"input_transactions"`
}

func (m *SignedBlockMessage) String() string {
	return "SignedBlockMessage"
}

func (m *SignedBlockMessage) GetType() P2PMessageType {
	return SignedBlockMessageType
}

type IDListMode uint8

const (
	none IDListMode = iota
	catch_up
	last_irr_catch_up
	normal
)

type OrderedTransactionIDs struct {
	Caca    [3]byte
	Mode    IDListMode    `json:"mode"`
	Pending uint32        `json:"pending"`
	IDs     []SHA256Bytes `json:"ids"`
}
type OrderedBlockIDs struct {
	Caca    [3]byte
	Mode    IDListMode    `json:"mode"`
	Pending uint32        `json:"pending"`
	IDs     []SHA256Bytes `json:"ids"`
}

type NoticeMessage struct {
	KnownTrx    OrderedBlockIDs `json:"known_trx"`
	KnownBlocks OrderedBlockIDs `json:"known_blocks"`
}

func (m *NoticeMessage) GetType() P2PMessageType {
	return NoticeMessageType
}

type SyncRequestMessage struct {
	StartBlock uint32 `json:"start_block"`
	EndBlock   uint32 `json:"end_block"`
}

func (m *SyncRequestMessage) GetType() P2PMessageType {
	return SyncRequestMessageType
}
func (m *SyncRequestMessage) String() string {
	return fmt.Sprintf("SyncRequest: Start Block [%d] End Block [%d]", m.StartBlock, m.EndBlock)
}

type RequestMessage struct {
	ReqTrx    OrderedBlockIDs `json:"req_trx"`
	ReqBlocks OrderedBlockIDs `json:"req_blocks"`
}

func (m *RequestMessage) GetType() P2PMessageType {
	return RequestMessageType
}

type SignedTransactionMessage struct {
	Signatures      []ecc.Signature `json:"signatures"`
	ContextFreeData []byte          `json:"context_free_data"`
}

func (m *SignedTransactionMessage) GetType() P2PMessageType {
	return SignedTransactionMessageType
}

type PackedTransactionMessage struct {
	Signatures            []ecc.Signature `json:"signatures"`
	Compression           CompressionType `json:"compression"`
	PackedContextFreeData []byte          `json:"packed_context_free_data"`
	PackedTrx             []byte          `json:"packed_trx"`
}

func (m *PackedTransactionMessage) GetType() P2PMessageType {
	return PackedTransactionMessageType
}
