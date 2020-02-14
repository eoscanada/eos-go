package ship

import "github.com/eoscanada/eos-go"

// Request

type Request struct {
	eos.BaseVariant
}

const (
	GetStatusRequestV0Type = iota
	GetBlocksRequestV0Type
	GetBlocksAckRequestV0Type
)

var requestVariantFactoryImplMap = map[uint32]eos.VariantImplFactory{
	GetStatusRequestV0Type:    func() interface{} { return new(GetStatusRequestV0) },
	GetBlocksRequestV0Type:    func() interface{} { return new(GetBlocksRequestV0) },
	GetBlocksAckRequestV0Type: func() interface{} { return new(GetBlocksAckRequestV0) },
}

func (r *Request) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, requestVariantFactoryImplMap)
}

// Result
const (
	GetStatusResultV0Type = iota
	GetBlocksResultV0Type
)

var resultVariantFactoryImplMap = map[uint32]eos.VariantImplFactory{
	GetStatusResultV0Type: func() interface{} { return new(GetStatusResultV0) },
	GetBlocksResultV0Type: func() interface{} { return new(GetBlocksResultV0) },
}

type Result struct {
	eos.BaseVariant
}

func (r *Result) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, resultVariantFactoryImplMap)
}

// TransactionTrace
const (
	TransactionTraceV0Type = iota
)

var TransactionTraceVariantFactoryImplMap = map[uint32]eos.VariantImplFactory{
	TransactionTraceV0Type: func() interface{} { return new(TransactionTraceV0) },
}

type TransactionTraceArray struct {
	Elem []*TransactionTrace
}

func (r *TransactionTraceArray) UnmarshalBinary(decoder *eos.Decoder) error {
	data, err := decoder.ReadByteArray()
	if err != nil {
		return err
	}
	return eos.UnmarshalBinary(data, &r.Elem)
}

type TransactionTrace struct {
	eos.BaseVariant
}

func (r *TransactionTrace) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, TransactionTraceVariantFactoryImplMap)
}

// ActionTrace
const (
	ActionTraceV0Type = iota
)

var ActionTraceVariantFactoryImplMap = map[uint32]eos.VariantImplFactory{
	ActionTraceV0Type: func() interface{} { return new(ActionTraceV0) },
}

type ActionTrace struct {
	eos.BaseVariant
}

func (r *ActionTrace) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, ActionTraceVariantFactoryImplMap)
}

// PartialTransaction
const (
	PartialTransactionV0Type = iota
)

var PartialTransactionVariantFactoryImplMap = map[uint32]eos.VariantImplFactory{
	PartialTransactionV0Type: func() interface{} { return new(PartialTransactionV0) },
}

type PartialTransaction struct {
	eos.BaseVariant
}

func (r *PartialTransaction) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, PartialTransactionVariantFactoryImplMap)
}

// TableDelta
const (
	TableDeltaV0Type = iota
)

var TableDeltaVariantFactoryImplMap = map[uint32]eos.VariantImplFactory{
	TableDeltaV0Type: func() interface{} { return new(TableDeltaV0) },
}

type TableDelta struct {
	eos.BaseVariant
}

func (d TableDelta) UnmarshalBinary(decoder *eos.Decoder) error {
	return d.BaseVariant.UnmarshalBinaryVariant(decoder, TableDeltaVariantFactoryImplMap)
}

type TableDeltaArray struct {
	Elem []*TableDelta
}

func (d TableDeltaArray) UnmarshalBinary(decoder *eos.Decoder) error {
	data, err := decoder.ReadByteArray()
	if err != nil {
		return err
	}
	return eos.UnmarshalBinary(data, &d.Elem)
}

// TransactionVariant
const (
	TransactionIDType = iota
	PackedTransactionType
)

var TransactionVariantFactoryImplMap = map[uint32]eos.VariantImplFactory{
	TransactionIDType:     func() interface{} { return new(TransactionID) },
	PackedTransactionType: func() interface{} { return new(eos.PackedTransaction) },
}

type TransactionVariant struct {
	eos.BaseVariant
}

func (d *TransactionVariant) UnmarshalBinary(decoder *eos.Decoder) error {
	return d.BaseVariant.UnmarshalBinaryVariant(decoder, TransactionVariantFactoryImplMap)
}

// ActionReceipt
const (
	ActionReceiptV0Type = iota
)

var ActionReceiptVariantFactoryImplMap = map[uint32]eos.VariantImplFactory{
	ActionReceiptV0Type: func() interface{} { return new(ActionReceiptV0) },
}

type ActionReceipt struct {
	eos.BaseVariant
}

func (r *ActionReceipt) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, ActionReceiptVariantFactoryImplMap)
}
