package eos

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type Transaction struct { // WARN: is a `variant` in C++, can be a SignedTransaction or a Transaction.
	Expiration     JSONTime `json:"expiration,omitempty"`
	Region         uint16   `json:"region"`
	RefBlockNum    uint16   `json:"ref_block_num,omitempty"`
	RefBlockPrefix uint32   `json:"ref_block_prefix,omitempty"`

	NetUsageWords Varuint32 `json:"net_usage_words"`
	KCPUUsage     Varuint32 `json:"kcpu_usage"`
	DelaySec      Varuint32 `json:"delay_sec"` // number of secs to delay, making it cancellable for that duration

	// TODO: implement the estimators and write that in `.Fill()`.. for the transaction.

	ContextFreeActions []*Action `json:"context_free_actions,omitempty"`
	Actions            []*Action `json:"actions,omitempty"`
}

// 69c9c15a 0000 1400 62f95d45 b31d 904e 00 00 020000000000ea305500000040258ab2c2010000000000ea305500000000a8ed

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

	Signatures      []string   `json:"signatures"`
	ContextFreeData []HexBytes `json:"context_free_data"`

	packed *PackedTransaction
}

func NewSignedTransaction(tx *Transaction) *SignedTransaction {
	return &SignedTransaction{
		Transaction:     tx,
		Signatures:      make([]string, 0),
		ContextFreeData: make([]HexBytes, 0),
	}
}

func (s *SignedTransaction) Pack(opts TxOptions) error {
	return s.pack(opts, -1)
}

func (s *SignedTransaction) pack(opts TxOptions, packedLen int) error {
	data, err := MarshalBinary(s.Transaction)
	if err != nil {
		return err
	}

	switch opts.Compress {
	case CompressionZlib:
		var buf bytes.Buffer
		_, _ = zlib.NewWriter(&buf).Write(data)
		data = buf.Bytes()
	}

	packed := &PackedTransaction{
		Signatures:      s.Signatures,
		ContextFreeData: s.ContextFreeData,
		Compression:     opts.Compress,
		Data:            data,
	}

	s.packed = packed

	if err := s.estimateResources(opts); err != nil {
		return err
	}

	if packedLen == -1 {
		return s.pack(opts, len(data))
	}
	return nil
}

func (tx *SignedTransaction) Repack(opts TxOptions) error {
	tx.NetUsageWords = 0
	tx.packed = nil
	tx.KCPUUsage = 0

	return tx.Pack(opts)
}

func (tx *SignedTransaction) estimateResources(opts TxOptions) error {
	// see programs/cleos/main.cpp for an estimation algo..
	if tx.NetUsageWords == 0 {
		// for signatures
		base := 5 /* varint for sig count */ + len(tx.Signatures)*65 /* bytes per sig */

		// for this resources varints, not yet accounted for
		base += 4

		if opts.Compress == CompressionZlib {
			// for new data (see C++ code .. not sure why here)
			base += 252 // 4 + 252 = 256
		}

		for _, cfa := range tx.ContextFreeData {
			base += len(cfa)
		}
		if len(tx.ContextFreeData) != 0 {
			base += 7 // for alignment ?
		}

		base += len(tx.packed.Data)

		tx.NetUsageWords = Varuint32(base / 8) // because it's a count of 8-bytes words.
	}
	if opts.KCPUUsage != 0 {
		tx.KCPUUsage = Varuint32(opts.KCPUUsage)
	}
	if tx.KCPUUsage == 0 {
		base := 2048 /* for good measure :P */
		// Estimated per context-free actions usage..
		base += 10000 * len(tx.ContextFreeActions)
		tx.KCPUUsage = Varuint32(base)  // should divide by 1024 ?!
	}
	return nil
}

// PackedTransaction represents a fully packed transaction, with
// signatures, and all. They circulate like that on the P2P net, and
// that's how they are stored.
type PackedTransaction struct {
	Signatures      []string        `json:"signatures"`
	ContextFreeData []HexBytes      `json:"context_free_data"`
	Compression     CompressionType `json:"compression"` // in C++, it's an enum, not sure how it Binary-marshals..
	Data            HexBytes        `json:"data"`
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
	Delay     time.Duration
	KCPUUsage uint32 // If you want to override the CPU usage (in counts of 1024)
	//ExtraKCPUUsage uint32 // If you want to *add* some CPU usage to the estimated amount (in counts of 1024)
	Compress CompressionType
}
