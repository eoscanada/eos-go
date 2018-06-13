package eos

import (
	"crypto/sha256"
	"fmt"

	"encoding/binary"
	"encoding/json"

	"github.com/eoscanada/eos-go/ecc"
)

type P2PMessage interface {
	GetType() P2PMessageType
}

type HandshakeMessage struct {
	// net_plugin/protocol.hpp handshake_message
	NetworkVersion           uint16        `json:"network_version"`
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

type ChainSizeMessage struct {
	LastIrreversibleBlockNum uint32      `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  SHA256Bytes `json:"last_irreversible_block_id"`
	HeadNum                  uint32      `json:"head_num"`
	HeadID                   SHA256Bytes `json:"head_id"`
}

func (m *ChainSizeMessage) GetType() P2PMessageType {
	return ChainSizeType
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
	TransactionStatusUnknown  = TransactionStatus(255)
)

func (s *TransactionStatus) UnmarshalJSON(data []byte) error {
	var decoded string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	switch decoded {
	case "executed":
		*s = TransactionStatusExecuted
	case "soft_fail":
		*s = TransactionStatusSoftFail

	case "hard_fail":
		*s = TransactionStatusHardFail
	case "delayed":
		*s = TransactionStatusDelayed
	default:
		*s = TransactionStatusUnknown
	}
	return nil
}

func (s TransactionStatus) MarshalJSON() (data []byte, err error) {
	out := "unknown"
	switch s {
	case TransactionStatusExecuted:
		out = "executed"
	case TransactionStatusSoftFail:
		out = "soft_fail"
	case TransactionStatusHardFail:
		out = "hard_fail"
	case TransactionStatusDelayed:
		out = "delayed"
	}
	return json.Marshal(out)
}
func (s TransactionStatus) String() string {

	switch s {
	case TransactionStatusExecuted:
		return "executed"
	case TransactionStatusSoftFail:
		return "soft fail"
	case TransactionStatusHardFail:
		return "hard fail"
	case TransactionStatusDelayed:
		return "delayed"
	default:
		return "unknown"
	}

}

//type TransactionID SHA256Bytes

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
	Region        uint16   `json:"region"`
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
	Timestamp        BlockTimestamp            `json:"timestamp"`
	Producer         AccountName               `json:"producer"`
	Confirmed        uint16                    `json:"confirmed"`
	Previous         SHA256Bytes               `json:"previous"`
	TransactionMRoot SHA256Bytes               `json:"transaction_mroot"`
	ActionMRoot      SHA256Bytes               `json:"action_mroot"`
	ScheduleVersion  uint32                    `json:"schedule_version"`
	NewProducers     *OptionalProducerSchedule `json:"new_producers" eos:"optional"`
	HeaderExtensions []*Extension              `json:"header_extensions"`
}

func (b *BlockHeader) BlockNumber() uint32 {
	return binary.BigEndian.Uint32(b.Previous[:4]) + 1
}

func (b *BlockHeader) BlockID() (SHA256Bytes, error) {
	cereal, err := MarshalBinary(b)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	_, _ = h.Write(cereal)
	hashed := h.Sum(nil)

	binary.BigEndian.PutUint32(hashed, b.BlockNumber())

	return SHA256Bytes(hashed), nil
}

type OptionalProducerSchedule struct {
	ProducerSchedule
}

type SignedBlockHeader struct {
	BlockHeader
	ProducerSignature ecc.Signature `json:"producer_signature"`
}

type SignedBlock struct {
	SignedBlockHeader
	Transactions    []TransactionReceipt `json:"transactions"`
	BlockExtensions []*Extension         `json:"block_extensions"`
}

func (m *SignedBlock) String() string {
	return "SignedBlock"
}

func (m *SignedBlock) GetType() P2PMessageType {
	return SignedBlockType
}

type TransactionReceiptHeader struct {
	Status               TransactionStatus `json:"status"`
	CPUUsageMicroSeconds uint32            `json:"cpu_usage_us"`
	NetUsageWords        Varuint32         `json:"net_usage_words"`
}

type TransactionReceipt struct {
	TransactionReceiptHeader
	Transaction TransactionWithID `json:"trx"`
}

type TransactionWithID struct {
	ID     uint8 //
	Packed PackedTransaction
}

func (t TransactionWithID) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		t.ID,
		t.Packed,
	})
}

func (t *TransactionWithID) UnmarshalJSON(data []byte) error {
	var packed PackedTransaction
	if data[0] == '{' {
		if err := json.Unmarshal(data, &packed); err != nil {
			return err
		}
		*t = TransactionWithID{
			Packed: packed,
		}

		return nil
	}

	var in []json.RawMessage
	err := json.Unmarshal(data, &in)
	if err != nil {
		return err
	}

	if len(in) != 2 {
		return fmt.Errorf("expected two params for TransactionWithID, got %d", len(in))
	}

	// ignore the ID field right now..
	err = json.Unmarshal(in[1], &packed)
	if err != nil {
		return err
	}

	*t = TransactionWithID{
		Packed: packed,
	}

	return nil
}

type IDListMode byte

const (
	none IDListMode = iota
	catch_up
	last_irr_catch_up
	normal
)

type OrderedTransactionIDs struct {
	Unknown [3]byte       `json:"-"` ///// WWUUuuuuuuuuuuuutzthat ?
	Mode    IDListMode    `json:"mode"`
	Pending uint32        `json:"pending"`
	IDs     []SHA256Bytes `json:"ids"`
}
type OrderedBlockIDs struct {
	Unknown [3]byte       `json:"-"` ///// wuuttzthat?
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

type PackedTransactionMessage struct {
	PackedTransaction
}

func (m *PackedTransactionMessage) GetType() P2PMessageType {
	return PackedTransactionMessageType
}
