package eos

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/tidwall/gjson"
)

type P2PMessage interface {
	fmt.Stringer
	GetType() P2PMessageType
}

type HandshakeMessage struct {
	// net_plugin/protocol.hpp handshake_message
	NetworkVersion           uint16        `json:"network_version"`
	ChainID                  Checksum256   `json:"chain_id"`
	NodeID                   Checksum256   `json:"node_id"` // sha256
	Key                      ecc.PublicKey `json:"key"`     // can be empty, producer key, or peer key
	Time                     Tstamp        `json:"time"`    // time?!
	Token                    Checksum256   `json:"token"`   // digest of time to prove we own the private `key`
	Signature                ecc.Signature `json:"sig"`     // can be empty if no key, signature of the digest above
	P2PAddress               string        `json:"p2p_address"`
	LastIrreversibleBlockNum uint32        `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  Checksum256   `json:"last_irreversible_block_id"`
	HeadNum                  uint32        `json:"head_num"`
	HeadID                   Checksum256   `json:"head_id"`
	OS                       string        `json:"os"`
	Agent                    string        `json:"agent"`
	Generation               int16         `json:"generation"`
}

func (m *HandshakeMessage) GetType() P2PMessageType {
	return HandshakeMessageType
}

type ChainSizeMessage struct {
	LastIrreversibleBlockNum uint32      `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  Checksum256 `json:"last_irreversible_block_id"`
	HeadNum                  uint32      `json:"head_num"`
	HeadID                   Checksum256 `json:"head_id"`
}

func (m *ChainSizeMessage) GetType() P2PMessageType {
	return ChainSizeType
}

func (m *HandshakeMessage) String() string {
	return fmt.Sprintf("handshake: Head [%d] 	Last Irreversible [%d] Time [%s]", m.HeadNum, m.LastIrreversibleBlockNum, m.Time)
}

type GoAwayReason uint8

// See plugins/net_plugin/include/eosio/net_plugin/protocol.hpp#L39
const (
	GoAwayNoReason = GoAwayReason(iota)
	GoAwaySelfConnect
	GoAwayDuplicate
	GoAwayWrongChain
	GoAwayWrongVersion
	GoAwayForked
	GoAwayUnlinkable
	GoAwayBadTransaction
	GoAwayValidation
	GoAwayBenignOther
	GoAwayFatalOther
	GoAwayAuthentication
)

var goAwayToStringMap = map[GoAwayReason]string{
	GoAwayNoReason:       "no reason",
	GoAwaySelfConnect:    "self connect",
	GoAwayDuplicate:      "duplicate",
	GoAwayWrongChain:     "wrong chain",
	GoAwayWrongVersion:   "wrong version",
	GoAwayForked:         "chain is forked",
	GoAwayUnlinkable:     "unlinkable block received",
	GoAwayBadTransaction: "bad transaction",
	GoAwayValidation:     "invalid block",
	GoAwayAuthentication: "authentication failure",
	GoAwayFatalOther:     "some other failure",
	GoAwayBenignOther:    "some other non-fatal condition, possibly unknown block",
}

func (r GoAwayReason) String() string {
	if value, exists := goAwayToStringMap[r]; exists {
		return value
	}

	return "some crazy reason"
}

type GoAwayMessage struct {
	Reason GoAwayReason `json:"reason"`
	NodeID Checksum256  `json:"node_id"`
}

func (m *GoAwayMessage) GetType() P2PMessageType {
	return GoAwayMessageType
}

func (m *GoAwayMessage) String() string {
	return fmt.Sprintf("go away: reason [%d]", m.Reason)
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
	TransactionStatusExpired                           ///< transaction expired
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
	case "expired":
		*s = TransactionStatusExpired
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
	case TransactionStatusExpired:
		out = "expired"
	}
	return json.Marshal(out)
}
func (s TransactionStatus) String() string {

	switch s {
	case TransactionStatusExecuted:
		return "executed"
	case TransactionStatusSoftFail:
		return "soft_fail"
	case TransactionStatusHardFail:
		return "hard_fail"
	case TransactionStatusDelayed:
		return "delayed"
	case TransactionStatusExpired:
		return "expired"
	default:
		return "unknown"
	}

}

//type TransactionID Checksum256

type ProducerKey struct {
	AccountName     AccountName   `json:"producer_name"`
	BlockSigningKey ecc.PublicKey `json:"block_signing_key"`
}

type ProducerSchedule struct {
	Version   uint32        `json:"version"`
	Producers []ProducerKey `json:"producers"`
}

type ProducerAuthoritySchedule struct {
	Version   uint32               `json:"version"`
	Producers []*ProducerAuthority `json:"producers"`
}

type ProducerAuthority struct {
	AccountName           AccountName            `json:"producer_name"`
	BlockSigningAuthority *BlockSigningAuthority `json:"authority"`
}

type MerkleRoot struct {
	ActiveNodes []Checksum256 `json:"_active_nodes"`
	NodeCount   uint64        `json:"_node_count"`
}

type PairAccountNameBlockNum struct {
	AccountName AccountName
	BlockNum    uint32
}

func (c PairAccountNameBlockNum) MarshalJSON() ([]byte, error) {
	return []byte(`["` + string(c.AccountName) + `",` + strconv.FormatUint(uint64(c.BlockNum), 10) + `]`), nil
}

func (c *PairAccountNameBlockNum) UnmarshalJSON(data []byte) error {
	var object [2]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return nil
	}

	accountName, ok := object[0].(string)
	if !ok {
		return fmt.Errorf("expecting first element to be a string, got %T", object[0])
	}

	blockNum, ok := object[1].(float64)
	if !ok {
		return fmt.Errorf("expecting second element to be a float64, got %T", object[1])
	}

	c.AccountName = AN(accountName)
	c.BlockNum = uint32(blockNum)

	return nil
}

func (c *PairAccountNameBlockNum) UnmarshalBinary(decoder *Decoder) error {
	accountName, err := decoder.ReadName()
	if err != nil {
		return fmt.Errorf("unable to read pair account name: %w", err)
	}

	blockNum, err := decoder.ReadUint32()
	if err != nil {
		return fmt.Errorf("unable to read pair block num: %w", err)
	}

	c.AccountName = AccountName(accountName)
	c.BlockNum = blockNum
	return nil
}

// FIXME: This structure supports both EOS 1.8.x as well as EOS 2.0.x. However, the binary encoding
//        format does only support the 2.0.x version for now. It's not clear how we would do thing
//        to propagate the information that encoding/decoding of binary should be performed with one
//        variant or the other. When this comment was added, the binary encoding/decoding was not
//        working for either version, so supporting EOS 2.0.x only is a fair improvements. Will need
//        to understand better if this is required for other chains for example.
type BlockState struct {
	BlockNum                         uint32 `json:"block_num"`
	DPoSProposedIrreversibleBlockNum uint32 `json:"dpos_proposed_irreversible_blocknum"`
	DPoSIrreversibleBlockNum         uint32 `json:"dpos_irreversible_blocknum"`

	// Hybrid (dynamic types)
	ActiveSchedule *ProducerScheduleOrAuthoritySchedule `json:"active_schedule"`

	BlockrootMerkle          *MerkleRoot               `json:"blockroot_merkle,omitempty"`
	ProducerToLastProduced   []PairAccountNameBlockNum `json:"producer_to_last_produced,omitempty"`
	ProducerToLastImpliedIRB []PairAccountNameBlockNum `json:"producer_to_last_implied_irb,omitempty"`

	// EOSIO 2.x
	ValidBlockSigningAuthorityV2 *BlockSigningAuthority `json:"valid_block_signing_authority,omitempty"`

	ConfirmCount []uint8 `json:"confirm_count,omitempty"`

	BlockID                   Checksum256                   `json:"id"`
	Header                    *SignedBlockHeader            `json:"header,omitempty"`
	PendingSchedule           *PendingSchedule              `json:"pending_schedule"`
	ActivatedProtocolFeatures *ProtocolFeatureActivationSet `json:"activated_protocol_features,omitempty" eos:"optional"`
	AdditionalSignatures      []ecc.Signature               `json:"additional_signatures"`

	SignedBlock *SignedBlock `json:"block,omitempty" eos:"optional"`
	Validated   bool         `json:"validated"`

	// EOSIO 1.x
	BlockSigningKeyV1 *ecc.PublicKey `json:"block_signing_key,omitempty" eos:"-"`
}

type ProducerScheduleOrAuthoritySchedule struct {
	// EOSIO 1.x
	V1 *ProducerSchedule

	// EOSIO 2.x
	V2 *ProducerAuthoritySchedule
}

func (p *ProducerScheduleOrAuthoritySchedule) MarshalJSON() ([]byte, error) {
	// In case of ambiguity, which arise only on empty `producers` array, the first one is picked since it does not matter (same JSON expectOutput)
	if p.V1 != nil {
		return json.Marshal(p.V1)
	}

	if p.V2 != nil {
		return json.Marshal(p.V2)
	}

	return nil, fmt.Errorf("both V1 and V2 were null, this is an error")
}

func (p *ProducerScheduleOrAuthoritySchedule) UnmarshalJSON(data []byte) error {
	versionResult := gjson.GetBytes(data, "version")
	if !versionResult.Exists() || versionResult.Type != gjson.Number {
		return fmt.Errorf("expected 'version' key of type 'number' to exist in %q", string(data))
	}

	producersResult := gjson.GetBytes(data, "producers")
	if !producersResult.Exists() || !producersResult.IsArray() {
		return fmt.Errorf("expected 'producers' key of type 'number' to exist in %q", string(data))
	}

	// We cannot infer anything, what should we do exactly? We could populate the two, but
	// what happens on marshal? Both are defined, that's what we choose for now, `eos-go` user
	// would then make the choice themselves.
	if len(producersResult.Array()) == 0 || producersResult.Get("0.block_signing_key").Exists() {
		p.V1 = new(ProducerSchedule)
		err := json.Unmarshal(data, p.V1)
		if err != nil {
			return fmt.Errorf("unable to unmarshal ProducerSchedule type: %w", err)
		}
	}

	if len(producersResult.Array()) == 0 || producersResult.Get("0.authority").Exists() {
		p.V2 = new(ProducerAuthoritySchedule)
		err := json.Unmarshal(data, p.V2)
		if err != nil {
			return fmt.Errorf("unable to unmarshal ProducerAuthoritySchedule type: %w", err)
		}
	}

	if p.V1 == nil && p.V2 == nil {
		return errors.New("unable to unmarshal producer authority or schedule, no type could be inferred from JSON")
	}

	return nil
}

func (p *ProducerScheduleOrAuthoritySchedule) UnmarshalBinary(decoder *Decoder) error {
	// FIXME: For now, we assume all blocks are in 2.0 format, so with V2 set. However,
	//        it's not really clear for me yet how is the versionning being handled by
	//        `nodeos` at the binary level. In the block header, there is some header
	//        extensions that can be used to determined if the WTMSIG protocol feature
	//        has been activated or not. But I'm not even sure the binary format changes
	//        or not (most probably that it does). Anyway, this compatibility stuff between
	//        1.8 and 2.0 with WTMSIG protocol feature needs to be researched and handled
	//        correctly everywhere in the library.
	if p.V2 == nil {
		p.V2 = new(ProducerAuthoritySchedule)
	}

	err := decoder.Decode(p.V2)
	if err != nil {
		return fmt.Errorf("unable to decode producer authority schedule (V2): %w", err)
	}

	return nil

}

// See libraries/chain/include/eosio/chain/producer_schedule.hpp#L161
var BlockSigningAuthorityVariant = NewVariantDefinition([]VariantType{
	{"block_signing_authority_v0", (*BlockSigningAuthorityV0)(nil)},
})

type BlockSigningAuthority struct {
	BaseVariant
}

func (a *BlockSigningAuthority) MarshalJSON() ([]byte, error) {
	return a.BaseVariant.MarshalJSON(BlockSigningAuthorityVariant)
}

func (a *BlockSigningAuthority) UnmarshalJSON(data []byte) error {
	return a.BaseVariant.UnmarshalJSON(data, BlockSigningAuthorityVariant)
}

func (a *BlockSigningAuthority) UnmarshalBinary(decoder *Decoder) error {
	return a.BaseVariant.UnmarshalBinaryVariant(decoder, BlockSigningAuthorityVariant)
}

// See libraries/chain/include/eosio/chain/producer_schedule.hpp#L100
type BlockSigningAuthorityV0 struct {
	Threshold uint32       `json:"threshold"`
	Keys      []*KeyWeight `json:"keys"`
}

type PendingSchedule struct {
	ScheduleLIBNum uint32                               `json:"schedule_lib_num"`
	ScheduleHash   Checksum256                          `json:"schedule_hash"`
	Schedule       *ProducerScheduleOrAuthoritySchedule `json:"schedule"`
}

type ProtocolFeatureActivationSet struct {
	ProtocolFeatures []Checksum256 `json:"protocol_features"`
}

type BlockHeader struct {
	Timestamp        BlockTimestamp `json:"timestamp"`
	Producer         AccountName    `json:"producer"`
	Confirmed        uint16         `json:"confirmed"`
	Previous         Checksum256    `json:"previous"`
	TransactionMRoot Checksum256    `json:"transaction_mroot"`
	ActionMRoot      Checksum256    `json:"action_mroot"`
	ScheduleVersion  uint32         `json:"schedule_version"`

	// EOSIO 1.x
	NewProducersV1 *ProducerSchedule `json:"new_producers,omitempty" eos:"optional"`

	HeaderExtensions []*Extension `json:"header_extensions"`
}

func (b *BlockHeader) BlockNumber() uint32 {
	return binary.BigEndian.Uint32(b.Previous[:4]) + 1
}

func (b *BlockHeader) BlockID() (Checksum256, error) {
	cereal, err := MarshalBinary(b)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	_, _ = h.Write(cereal)
	hashed := h.Sum(nil)

	binary.BigEndian.PutUint32(hashed, b.BlockNumber())

	return Checksum256(hashed), nil
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
	return fmt.Sprintf("SignedBlock [%d] with %d txs", m.BlockNumber(), len(m.Transactions))
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
	ID     Checksum256
	Packed *PackedTransaction
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

		id, err := packed.ID()
		if err != nil {
			return fmt.Errorf("get id: %w", err)
		}

		*t = TransactionWithID{
			ID:     id,
			Packed: &packed,
		}

		return nil
	} else if data[0] == '"' {
		var id string
		err := json.Unmarshal(data, &id)
		if err != nil {
			return err
		}

		shaID, err := hex.DecodeString(id)
		if err != nil {
			return fmt.Errorf("decoding id in trx: %w", err)
		}

		*t = TransactionWithID{
			ID: Checksum256(shaID),
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

	typ := string(in[0])
	switch typ {
	case "0":
		var s string
		if err := json.Unmarshal(in[1], &s); err != nil {
			return err
		}

		*t = TransactionWithID{}
		if err := json.Unmarshal(in[1], &t.ID); err != nil {
			return err
		}
	case "1":

		// ignore the ID field right now..
		err = json.Unmarshal(in[1], &packed)
		if err != nil {
			return err
		}

		id, err := packed.ID()
		if err != nil {
			return fmt.Errorf("get id: %w", err)
		}

		*t = TransactionWithID{
			ID:     id,
			Packed: &packed,
		}
	default:
		return fmt.Errorf("unsupported multi-variant trx serialization type from C++ code into Go: %q", typ)
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
	Mode    [4]byte       `json:"mode"`
	Pending uint32        `json:"pending"`
	IDs     []Checksum256 `json:"ids"`
}
type OrderedBlockIDs struct {
	Mode    [4]byte       `json:"mode"`
	Pending uint32        `json:"pending"`
	IDs     []Checksum256 `json:"ids"`
}

func (o *OrderedBlockIDs) String() string {

	ids := ""
	for _, id := range o.IDs {
		ids += fmt.Sprintf("%s,", id)
	}
	return fmt.Sprintf("Mode %d, Pending %d, ids [%s]", o.Mode, o.Pending, ids)
}

type NoticeMessage struct {
	KnownTrx    OrderedBlockIDs `json:"known_trx"`
	KnownBlocks OrderedBlockIDs `json:"known_blocks"`
}

func (n *NoticeMessage) String() string {
	return fmt.Sprintf("KnownTrx %s :: KnownBlocks %s", n.KnownTrx.String(), n.KnownBlocks.String())
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

func (r *RequestMessage) String() string {
	return fmt.Sprintf("ReqTrx %s :: ReqBlocks %s", r.ReqTrx.String(), r.ReqBlocks.String())
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

func (m PackedTransactionMessage) String() string {
	signTrx, err := m.Unpack()
	if err != nil {
		return fmt.Sprintf("err trx msg unpack by %s", err.Error())
	}
	return signTrx.String()
}
