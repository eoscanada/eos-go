package eos

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"time"

	"io"

	"encoding/json"

	"io/ioutil"

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

type Transaction struct { // WARN: is a `variant` in C++, can be a SignedTransaction or a Transaction.
	TransactionHeader

	ContextFreeActions []*Action    `json:"context_free_actions"`
	Actions            []*Action    `json:"actions"`
	Extensions         []*Extension `json:"transaction_extensions"`
}

type Extension struct {
	Type uint16   `json:"type"`
	Data HexBytes `json:"data"`
}

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
	if tx.Extensions == nil {
		tx.Extensions = make([]*Extension, 0, 0)
	}

	tx.setRefBlock(info.HeadBlockID)

	/// TODO: configure somewhere the default time for transactions,
	/// etc.. add a `.Timeout` with that duration, default to 30
	/// seconds ?
	tx.Expiration = JSONTime{info.HeadBlockTime.Add(30 * time.Second)}
	//tx.DelaySec = 30

	return info.HeadBlockID, nil
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

func (s *SignedTransaction) String() string {

	data, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (tx *Transaction) ID() string {
	return "ID here" //todo
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

func (tx *SignedTransaction) estimateResources(opts TxOptions, maxcpu uint8, maxnet uint32) {
	// see programs/cleos/main.cpp for an estimation algo..
	if opts.MaxNetUsageWords != 0 {
		tx.MaxNetUsageWords = Varuint32(opts.MaxNetUsageWords)
	} else {
		tx.MaxNetUsageWords = Varuint32(maxnet)
	}

	if opts.MaxCPUUsageMS != 0 {
		tx.MaxCPUUsageMS = opts.MaxCPUUsageMS
	} else {
		tx.MaxCPUUsageMS = maxcpu
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

func (p *PackedTransaction) Unpack() (signedTx *SignedTransaction, err error) {
	var txReader io.Reader
	txReader = bytes.NewBuffer(p.PackedTransaction)

	var freeDataReader io.Reader
	freeDataReader = bytes.NewBuffer(p.PackedContextFreeData)

	switch p.Compression {
	case CompressionZlib:
		txReader, err = zlib.NewReader(txReader)
		if err != nil {
			return
		}

		freeDataReader, err = zlib.NewReader(freeDataReader)
		if err != nil {
			return
		}
	}

	data, err := ioutil.ReadAll(txReader)
	if err != nil {
		return
	}
	decoder := NewDecoder(data)

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

// TxOptions represents options you want to pass to the transaction
// you're sending.
type TxOptions struct {
	MaxNetUsageWords uint32
	Delay            time.Duration
	MaxCPUUsageMS    uint8 // If you want to override the CPU usage (in counts of 1024)
	//ExtraKCPUUsage uint32 // If you want to *add* some CPU usage to the estimated amount (in counts of 1024)
	Compress CompressionType
}
