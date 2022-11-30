package eos

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"time"

	"github.com/eoscanada/eos-go/ecc"
)

type TransactionHeader struct {
	Expiration     JSONTime `json:"expiration"`
	RefBlockNum    uint16   `json:"ref_block_num"`
	RefBlockPrefix uint32   `json:"ref_block_prefix"`

	MaxNetUsageWords Varuint32 `json:"max_net_usage_words"`
	MaxCPUUsageMS    uint8     `json:"max_cpu_usage_ms"`
	DelaySec         Varuint32 `json:"delay_sec"` // number of secs to delay, making it cancellable for that duration
}

// Transaction represents an EOS transaction that is no signed yet.
//
// **Note** In EOSIO codebase, used within both `signed_transaction` and `transaction`.
type Transaction struct {
	TransactionHeader

	ContextFreeActions []*Action    `json:"context_free_actions"`
	Actions            []*Action    `json:"actions"`
	Extensions         []*Extension `json:"transaction_extensions"`
}

// NewTransaction creates a transaction. Unless you plan on adding HeadBlockID later, to be complete, opts should contain it.  Sign
func NewTransaction(actions []*Action, opts *TxOptions) *Transaction {
	if opts == nil {
		opts = &TxOptions{}
	}

	tx := &Transaction{Actions: actions}
	tx.Fill(opts.HeadBlockID, opts.DelaySecs, opts.MaxNetUsageWords, opts.MaxCPUUsageMS)
	return tx
}

func (tx *Transaction) SetExpiration(in time.Duration) {
	tx.Expiration = JSONTime{time.Now().UTC().Add(in)}
}

const (
	EOS_ProtocolFeatureActivation BlockHeaderExtensionType = iota
	EOS_ProducerScheduleChangeExtension
)

type BlockHeaderExtension interface {
	TypeID() BlockHeaderExtensionType
}

type BlockHeaderExtensionType uint16

type blockHeaderExtensionMap = map[BlockHeaderExtensionType]newBlockHeaderExtension
type newBlockHeaderExtension func() BlockHeaderExtension

var blockHeaderExtensions = map[string]blockHeaderExtensionMap{
	"EOS": {
		EOS_ProtocolFeatureActivation:       func() BlockHeaderExtension { return new(ProtocolFeatureActivationExtension) },
		EOS_ProducerScheduleChangeExtension: func() BlockHeaderExtension { return new(ProducerScheduleChangeExtension) },
	},
}

type Extension struct {
	Type uint16
	Data HexBytes
}

func (e *Extension) MarshalJSON() ([]byte, error) {
	pair := make([]interface{}, 2)
	pair[0] = e.Type
	pair[1] = e.Data

	return json.Marshal(&pair)
}

func (e *Extension) UnmarshalJSON(data []byte) error {
	var pair []interface{}
	err := json.Unmarshal(data, &pair)
	if err != nil {
		return err
	}

	typeFloat, ok := pair[0].(float64)
	if !ok {
		return unmarshalTypeError(pair[0], typeFloat, e, "Type")
	}

	if typeFloat < 0 || typeFloat > math.MaxUint16 {
		return unmarshalTypeError(pair[0], e.Type, e, "Type")
	}

	e.Type = uint16(typeFloat)

	dataHexString, ok := pair[1].(string)
	if !ok {
		return unmarshalTypeError(pair[1], dataHexString, e, "Data")
	}

	e.Data, err = hex.DecodeString(dataHexString)
	if err != nil {
		return unmarshalTypeError(pair[1], e.Data, e, "Data")
	}

	return nil
}

func (e *Extension) UnmarshalBinary(decoder *Decoder) error {
	typeID, err := decoder.ReadUint16()
	if err != nil {
		return fmt.Errorf("unable to read extension type id: %w", err)
	}

	data, err := decoder.ReadByteArray()
	if err != nil {
		return fmt.Errorf("unable to read extension data: %w", err)
	}

	e.Type = typeID
	e.Data = HexBytes(data)
	return nil
}

// AsBlockHeaderExtension turns the given `Extension` object into one of the known
// `BlockHeaderExtension` concrete type.
func (e *Extension) AsBlockHeaderExtension(chain string) (BlockHeaderExtension, error) {
	knownExtensions := blockHeaderExtensions[chain]
	if len(knownExtensions) == 0 {
		return nil, fmt.Errorf("unknown chain identifier %q", chain)
	}

	newPointer := knownExtensions[BlockHeaderExtensionType(e.Type)]
	if newPointer == nil {
		return nil, fmt.Errorf("unknown block header extension type %d for chain %s", e.Type, chain)
	}

	element := newPointer()
	decoder := NewDecoder(e.Data)
	err := decoder.Decode(element)
	if err != nil {
		return nil, fmt.Errorf("unable to decode block header extension: %w", err)
	}

	return element, nil
}

// ProtocolFeatureActivationExtension is a block header extension present in the signed
// block when a particular set of protocol features has been activated by this blockl.
type ProtocolFeatureActivationExtension struct {
	FeatureDigests []Checksum256 `json:"protocol_features"`
}

func (e *ProtocolFeatureActivationExtension) TypeID() BlockHeaderExtensionType {
	return EOS_ProtocolFeatureActivation
}

// ProducerScheduleChangeExtension is a block header extension present in the signed
// block when a new producer authority schedule should be applied.
type ProducerScheduleChangeExtension struct {
	ProducerAuthoritySchedule
}

func (e *ProducerScheduleChangeExtension) TypeID() BlockHeaderExtensionType {
	return EOS_ProducerScheduleChangeExtension
}

func unmarshalTypeError(value interface{}, reflectTypeHost interface{}, target interface{}, field string) *json.UnmarshalTypeError {
	return &json.UnmarshalTypeError{
		Value:  fmt.Sprintf("%T", value),
		Type:   reflect.TypeOf(reflectTypeHost),
		Struct: fmt.Sprintf("%T", target),
		Field:  field,
	}
}

// Fill sets the fields on a transaction.  If you pass `headBlockID`, then `api` can be nil. If you don't pass `headBlockID`, then the `api` is going to be called to fetch
func (tx *Transaction) Fill(headBlockID Checksum256, delaySecs, maxNetUsageWords uint32, maxCPUUsageMS uint8) {
	tx.setRefBlock(headBlockID)

	if tx.ContextFreeActions == nil {
		tx.ContextFreeActions = make([]*Action, 0, 0)
	}
	if tx.Extensions == nil {
		tx.Extensions = make([]*Extension, 0, 0)
	}

	tx.MaxNetUsageWords = Varuint32(maxNetUsageWords)
	tx.MaxCPUUsageMS = maxCPUUsageMS
	tx.DelaySec = Varuint32(delaySecs)

	tx.SetExpiration(30 * time.Second)
}

func (tx *Transaction) setRefBlock(blockID []byte) {
	if len(blockID) == 0 {
		return
	}
	tx.RefBlockNum = uint16(binary.BigEndian.Uint32(blockID[:4]))
	tx.RefBlockPrefix = binary.LittleEndian.Uint32(blockID[8:16])
}

type TransactionTrace struct {
	ID              Checksum256               `json:"id"`
	BlockNum        uint32                    `json:"block_num"`
	BlockTime       BlockTimestamp            `json:"block_time"`
	ProducerBlockID Checksum256               `json:"producer_block_id" eos:"optional"`
	Receipt         *TransactionReceiptHeader `json:"receipt,omitempty" eos:"optional"`
	Elapsed         Int64                     `json:"elapsed"`
	NetUsage        Uint64                    `json:"net_usage"`
	Scheduled       bool                      `json:"scheduled"`
	ActionTraces    []ActionTrace             `json:"action_traces"`
	AccountRamDelta *struct {
		AccountName AccountName `json:"account_name"`
		Delta       Int64       `json:"delta"`
	} `json:"account_ram_delta" eos:"optional"`
	FailedDtrxTrace *TransactionTrace `json:"failed_dtrx_trace,omitempty" eos:"optional"`
	Except          *Except           `json:"except,omitempty" eos:"optional"`
	ErrorCode       *Uint64           `json:"error_code,omitempty" eos:"optional"`
}

type SignedTransaction struct {
	*Transaction

	Signatures      []ecc.Signature `json:"signatures"`
	ContextFreeData []HexBytes      `json:"context_free_data"`

	packed *PackedTransaction
}

func NewSignedTransaction(tx *Transaction) *SignedTransaction {
	return &SignedTransaction{
		Transaction:     tx,
		Signatures:      make([]ecc.Signature, 0),
		ContextFreeData: make([]HexBytes, 0),
	}
}

func (s *SignedTransaction) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (s *SignedTransaction) SignedByKeys(chainID Checksum256) (out []ecc.PublicKey, err error) {
	trx, cfd, err := s.PackedTransactionAndCFD()
	if err != nil {
		return
	}

	for _, sig := range s.Signatures {
		pubKey, err := sig.PublicKey(SigDigest(chainID, trx, cfd))
		if err != nil {
			return nil, err
		}

		out = append(out, pubKey)
	}

	return
}

func (s *SignedTransaction) PackedTransactionAndCFD() ([]byte, []byte, error) {
	rawtrx, err := MarshalBinary(s.Transaction)
	if err != nil {
		return nil, nil, err
	}

	rawcfd := []byte{}
	if len(s.ContextFreeData) > 0 {
		rawcfd, err = MarshalBinary(s.ContextFreeData)
		if err != nil {
			return nil, nil, err
		}
	}

	return rawtrx, rawcfd, nil
}

func (s *SignedTransaction) Pack(compression CompressionType) (*PackedTransaction, error) {
	rawtrx, rawcfd, err := s.PackedTransactionAndCFD()
	if err != nil {
		return nil, err
	}

	switch compression {
	case CompressionZlib:
		var trx bytes.Buffer
		var cfd bytes.Buffer

		// Compress Trx
		writer, _ := zlib.NewWriterLevel(&trx, flate.BestCompression) // can only fail if invalid `level`..
		writer.Write(rawtrx)                                          // ignore error, could only bust memory
		err = writer.Close()
		if err != nil {
			return nil, fmt.Errorf("tx writer close %s", err)
		}
		rawtrx = trx.Bytes()

		// Compress ContextFreeData
		writer, _ = zlib.NewWriterLevel(&cfd, flate.BestCompression) // can only fail if invalid `level`..
		writer.Write(rawcfd)                                         // ignore errors, memory errors only
		err = writer.Close()
		if err != nil {
			return nil, fmt.Errorf("cfd writer close %s", err)
		}
		rawcfd = cfd.Bytes()

	}

	packed := &PackedTransaction{
		Signatures:            s.Signatures,
		Compression:           compression,
		PackedContextFreeData: rawcfd,
		PackedTransaction:     rawtrx,
		wasPackedLocally:      true,
	}

	return packed, nil
}

// PackedTransaction represents a fully packed transaction, with
// signatures, and all. They circulate like that on the P2P net, and
// that's how they are stored.
type PackedTransaction struct {
	Signatures            []ecc.Signature `json:"signatures"`
	Compression           CompressionType `json:"compression"` // in C++, it's an enum, not sure how it Binary-marshals..
	PackedContextFreeData HexBytes        `json:"packed_context_free_data"`
	PackedTransaction     HexBytes        `json:"packed_trx"`
	ContextFreeData       []string        `json:"context_free_data,omitempty" eos:"-"`
	Transaction           *Transaction    `json:"transaction,omitempty" eos:"-"`

	wasPackedLocally bool
}

// ID returns the hash of a transaction.
func (p *PackedTransaction) ID() (Checksum256, error) {
	h := sha256.New()

	if p.wasPackedLocally {
		_, _ = h.Write(p.PackedTransaction)
		return h.Sum(nil), nil
	}

	signed, err := p.UnpackBare()
	if err != nil {
		return nil, err
	}

	repacked, err := signed.Pack(CompressionNone)
	if err != nil {
		return nil, err
	}

	_, _ = h.Write(repacked.PackedTransaction)
	return h.Sum(nil), nil
}

// Unpack decodes the bytestream of the transaction, and attempts to
// decode the registered actions.
func (p *PackedTransaction) Unpack() (signedTx *SignedTransaction, err error) {
	return p.unpack(false)
}

// UnpackBare decodes the transcation payload, but doesn't decode the
// nested action data structure.  See also `Unpack`.
func (p *PackedTransaction) UnpackBare() (signedTx *SignedTransaction, err error) {
	return p.unpack(true)
}

func (p *PackedTransaction) unpack(bare bool) (signedTx *SignedTransaction, err error) {
	var txReader io.Reader
	txReader = bytes.NewBuffer(p.PackedTransaction)

	var freeDataReader io.Reader
	freeDataReader = bytes.NewBuffer(p.PackedContextFreeData)

	switch p.Compression {
	case CompressionZlib:
		txReader, err = zlib.NewReader(txReader)
		if err != nil {
			return nil, fmt.Errorf("new reader for tx, %s", err)
		}

		if len(p.PackedContextFreeData) > 0 {
			freeDataReader, err = zlib.NewReader(freeDataReader)
			if err != nil {
				return nil, fmt.Errorf("new reader for free data, %s", err)
			}
		}
	}

	data, err := ioutil.ReadAll(txReader)
	if err != nil {
		return nil, fmt.Errorf("unpack read all, %s", err)
	}
	decoder := NewDecoder(data)
	decoder.DecodeActions(!bare)

	var tx Transaction
	err = decoder.Decode(&tx)
	if err != nil {
		return nil, fmt.Errorf("unpacking Transaction, %s", err)
	}

	// TODO: wire that in
	//decoder = NewDecoder(freeDataReader)
	//var contextFreeData []HexBytes
	//err = decoder.Decode(&contextFreeData)
	//if err != nil {
	//	fmt.Println("PackedTransaction@freedata err: ", err)
	//	return
	//}

	signedTx = NewSignedTransaction(&tx)
	//signedTx.ContextFreeData = contextFreeData
	signedTx.Signatures = p.Signatures
	signedTx.packed = p

	return
}

type DeferredTransaction struct {
	*Transaction

	SenderID   uint32      `json:"sender_id"`
	Sender     AccountName `json:"sender"`
	DelayUntil JSONTime    `json:"delay_until"`
}

type ScheduledTransaction struct {
	TransactionID Checksum256 `json:"trx_id"`
	Sender        AccountName `json:"sender"`
	SenderID      string      `json:"sender_id"`
	Payer         AccountName `json:"payer"`
	DelayUntil    JSONTime    `json:"delay_until"`
	Expiration    JSONTime    `json:"expiration"`
	Published     JSONTime    `json:"published"`

	Transaction *Transaction `json:"transaction"`
}

// TxOptions represents options you want to pass to the transaction
// you're sending.
type TxOptions struct {
	ChainID          Checksum256 // If specified, we won't hit the API to fetch it
	HeadBlockID      Checksum256 // If provided, don't hit API to fetch it.  This allows offline transaction signing.
	MaxNetUsageWords uint32
	DelaySecs        uint32
	MaxCPUUsageMS    uint8 // If you want to override the CPU usage (in counts of 1024)
	//ExtraKCPUUsage uint32 // If you want to *add* some CPU usage to the estimated amount (in counts of 1024)
	Compress CompressionType
}

// FillFromChain will load ChainID (for signing transactions) and
// HeadBlockID (to fill transaction with TaPoS data).
func (opts *TxOptions) FillFromChain(ctx context.Context, api *API) error {
	if opts == nil {
		return errors.New("TxOptions should not be nil, send an object")
	}

	if opts.HeadBlockID == nil || opts.ChainID == nil {
		info, err := api.cachedGetInfo(ctx)
		if err != nil {
			return err
		}

		if opts.HeadBlockID == nil {
			opts.HeadBlockID = info.HeadBlockID
		}
		if opts.ChainID == nil {
			opts.ChainID = info.ChainID
		}
	}

	return nil
}
