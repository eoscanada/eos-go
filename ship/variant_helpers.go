package ship

import (
	"github.com/eoscanada/eos-go"
)

// Request
var RequestVariant = eos.NewVariantDefinition([]eos.VariantType{
	{"get_status_request_v0", (*GetStatusRequestV0)(nil)},
	{"get_blocks_request_v0", (*GetBlocksRequestV0)(nil)},
	{"get_blocks_ack_request_v0", (*GetBlocksAckRequestV0)(nil)},
})

type Request struct {
	eos.BaseVariant
}

func (r *Request) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, RequestVariant)
}

// Result
var ResultVariant = eos.NewVariantDefinition([]eos.VariantType{
	{"get_status_result_v0", (*GetStatusResultV0)(nil)},
	{"get_blocks_result_v0", (*GetBlocksResultV0)(nil)},
})

type Result struct {
	eos.BaseVariant
}

func (r *Result) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, ResultVariant)
}

// TransactionTrace
var TransactionTraceVariant = eos.NewVariantDefinition([]eos.VariantType{
	{"transaction_trace_v0", (*TransactionTraceV0)(nil)},
})

type TransactionTrace struct {
	eos.BaseVariant
}

type TransactionTraceArray struct {
	Elem []*TransactionTrace
}

func (t *TransactionTraceArray) AsTransactionTracesV0() (out []*TransactionTraceV0) {
	if t == nil || t.Elem == nil {
		return nil
	}
	for _, e := range t.Elem {
		switch v := e.Impl.(type) {
		case *TransactionTraceV0:
			out = append(out, v)

		default:
			panic("wrong type for conversion")
		}
	}
	return out
}

func (r TransactionTraceArray) MarshalBinary(enc *eos.Encoder) error {
	data, err := eos.MarshalBinary(r.Elem)
	if err != nil {
		return err
	}
	return enc.Encode(data)
}

func (r *TransactionTraceArray) UnmarshalBinary(decoder *eos.Decoder) error {
	data, err := decoder.ReadByteArray()
	if err != nil {
		return err
	}
	return eos.UnmarshalBinary(data, &r.Elem)
}

func (r *TransactionTrace) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, TransactionTraceVariant)
}

// ActionTrace
var ActionTraceVariant = eos.NewVariantDefinition([]eos.VariantType{
	{"action_trace_v0", (*ActionTraceV0)(nil)},
	{"action_trace_v1", (*ActionTraceV1)(nil)},
})

type ActionTrace struct {
	eos.BaseVariant
}

func (r *ActionTrace) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, ActionTraceVariant)
}

// PartialTransaction
var PartialTransactionVariant = eos.NewVariantDefinition([]eos.VariantType{
	{"partial_transaction_v0", (*PartialTransactionV0)(nil)},
})

type PartialTransaction struct {
	eos.BaseVariant
}

func (r *PartialTransaction) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, PartialTransactionVariant)
}

// TableDelta
var TableDeltaVariant = eos.NewVariantDefinition([]eos.VariantType{
	{"table_delta_v0", (*TableDeltaV0)(nil)},
})

type TableDelta struct {
	eos.BaseVariant
}

func (d *TableDelta) UnmarshalBinary(decoder *eos.Decoder) error {
	return d.BaseVariant.UnmarshalBinaryVariant(decoder, TableDeltaVariant)
}

type TableDeltaArray struct {
	Elem []*TableDelta
}

func (d TableDeltaArray) MarshalBinary(enc *eos.Encoder) error {
	data, err := eos.MarshalBinary(d.Elem)
	if err != nil {
		return err
	}
	return enc.Encode(data)
}

func (d *TableDeltaArray) UnmarshalBinary(decoder *eos.Decoder) error {
	data, err := decoder.ReadByteArray()
	if err != nil {
		return err
	}
	return eos.UnmarshalBinary(data, &d.Elem)
}

func (t *TableDeltaArray) AsTableDeltasV0() (out []*TableDeltaV0) {
	if t == nil || t.Elem == nil {
		return nil
	}
	for _, e := range t.Elem {
		switch v := e.Impl.(type) {
		case *TableDeltaV0:
			out = append(out, v)

		default:
			panic("wrong type for conversion")
		}
	}
	return out
}

// Transaction
var TransactionVariant = eos.NewVariantDefinition([]eos.VariantType{
	{"transaction_id", (*eos.Checksum256)(nil)},
	{"packed_transaction", (*eos.PackedTransaction)(nil)},
})

type Transaction struct {
	eos.BaseVariant
}

func (d *Transaction) UnmarshalBinary(decoder *eos.Decoder) error {
	return d.BaseVariant.UnmarshalBinaryVariant(decoder, TransactionVariant)
}

// ActionReceipt
var ActionReceiptVariant = eos.NewVariantDefinition([]eos.VariantType{
	{"action_receipt_v0", (*ActionReceiptV0)(nil)},
})

type ActionReceipt struct {
	eos.BaseVariant
}

func (r *ActionReceipt) UnmarshalBinary(decoder *eos.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, ActionReceiptVariant)
}
