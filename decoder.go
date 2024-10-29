package eos

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/eoscanada/eos-go/ecc"
	"go.uber.org/zap"
)

// UnmarshalerBinary is the interface implemented by types
// that can unmarshal an EOSIO binary description of themselves.
//
// **Warning** This is experimental, exposed only for internal usage for now.
type UnmarshalerBinary interface {
	UnmarshalBinary(decoder *Decoder) error
}

var TypeSize = struct {
	Bool int
	Byte int

	Int8  int
	Int16 int

	Uint8   int
	Uint16  int
	Uint32  int
	Uint64  int
	Uint128 int

	Float32 int
	Float64 int

	Checksum160 int
	Checksum256 int
	Checksum512 int

	PublicKey int
	Signature int

	Tstamp         int
	BlockTimestamp int

	CurrencyName int

	// Deprecated Fields

	// Deprecated: use Uint8 instead
	UInt8 int
	// Deprecated: use Uint16 instead
	UInt16 int
	// Deprecated: use Uint32 instead
	UInt32 int
	// Deprecated: use Uint64 instead
	UInt64 int
	// Deprecated: use Uint128 instead
	UInt128 int
}{
	Byte: 1,
	Bool: 1,

	Int8:  1,
	Int16: 2,

	Uint8:   1,
	Uint16:  2,
	Uint32:  4,
	Uint64:  8,
	Uint128: 16,

	Float32: 4,
	Float64: 8,

	Checksum160: 20,
	Checksum256: 32,
	Checksum512: 64,

	PublicKey: 34,
	Signature: 66,

	Tstamp:         8,
	BlockTimestamp: 4,

	CurrencyName: 7,
}

func init() {
	// Deprecated fields initialization
	TypeSize.UInt8 = TypeSize.Uint8
	TypeSize.UInt16 = TypeSize.Uint16
	TypeSize.UInt32 = TypeSize.Uint32
	TypeSize.UInt64 = TypeSize.Uint64
	TypeSize.UInt128 = TypeSize.Uint128
}

var RegisteredActions = map[AccountName]map[ActionName]reflect.Type{}

// Registers Action objects..
func RegisterAction(accountName AccountName, actionName ActionName, obj interface{}) {
	// TODO: lock or som'th.. unless we never call after boot time..
	if RegisteredActions[accountName] == nil {
		RegisteredActions[accountName] = make(map[ActionName]reflect.Type)
	}
	RegisteredActions[accountName][actionName] = reflect.TypeOf(obj)
}

// Decoder implements the EOS unpacking, similar to FC_BUFFER
type Decoder struct {
	data             []byte
	pos              int
	decodeP2PMessage bool
	decodeActions    bool
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		data:             data,
		decodeP2PMessage: true,
		decodeActions:    true,
	}
}

func (d *Decoder) DecodeP2PMessage(decode bool) {
	d.decodeP2PMessage = decode
}

func (d *Decoder) DecodeActions(decode bool) {
	d.decodeActions = decode
}

type DecodeOption = interface{}

type optionalFieldType bool

const OptionalField optionalFieldType = true

func (d *Decoder) LastPos() int {
	return d.pos
}

func (d *Decoder) Decode(v interface{}, options ...DecodeOption) (err error) {
	optionalField := false
	for _, option := range options {
		if _, isOptionalField := option.(optionalFieldType); isOptionalField {
			optionalField = true
		}
	}

	rv := reflect.Indirect(reflect.ValueOf(v))
	if !rv.CanAddr() {
		return fmt.Errorf("can only decode to pointer type, got %T", v)
	}
	t := rv.Type()

	if tracer.Enabled() {
		zlog.Debug("decode type", typeField("type", v), zap.Bool("optional", optionalField))
	}

	if optionalField {
		var isPresent byte
		if d.hasRemaining() {
			isPresent, err = d.ReadByte()
			if err != nil {
				err = fmt.Errorf("decode: %t isPresent, %s", v, err)
				return
			}
		}

		if isPresent == 0 {
			if tracer.Enabled() {
				zlog.Debug("skipping optional", typeField("type", v))
			}

			rv.Set(reflect.Zero(t))
			return
		}
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		newRV := reflect.New(t)
		rv.Set(newRV)

		// At this point, `newRV` is a pointer to our target type, we need to check here because
		// after that, when `reflect.Indirect` is used, we get a `**<Type>` instead of a `*<Type>`
		// which breaks the interface checking.
		//
		// Ultimetaly, I think this could should be re-written, I don't think the `**<Type>` is necessary.
		if u, ok := newRV.Interface().(UnmarshalerBinary); ok {
			if tracer.Enabled() {
				zlog.Debug("using UnmarshalBinary method to decode type", typeField("type", v))
			}
			return u.UnmarshalBinary(d)
		}

		rv = reflect.Indirect(newRV)
	} else {
		// We check if `v` directly is `UnmarshalerBinary` this is to overcome our bad code that
		// has problem dealing with non-pointer type, which should still be possible here, by allocating
		// the empty value for it can then unmarshalling using the address of it. See comment above about
		// `newRV` being turned into `**<Type>`.
		//
		// We should re-code all the logic to determine the type and indirection using Golang `json` package
		// logic. See here: https://github.com/golang/go/blob/54697702e435bddb69c0b76b25b3209c78d2120a/src/encoding/json/decode.go#L439
		if u, ok := v.(UnmarshalerBinary); ok {
			if tracer.Enabled() {
				zlog.Debug("using UnmarshalBinary method to decode type", typeField("type", v))
			}
			return u.UnmarshalBinary(d)
		}
	}

	switch v.(type) {
	case *string:
		s, e := d.ReadString()
		if e != nil {
			err = e
			return
		}
		rv.SetString(s)
		return
	case *Name, *AccountName, *PermissionName, *ActionName, *TableName, *ScopeName:
		var n uint64
		n, err = d.ReadUint64()
		name := NameToString(n)
		if tracer.Enabled() {
			zlog.Debug("read name", zap.String("name", name))
		}
		rv.SetString(name)
		return

	case *byte, *P2PMessageType, *TransactionStatus, *CompressionType, *IDListMode, *GoAwayReason:
		var n byte
		n, err = d.ReadByte()
		rv.SetUint(uint64(n))
		return
	case *int8:
		var n int8
		n, err = d.ReadInt8()
		rv.SetInt(int64(n))
		return
	case *int16:
		var n int16
		n, err = d.ReadInt16()
		rv.SetInt(int64(n))
		return
	case *int32:
		var n int32
		n, err = d.ReadInt32()
		rv.SetInt(int64(n))
		return
	case *int64:
		var n int64
		n, err = d.ReadInt64()
		rv.SetInt(int64(n))
		return
	case *Int64:
		var n int64
		n, err = d.ReadInt64()
		rv.SetInt(int64(n))
		return

	// This is so hackish, doing it right now, but the decoder needs to handle those
	// case (a struct field that is itself a pointer) by itself.
	case **Uint64:
		var n uint64
		n, err = d.ReadUint64()
		if err == nil {
			rv.Set(reflect.ValueOf((Uint64)(n)))
		}

		return
	case *Uint64:
		var n uint64
		n, err = d.ReadUint64()
		rv.SetUint(uint64(n))
		return
	case *Float64:
		var n float64
		n, err = d.ReadFloat64()
		rv.SetFloat(n)
		return
	case *Uint128:
		var n Uint128
		n, err = d.ReadUint128("uint128")
		rv.Set(reflect.ValueOf(n))
		return
	case *Int128:
		var n Uint128
		n, err = d.ReadUint128("int128")
		rv.Set(reflect.ValueOf(Int128(n)))
		return
	case *Float128:
		var n Uint128
		n, err = d.ReadUint128("float128")
		rv.Set(reflect.ValueOf(Float128(n)))
		return
	case *uint16:
		var n uint16
		n, err = d.ReadUint16()
		rv.SetUint(uint64(n))
		return
	case *uint32:
		var n uint32
		n, err = d.ReadUint32()
		rv.SetUint(uint64(n))
		return
	case *uint64:
		var n uint64
		n, err = d.ReadUint64()
		rv.SetUint(n)
		return
	case *Varuint32:
		var r uint64
		r, err = d.ReadUvarint64()
		rv.SetUint(r)
		return
	case *bool:
		var r bool
		r, err = d.ReadBool()
		rv.SetBool(r)
		return
	case *Bool:
		var r bool
		r, err = d.ReadBool()
		rv.SetBool(r)
		return
	case *HexBytes:
		var data []byte
		data, err = d.ReadByteArray()
		rv.SetBytes(data)
		return
	case *[]byte:
		var data []byte
		data, err = d.ReadByteArray()
		rv.SetBytes(data)
		return
	case *Checksum256:
		var s Checksum256
		s, err = d.ReadChecksum256()
		rv.SetBytes(s)
		return
	case *ecc.PublicKey:
		var p ecc.PublicKey
		p, err = d.ReadPublicKey()
		rv.Set(reflect.ValueOf(p))
		return
	case *ecc.Signature:
		var s ecc.Signature
		s, err = d.ReadSignature()
		rv.Set(reflect.ValueOf(s))
		return
	case *Tstamp:
		var ts Tstamp
		ts, err = d.ReadTstamp()
		rv.Set(reflect.ValueOf(ts))
		return
	case *TimePoint:
		var tp TimePoint
		tp, err = d.ReadTimePoint()
		rv.Set(reflect.ValueOf(tp))
		return
	case *TimePointSec:
		var tp TimePointSec
		tp, err = d.ReadTimePointSec()
		rv.Set(reflect.ValueOf(tp))
		return
	case *BlockTimestamp:
		var bt BlockTimestamp
		bt, err = d.ReadBlockTimestamp()
		rv.Set(reflect.ValueOf(bt))
		return
	case *JSONTime:
		var jt JSONTime
		jt, err = d.ReadJSONTime()
		rv.Set(reflect.ValueOf(jt))
		return
	case *CurrencyName:
		var cur CurrencyName
		cur, err = d.ReadCurrencyName()
		rv.Set(reflect.ValueOf(cur))
		return
	case *Symbol:
		var symbol *Symbol
		symbol, err = d.ReadSymbol()
		rv.Set(reflect.ValueOf(*symbol))
		return
	case *SymbolCode:
		var sc SymbolCode
		sc, err = d.ReadSymbolCode()
		rv.Set(reflect.ValueOf(sc))
		return
	case *Asset:
		var asset Asset
		asset, err = d.ReadAsset()
		rv.Set(reflect.ValueOf(asset))
		return
	case *TransactionWithID:
		t, e := d.ReadByte()
		if err != nil {
			err = fmt.Errorf("failed to read TransactionWithID type byte: %w", e)
			return
		}

		if tracer.Enabled() {
			zlog.Debug("type byte value", zap.Uint8("val", t))
		}

		if t == 0 {
			id, e := d.ReadChecksum256()
			if err != nil {
				err = fmt.Errorf("failed to read TransactionWithID id: %w", e)
				return
			}

			trx := TransactionWithID{ID: id}
			rv.Set(reflect.ValueOf(trx))
			return nil

		} else {
			packedTrx := &PackedTransaction{}
			if err := d.Decode(packedTrx); err != nil {
				return fmt.Errorf("packed transaction: %w", err)
			}

			id, err := packedTrx.ID()
			if err != nil {
				return fmt.Errorf("packed transaction id: %w", err)
			}

			trx := TransactionWithID{ID: id, Packed: packedTrx}
			rv.Set(reflect.ValueOf(trx))
			return nil
		}

	case **Action:
		err = d.decodeStruct(v, t, rv)
		if err != nil {
			return
		}
		action := rv.Interface().(Action)

		if d.decodeActions {
			err = d.ReadActionData(&action)
		}

		rv.Set(reflect.ValueOf(action))
		return

	case *Packet:

		envelope, e := d.ReadP2PMessageEnvelope()
		if e != nil {
			err = fmt.Errorf("decode, %s", e)
			return
		}

		if d.decodeP2PMessage {
			attr, ok := envelope.Type.reflectTypes()
			if !ok {
				return fmt.Errorf("decode, unknown p2p message type [%d]", envelope.Type)
			}
			msg := reflect.New(attr.ReflectType)
			subDecoder := NewDecoder(envelope.Payload)

			err = subDecoder.Decode(msg.Interface())

			decoded := msg.Interface().(P2PMessage)
			envelope.P2PMessage = decoded
		}

		rv.Set(reflect.ValueOf(*envelope))

		return
	}

	switch t.Kind() {
	case reflect.Array:
		if tracer.Enabled() {
			zlog.Debug("reading array")
		}
		len := t.Len()
		for i := 0; i < int(len); i++ {
			if err = d.Decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		var l uint64
		if l, err = d.ReadUvarint64(); err != nil {
			return
		}
		if tracer.Enabled() {
			zlog.Debug("reading slice", zap.Uint64("len", l), typeField("type", v))
		}
		rv.Set(reflect.MakeSlice(t, int(l), int(l)))
		for i := 0; i < int(l); i++ {
			if err = d.Decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}

	case reflect.Struct:

		err = d.decodeStruct(v, t, rv)
		if err != nil {
			return
		}

	default:
		return errors.New("decode, unsupported type " + t.String())
	}

	return
}

func (d *Decoder) decodeStruct(v interface{}, t reflect.Type, rv reflect.Value) (err error) {
	l := rv.NumField()

	seenBinaryExtensionField := false
	for i := 0; i < l; i++ {
		structField := t.Field(i)
		tag := structField.Tag.Get("eos")
		if tag == "-" {
			continue
		}

		if tag != "binary_extension" && seenBinaryExtensionField {
			panic(fmt.Sprintf("the `eos: \"binary_extension\"` tags must be packed together at the end of struct fields, problematic field %s", structField.Name))
		}

		if tag == "binary_extension" {
			seenBinaryExtensionField = true

			// FIXME: This works only if what is in `d.data` is the actual full data buffer that
			//        needs to be decoded. If there is for example two structs in the buffer, this
			//        will not work as we would continue into the next struct.
			//
			//        But at the same time, does it make sense otherwise? What would be the inference
			//        rule in the case of extra bytes available? Continue decoding and revert if it's
			//        not working? But how to detect valid errors?
			if len(d.data[d.pos:]) <= 0 {
				continue
			}
		}

		if v := rv.Field(i); v.CanSet() && structField.Name != "_" {
			var options []DecodeOption
			if tag == "optional" {
				options = append(options, OptionalField)
			}

			value := v.Addr().Interface()

			if tracer.Enabled() {
				zlog.Debug("struct field", typeField(structField.Name, value), zap.String("tag", tag))
			}

			if err = d.Decode(value, options...); err != nil {
				return
			}
		}
	}
	return
}

var ErrVarIntBufferSize = errors.New("varint: invalid buffer size")

func (d *Decoder) ReadUvarint64() (uint64, error) {
	l, read := binary.Uvarint(d.data[d.pos:])
	if read <= 0 {
		return l, ErrVarIntBufferSize
	}
	if tracer.Enabled() {
		zlog.Debug("read uvarint64", zap.Uint64("val", l))
	}
	d.pos += read
	return l, nil
}
func (d *Decoder) ReadVarint64() (out int64, err error) {
	l, read := binary.Varint(d.data[d.pos:])
	if read <= 0 {
		return l, ErrVarIntBufferSize
	}
	if tracer.Enabled() {
		zlog.Debug("read varint", zap.Int64("val", l))
	}
	d.pos += read
	return l, nil
}

func (d *Decoder) ReadVarint32() (out int32, err error) {
	n, err := d.ReadVarint64()
	if err != nil {
		return out, err
	}
	out = int32(n)
	if tracer.Enabled() {
		zlog.Debug("read varint32", zap.Int32("val", out))
	}
	return
}
func (d *Decoder) ReadUvarint32() (out uint32, err error) {

	n, err := d.ReadUvarint64()
	if err != nil {
		return out, err
	}
	out = uint32(n)
	if tracer.Enabled() {
		zlog.Debug("read uvarint32", zap.Uint32("val", out))
	}
	return
}

func (d *Decoder) ReadByteArray() (out []byte, err error) {

	l, err := d.ReadUvarint64()
	if err != nil {
		return nil, err
	}

	if l > math.MaxInt || d.pos > math.MaxInt-int(l) {
		return nil, errors.New("byte array: varlen is overflowing int")
	}

	if len(d.data) < d.pos+int(l) {
		return nil, fmt.Errorf("byte array: varlen=%d, missing %d bytes", l, d.pos+int(l)-len(d.data))
	}

	out = d.data[d.pos : d.pos+int(l)]
	d.pos += int(l)
	if tracer.Enabled() {
		zlog.Debug("read byte array", zap.Stringer("hex", HexBytes(out)))
	}
	return
}

func (d *Decoder) ReadByte() (out byte, err error) {
	if d.remaining() < TypeSize.Byte {
		err = fmt.Errorf("required [1] byte, remaining [%d]", d.remaining())
		return
	}

	out = d.data[d.pos]
	d.pos++
	if tracer.Enabled() {
		zlog.Debug("read byte", zap.Uint8("byte", out))
	}
	return
}

func (d *Decoder) ReadBool() (out bool, err error) {
	if d.remaining() < TypeSize.Bool {
		err = fmt.Errorf("bool required [%d] byte, remaining [%d]", TypeSize.Bool, d.remaining())
		return
	}

	b, err := d.ReadByte()

	if err != nil {
		err = fmt.Errorf("readBool, %s", err)
	}
	out = b != 0
	if tracer.Enabled() {
		zlog.Debug("read bool", zap.Bool("val", out))
	}
	return

}

func (d *Decoder) ReadUint8() (out uint8, err error) {
	out, err = d.ReadByte()
	return
}

// Deprecated: Use `ReadUint8` (with a lower case `i`) instead
func (d *Decoder) ReadUInt8() (out uint8, err error) {
	return d.ReadUint8()
}

func (d *Decoder) ReadInt8() (out int8, err error) {
	b, err := d.ReadByte()
	out = int8(b)
	if tracer.Enabled() {
		zlog.Debug("read int8", zap.Int8("val", out))
	}
	return
}

func (d *Decoder) ReadUint16() (out uint16, err error) {
	if d.remaining() < TypeSize.Uint16 {
		err = fmt.Errorf("uint16 required [%d] bytes, remaining [%d]", TypeSize.Uint16, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint16(d.data[d.pos:])
	d.pos += TypeSize.Uint16
	if tracer.Enabled() {
		zlog.Debug("read uint16", zap.Uint16("val", out))
	}
	return
}

func (d *Decoder) ReadInt16() (out int16, err error) {
	n, err := d.ReadUint16()
	out = int16(n)
	if tracer.Enabled() {
		zlog.Debug("read int16", zap.Int16("val", out))
	}
	return
}
func (d *Decoder) ReadInt64() (out int64, err error) {
	n, err := d.ReadUint64()
	out = int64(n)
	if tracer.Enabled() {
		zlog.Debug("read int64", zap.Int64("val", out))
	}
	return
}

func (d *Decoder) ReadUint32() (out uint32, err error) {
	if d.remaining() < TypeSize.Uint32 {
		err = fmt.Errorf("uint32 required [%d] bytes, remaining [%d]", TypeSize.Uint32, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint32(d.data[d.pos:])
	d.pos += TypeSize.Uint32
	if tracer.Enabled() {
		zlog.Debug("read uint32", zap.Uint32("val", out))
	}
	return
}
func (d *Decoder) ReadInt32() (out int32, err error) {
	n, err := d.ReadUint32()
	out = int32(n)
	if tracer.Enabled() {
		zlog.Debug("read int32", zap.Int32("val", out))
	}
	return
}

func (d *Decoder) ReadUint64() (out uint64, err error) {
	if d.remaining() < TypeSize.Uint64 {
		err = fmt.Errorf("uint64 required [%d] bytes, remaining [%d]", TypeSize.Uint64, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.Uint64]
	out = binary.LittleEndian.Uint64(data)
	d.pos += TypeSize.Uint64
	if tracer.Enabled() {
		zlog.Debug("read uint64", zap.Uint64("val", out), zap.Stringer("hex", HexBytes(data)))
	}
	return
}

func (d *Decoder) ReadInt128() (out Int128, err error) {
	v, err := d.ReadUint128("int128")
	if err != nil {
		return
	}

	return Int128(v), nil
}

func (d *Decoder) ReadUint128(typeName string) (out Uint128, err error) {
	if d.remaining() < TypeSize.Uint128 {
		err = fmt.Errorf("%s required [%d] bytes, remaining [%d]", typeName, TypeSize.Uint128, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.Uint128]
	out.Lo = binary.LittleEndian.Uint64(data)
	out.Hi = binary.LittleEndian.Uint64(data[8:])

	d.pos += TypeSize.Uint128
	if tracer.Enabled() {
		zlog.Debug("read uint128", zap.Stringer("hex", out), zap.Uint64("hi", out.Hi), zap.Uint64("lo", out.Lo))
	}
	return
}

func (d *Decoder) ReadFloat32() (out float32, err error) {
	if d.remaining() < TypeSize.Float32 {
		err = fmt.Errorf("float32 required [%d] bytes, remaining [%d]", TypeSize.Float32, d.remaining())
		return
	}

	n := binary.LittleEndian.Uint32(d.data[d.pos:])
	out = math.Float32frombits(n)
	d.pos += TypeSize.Float32
	if tracer.Enabled() {
		zlog.Debug("read float32", zap.Float32("val", out))
	}
	return
}

func (d *Decoder) ReadNodeosFloat32() (out float32, err error) {
	if d.remaining() < TypeSize.Float32 {
		err = fmt.Errorf("float32 required [%d] bytes, remaining [%d]", TypeSize.Float32, d.remaining())
		return
	}

	n := binary.LittleEndian.Uint32(d.data[d.pos:])
	out = math.Float32frombits(n)
	d.pos += TypeSize.Float32
	if tracer.Enabled() {
		zlog.Debug("read float32", zap.Float32("val", out))
	}
	return
}

func (d *Decoder) ReadFloat64() (out float64, err error) {
	if d.remaining() < TypeSize.Float64 {
		err = fmt.Errorf("float64 required [%d] bytes, remaining [%d]", TypeSize.Float64, d.remaining())
		return
	}

	n := binary.LittleEndian.Uint64(d.data[d.pos:])
	out = math.Float64frombits(n)
	d.pos += TypeSize.Float64
	if tracer.Enabled() {
		zlog.Debug("read Float64", zap.Float64("val", float64(out)))
	}
	return
}

func fixUtf(r rune) rune {
	if r == utf8.RuneError {
		return '�'
	}
	return r
}
func (d *Decoder) SafeReadUTF8String() (out string, err error) {
	data, err := d.ReadByteArray()
	out = strings.Map(fixUtf, string(data))
	if tracer.Enabled() {
		zlog.Debug("read safe UTF8 string", zap.String("val", out))
	}
	return
}

func (d *Decoder) ReadString() (out string, err error) {
	data, err := d.ReadByteArray()
	out = string(data)
	if tracer.Enabled() {
		zlog.Debug("read string", zap.String("val", out))
	}
	return
}

func (d *Decoder) ReadChecksum160() (out Checksum160, err error) {
	if d.remaining() < TypeSize.Checksum160 {
		err = fmt.Errorf("checksum 160 required [%d] bytes, remaining [%d]", TypeSize.Checksum160, d.remaining())
		return
	}

	out = make(Checksum160, TypeSize.Checksum160)
	copy(out, d.data[d.pos:d.pos+TypeSize.Checksum160])
	d.pos += TypeSize.Checksum160
	if tracer.Enabled() {
		zlog.Debug("read checksum160", zap.Stringer("hex", HexBytes(out)))
	}
	return
}

func (d *Decoder) ReadChecksum256() (out Checksum256, err error) {
	if d.remaining() < TypeSize.Checksum256 {
		err = fmt.Errorf("checksum 256 required [%d] bytes, remaining [%d]", TypeSize.Checksum256, d.remaining())
		return
	}

	out = make(Checksum256, TypeSize.Checksum256)
	copy(out, d.data[d.pos:d.pos+TypeSize.Checksum256])
	d.pos += TypeSize.Checksum256
	if tracer.Enabled() {
		zlog.Debug("read checksum256", zap.Stringer("hex", HexBytes(out)))
	}
	return
}

func (d *Decoder) ReadChecksum512() (out Checksum512, err error) {
	if d.remaining() < TypeSize.Checksum512 {
		err = fmt.Errorf("checksum 512 required [%d] bytes, remaining [%d]", TypeSize.Checksum512, d.remaining())
		return
	}

	out = make(Checksum512, TypeSize.Checksum512)
	copy(out, d.data[d.pos:d.pos+TypeSize.Checksum512])
	d.pos += TypeSize.Checksum512
	if tracer.Enabled() {
		zlog.Debug("read checksum512", zap.Stringer("hex", HexBytes(out)))
	}
	return
}

func (d *Decoder) ReadPublicKey() (out ecc.PublicKey, err error) {
	typeID, err := d.ReadUint8()
	if err != nil {
		return out, fmt.Errorf("unable to read public key type: %w", err)
	}

	curveID := ecc.CurveID(typeID)
	var keyMaterial []byte

	if curveID == ecc.CurveK1 || curveID == ecc.CurveR1 {
		// Minus 1 because we already read the curveID which is 1 out of the 34 bytes of a full "legacy" PublicKey
		keyMaterial, err = d.readPublicKeyMaterial(curveID, TypeSize.PublicKey-1)
	} else if curveID == ecc.CurveWA {
		keyMaterial, err = d.readWAPublicKeyMaterial()
	} else {
		err = fmt.Errorf("unsupported curve ID %d (%s)", uint8(curveID), curveID)
	}

	if err != nil {
		return out, fmt.Errorf("unable to read public key material for curve %s: %w", curveID, err)
	}

	data := append([]byte{byte(curveID)}, keyMaterial...)
	out, err = ecc.NewPublicKeyFromData(data)
	if err != nil {
		return out, fmt.Errorf("new public key from data: %w", err)
	}

	if tracer.Enabled() {
		zlog.Debug("read public key", zap.Stringer("pubkey", out))
	}

	return
}

func (d *Decoder) readPublicKeyMaterial(curveID ecc.CurveID, keyMaterialSize int) (out []byte, err error) {
	if d.remaining() < keyMaterialSize {
		err = fmt.Errorf("publicKey %s key material requires [%d] bytes, remaining [%d]", curveID, keyMaterialSize, d.remaining())
		return
	}

	out = make([]byte, keyMaterialSize)
	copy(out, d.data[d.pos:d.pos+keyMaterialSize])
	d.pos += keyMaterialSize

	return
}

func (d *Decoder) readWAPublicKeyMaterial() (out []byte, err error) {
	begin := d.pos
	if d.remaining() < 35 {
		err = fmt.Errorf("publicKey WA key material requires at least [35] bytes, remaining [%d]", d.remaining())
		return
	}

	d.pos += 34
	remainderDataSize, err := d.ReadUvarint32()
	if err != nil {
		return out, fmt.Errorf("unable to read public key WA key material size: %w", err)
	}

	if d.remaining() < int(remainderDataSize) {
		err = fmt.Errorf("publicKey WA remainder key material requires [%d] bytes, remaining [%d]", remainderDataSize, d.remaining())
		return
	}

	d.pos += int(remainderDataSize)
	keyMaterialSize := d.pos - begin

	out = make([]byte, keyMaterialSize)
	copy(out, d.data[begin:begin+keyMaterialSize])

	return
}

func (d *Decoder) ReadSignature() (out ecc.Signature, err error) {
	typeID, err := d.ReadUint8()
	if err != nil {
		return out, fmt.Errorf("unable to read signature type: %w", err)
	}

	curveID := ecc.CurveID(typeID)
	if tracer.Enabled() {
		zlog.Debug("read signature curve id", zap.Stringer("curve", curveID))
	}

	var data []byte
	if curveID == ecc.CurveK1 || curveID == ecc.CurveR1 {
		// Minus 1 because we already read the curveID which is 1 out of the 34 bytes of a full "legacy" PublicKey
		if d.remaining() < TypeSize.Signature-1 {
			return out, fmt.Errorf("signature required [%d] bytes, remaining [%d]", TypeSize.Signature-1, d.remaining())
		}

		data = make([]byte, 66)
		data[0] = byte(curveID)
		copy(data[1:], d.data[d.pos:d.pos+TypeSize.Signature-1])
		if tracer.Enabled() {
			zlog.Debug("read signature data", zap.Stringer("data", HexBytes(data)))
		}

		d.pos += TypeSize.Signature - 1
	} else if curveID == ecc.CurveWA {
		data, err = d.readWASignatureData()
		if err != nil {
			return out, fmt.Errorf("unable to read WA signature: %w", err)
		}
	} else {
		return out, fmt.Errorf("unsupported curve ID %d (%s)", uint8(curveID), curveID)
	}

	out, err = ecc.NewSignatureFromData(data)
	if err != nil {
		return out, fmt.Errorf("new signature: %w", err)
	}

	if tracer.Enabled() {
		zlog.Debug("read signature", zap.Stringer("sig", out))
	}

	return
}

func (d *Decoder) readWASignatureData() (out []byte, err error) {
	begin := d.pos
	if d.remaining() < 66 {
		err = fmt.Errorf("signature WA key material requires at least [66] bytes, remaining [%d]", d.remaining())
		return
	}

	// Skip key recover param id (1 byte), R value (32 bytes) and S value (32 bytes)
	d.pos += 65
	authenticatorDataSize, err := d.ReadUvarint32()
	if err != nil {
		return out, fmt.Errorf("unable to read signature WA authenticator data size: %w", err)
	}

	if d.remaining() < int(authenticatorDataSize) {
		err = fmt.Errorf("signature WA authenticator data requires [%d] bytes, remaining [%d]", authenticatorDataSize, d.remaining())
		return
	}
	d.pos += int(authenticatorDataSize)

	clientDataJSONSize, err := d.ReadUvarint32()
	if err != nil {
		return out, fmt.Errorf("unable to read signature WA client data JSON size: %w", err)
	}

	if d.remaining() < int(clientDataJSONSize) {
		err = fmt.Errorf("signature WA client data JSON requires [%d] bytes, remaining [%d]", clientDataJSONSize, d.remaining())
		return
	}
	d.pos += int(clientDataJSONSize)

	signatureMaterialSize := d.pos - begin

	out = make([]byte, signatureMaterialSize+1)
	out[0] = byte(ecc.CurveWA)
	copy(out[1:], d.data[begin:begin+signatureMaterialSize])
	if tracer.Enabled() {
		zlog.Debug("read wa signature data", zap.Stringer("data", HexBytes(out)))
	}

	return
}

func (d *Decoder) ReadTstamp() (out Tstamp, err error) {
	if d.remaining() < TypeSize.Tstamp {
		err = fmt.Errorf("tstamp required [%d] bytes, remaining [%d]", TypeSize.Tstamp, d.remaining())
		return
	}

	unixNano, err := d.ReadUint64()
	out.Time = time.Unix(0, int64(unixNano))
	if tracer.Enabled() {
		zlog.Debug("read tstamp", zap.Time("time", out.Time))
	}
	return
}

func (d *Decoder) ReadBlockTimestamp() (out BlockTimestamp, err error) {
	if d.remaining() < TypeSize.BlockTimestamp {
		err = fmt.Errorf("blockTimestamp required [%d] bytes, remaining [%d]", TypeSize.BlockTimestamp, d.remaining())
		return
	}

	// Encoded value of block timestamp is the slot, which represents the amount of 500 ms that
	// has elapsed since block epoch which is Januaray 1st, 2000 (946684800000 Unix Timestamp Milliseconds)
	n, err := d.ReadUint32()
	milliseconds := int64(n)*500 + 946684800000

	out.Time = time.Unix(0, milliseconds*1000*1000)
	if tracer.Enabled() {
		zlog.Debug("read block timestamp", zap.Time("time", out.Time))
	}
	return
}

func (d *Decoder) ReadTimePoint() (out TimePoint, err error) {
	n, err := d.ReadUint64()
	out = TimePoint(n)
	if tracer.Enabled() {
		zlog.Debug("read TimePoint", zap.Uint64("us", uint64(out)))
	}
	return

}
func (d *Decoder) ReadTimePointSec() (out TimePointSec, err error) {
	n, err := d.ReadUint32()
	out = TimePointSec(n)
	if tracer.Enabled() {
		zlog.Debug("read TimePointSec", zap.Uint32("secs", uint32(out)))
	}
	return

}

func (d *Decoder) ReadJSONTime() (jsonTime JSONTime, err error) {
	n, err := d.ReadUint32()
	jsonTime = JSONTime{time.Unix(int64(n), 0).UTC()}
	if tracer.Enabled() {
		zlog.Debug("read json time", zap.Time("time", jsonTime.Time))
	}
	return
}

func (d *Decoder) ReadName() (out Name, err error) {
	n, err := d.ReadUint64()
	out = Name(NameToString(n))
	if tracer.Enabled() {
		zlog.Debug("read name", zap.String("name", string(out)))
	}
	return
}

func (d *Decoder) ReadCurrencyName() (out CurrencyName, err error) {
	data := d.data[d.pos : d.pos+TypeSize.CurrencyName]
	d.pos += TypeSize.CurrencyName
	out = CurrencyName(strings.TrimRight(string(data), "\x00"))
	if tracer.Enabled() {
		zlog.Debug("read currency name", zap.String("name", string(out)))
	}
	return
}

func (d *Decoder) ReadAsset() (out Asset, err error) {

	amount, err := d.ReadInt64()
	precision, err := d.ReadByte()
	if err != nil {
		return out, fmt.Errorf("readSymbol precision, %s", err)
	}

	if d.remaining() < 7 {
		err = fmt.Errorf("asset symbol required [%d] bytes, remaining [%d]", 7, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+7]
	d.pos += 7

	out = Asset{}
	out.Amount = Int64(amount)
	out.Precision = precision
	out.Symbol.Symbol = strings.TrimRight(string(data), "\x00")
	if tracer.Enabled() {
		zlog.Debug("read asset", zap.Stringer("value", out))
	}
	return
}

func (d *Decoder) ReadExtendedAsset() (out ExtendedAsset, err error) {
	asset, err := d.ReadAsset()
	if err != nil {
		return out, fmt.Errorf("read extended asset: read asset: %w", err)
	}

	contract, err := d.ReadName()
	if err != nil {
		return out, fmt.Errorf("read extended asset: read name: %w", err)
	}

	extendedAsset := ExtendedAsset{
		Asset:    asset,
		Contract: AccountName(contract),
	}

	if tracer.Enabled() {
		zlog.Debug("read extended asset")
	}

	return extendedAsset, err
}

func (d *Decoder) ReadSymbol() (out *Symbol, err error) {
	rawValue, err := d.ReadUint64()
	if err != nil {
		return out, fmt.Errorf("read symbol: %w", err)
	}

	precision := uint8(rawValue & 0xFF)
	symbolCode := SymbolCode(rawValue >> 8).String()

	out = &Symbol{
		Precision: precision,
		Symbol:    symbolCode,
	}

	if tracer.Enabled() {
		zlog.Debug("read symbol", zap.String("symbol", symbolCode), zap.Uint8("precision", precision))
	}
	return
}

func (d *Decoder) ReadSymbolCode() (out SymbolCode, err error) {
	n, err := d.ReadUint64()
	out = SymbolCode(n)

	if tracer.Enabled() {
		zlog.Debug("read symbol code")
	}
	return
}

func (d *Decoder) ReadActionData(action *Action) (err error) {
	actionMap := RegisteredActions[action.Account]

	var decodeInto reflect.Type
	if actionMap != nil {

		objType := actionMap[action.Name]
		if objType != nil {
			if tracer.Enabled() {
				zlog.Debug("read object", zap.String("type", objType.Name()))
			}
			decodeInto = objType
		}
	}
	if decodeInto == nil {
		return
	}

	if tracer.Enabled() {
		zlog.Debug("reflect type", zap.String("type", decodeInto.Name()))
	}
	obj := reflect.New(decodeInto)
	iface := obj.Interface()
	if tracer.Enabled() {
		zlog.Debug("reflect object", typeField("type", iface), zap.Reflect("obj", obj))
	}
	err = UnmarshalBinary(action.ActionData.HexData, iface)
	if err != nil {
		return fmt.Errorf("decoding Action [%s], %s", obj.Type().Name(), err)
	}

	action.ActionData.Data = iface

	return
}

func (d *Decoder) ReadP2PMessageEnvelope() (out *Packet, err error) {

	out = &Packet{}
	l, err := d.ReadUint32()
	if err != nil {
		err = fmt.Errorf("p2p envelope length: %w", err)
		return
	}
	out.Length = l
	b, err := d.ReadByte()
	if err != nil {
		err = fmt.Errorf("p2p envelope type: %w", err)
		return
	}
	out.Type = P2PMessageType(b)

	payloadLength := int(l - 1)
	if d.remaining() < payloadLength {
		err = fmt.Errorf("p2p envelope payload required [%d] bytes, remaining [%d]", l, d.remaining())
		return
	}

	out.Payload = make([]byte, int(payloadLength))
	copy(out.Payload, d.data[d.pos:d.pos+int(payloadLength)])

	d.pos += int(out.Length)
	return
}

func (d *Decoder) remaining() int {
	return len(d.data) - d.pos
}

func (d *Decoder) hasRemaining() bool {
	return d.remaining() > 0
}

func UnmarshalBinaryReader(reader io.Reader, v interface{}) (err error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	return UnmarshalBinary(data, v)
}

func UnmarshalBinary(data []byte, v interface{}) (err error) {
	decoder := NewDecoder(data)
	return decoder.Decode(v)
}
