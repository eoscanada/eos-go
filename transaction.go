package eos

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/eoscanada/eos-go/ecc"
)

type Transaction struct { // WARN: is a `variant` in C++, can be a SignedTransaction or a Transaction.
	Expiration     JSONTime `json:"expiration"`
	Region         uint16   `json:"region"`
	RefBlockNum    uint16   `json:"ref_block_num"`
	RefBlockPrefix uint32   `json:"ref_block_prefix"`

	MaxNetUsageWords Varuint32 `json:"max_net_usage_words"`
	MaxKCPUUsage     Varuint32 `json:"max_kcpu_usage"`
	DelaySec         Varuint32 `json:"delay_sec"` // number of secs to delay, making it cancellable for that duration

	// TODO: implement the estimators and write that in `.Fill()`.. for the transaction.

	ContextFreeActions []*Action `json:"context_free_actions"`
	Actions            []*Action `json:"actions"`
}

// 69c9c15a 0000 1400 62f95d45 b31d 904e 00 00 020000000000ea305500000040258ab2c2010000000000ea305500000000a8ed

func (tx *Transaction) Fill(api *API) ([]byte, error) {
	var info *InfoResp
	var err error

	api.lastGetInfoLock.Lock()
	if !api.lastGetInfoStamp.IsZero() && time.Now().Add(-1*time.Second).Before(api.lastGetInfoStamp) {
		info = api.lastGetInfo
	} else {
		info, err = api.GetInfo()
		if err != nil {
			return nil, err
		}
		api.lastGetInfoStamp = time.Now()
		api.lastGetInfo = info
	}
	api.lastGetInfoLock.Unlock()
	if err != nil {
		return nil, err
	}

	if tx.ContextFreeActions == nil {
		tx.ContextFreeActions = make([]*Action, 0, 0)
	}

	blockID, err := hex.DecodeString(info.HeadBlockID)
	if err != nil {
		return nil, fmt.Errorf("decode hex: %s", err)
	}

	tx.setRefBlock(blockID)

	/// TODO: configure somewhere the default time for transactions,
	/// etc.. add a `.Timeout` with that duration, default to 30
	/// seconds ?
	tx.Expiration = JSONTime{info.HeadBlockTime.Add(30 * time.Second)}

	return blockID, nil
}

func (tx *Transaction) setRefBlock(blockID []byte) {
	tx.RefBlockNum = uint16(binary.BigEndian.Uint32(blockID[:4]))
	tx.RefBlockPrefix = binary.LittleEndian.Uint32(blockID[8:16])
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

func (s *SignedTransaction) Pack(opts TxOptions) (*PackedTransaction, error) {
	rawtrx, err := MarshalBinary(s.Transaction)
	if err != nil {
		return nil, err
	}

	rawcfd, err := MarshalBinary(s.ContextFreeData)
	if err != nil {
		return nil, err
	}

	// Is it so ?
	if len(s.ContextFreeData) == 0 {
		rawcfd = []byte{}
	}

	switch opts.Compress {
	case CompressionZlib:
		var trx bytes.Buffer
		var cfd bytes.Buffer

		// Compress Trx
		writer, _ := zlib.NewWriterLevel(&trx, flate.BestCompression) // can only fail if invalid `level`..
		writer.Write(rawtrx)                                          // ignore error, could only bust memory
		rawtrx = trx.Bytes()

		// Compress ContextFreeData
		writer, _ = zlib.NewWriterLevel(&cfd, flate.BestCompression) // can only fail if invalid `level`..
		writer.Write(rawcfd)                                         // ignore errors, memory errors only
		rawcfd = cfd.Bytes()

	}

	packed := &PackedTransaction{
		Signatures:            s.Signatures,
		Compression:           opts.Compress,
		PackedContextFreeData: rawcfd,
		PackedTransaction:     rawtrx,
	}

	return packed, nil
}

func (tx *SignedTransaction) estimateResources(opts TxOptions, maxcpu, maxnet uint32) {
	// see programs/cleos/main.cpp for an estimation algo..
	if opts.MaxNetUsageWords != 0 {
		tx.MaxNetUsageWords = Varuint32(opts.MaxNetUsageWords)
	} else {
		tx.MaxNetUsageWords = Varuint32(maxnet)
	}

	if opts.MaxKCPUUsage != 0 {
		tx.MaxKCPUUsage = Varuint32(opts.MaxKCPUUsage)
	} else {
		tx.MaxKCPUUsage = Varuint32(maxcpu)
	}
}

// PackedTransaction represents a fully packed transaction, with
// signatures, and all. They circulate like that on the P2P net, and
// that's how they are stored.
type PackedTransaction struct {
	Signatures            []ecc.Signature `json:"signatures"`
	Compression           CompressionType `json:"compression"` // in C++, it's an enum, not sure how it Binary-marshals..
	PackedContextFreeData HexBytes        `json:"packed_context_free_data"`
	PackedTransaction     HexBytes        `json:"packed_trx"`
}

type DeferredTransaction struct {
	*Transaction

	SenderID   uint32      `json:"sender_id"`
	Sender     AccountName `json:"sender"`
	DelayUntil JSONTime    `json:"delay_until"`
}

// TxOptions represents options you want to pass to the transaction
// you're sending.
type TxOptions struct {
	MaxNetUsageWords uint32
	Delay            time.Duration
	MaxKCPUUsage     uint32 // If you want to override the CPU usage (in counts of 1024)
	//ExtraKCPUUsage uint32 // If you want to *add* some CPU usage to the estimated amount (in counts of 1024)
	Compress CompressionType
}
