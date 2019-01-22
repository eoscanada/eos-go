package eos

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"time"

	"errors"
	"reflect"

	"strings"

	"io/ioutil"

	"github.com/eoscanada/eos-go/ecc"
	"go.uber.org/zap"
)

var TypeSize = struct {
	Byte           int
	Int8           int
	UInt8          int
	UInt16         int
	Int16          int
	UInt32         int
	UInt64         int
	UInt128        int
	Float32        int
	Float64        int
	Checksum160    int
	Checksum256    int
	Checksum512    int
	PublicKey      int
	Signature      int
	Tstamp         int
	BlockTimestamp int
	CurrencyName   int
	Bool           int
}{
	Byte:           1,
	Int8:           1,
	UInt8:          1,
	UInt16:         2,
	Int16:          2,
	UInt32:         4,
	UInt64:         8,
	UInt128:        16,
	Float32:        4,
	Float64:        8,
	Checksum160:    20,
	Checksum256:    32,
	Checksum512:    64,
	PublicKey:      34,
	Signature:      66,
	Tstamp:         8,
	BlockTimestamp: 4,
	CurrencyName:   7,
	Bool:           1,
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
	data               []byte
	pos                int
	decodeP2PMessage   bool
	decodeTransactions bool
	decodeActions      bool
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		data:               data,
		decodeP2PMessage:   true,
		decodeTransactions: true,
		decodeActions:      true,
	}
}

func (d *Decoder) DecodeP2PMessage(decode bool) {
	d.decodeP2PMessage = decode
}

func (d *Decoder) DecodeActions(decode bool) {
	d.decodeActions = decode
}

func (d *Decoder) Decode(v interface{}) (err error) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if !rv.CanAddr() {
		return errors.New("decode, can only Decode to pointer type")
	}
	t := rv.Type()

	decoderLog.Debug("decode type", typeField("type", v))
	if !rv.CanAddr() {
		return errors.New("binary: can only Decode to pointer type")
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		newRV := reflect.New(t)
		rv.Set(newRV)
		rv = reflect.Indirect(newRV)
	}

	switch realV := v.(type) {
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
		decoderLog.Debug("read name", zap.String("name", name))
		rv.SetString(name)
		return
	case *byte, *P2PMessageType, *TransactionStatus, *CompressionType, *IDListMode, *GoAwayReason:
		var n byte
		n, err = d.ReadByte()
		rv.SetUint(uint64(n))
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
	case *JSONFloat64:
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
	case *Asset:
		var asset Asset
		asset, err = d.ReadAsset()
		rv.Set(reflect.ValueOf(asset))
		return

	case *TransactionWithID:

		t, e := d.ReadByte()
		if err != nil {
			err = fmt.Errorf("decode: TransactionWithID failed to read type byte: %s", e)
			return
		}

		decoderLog.Debug("type byte value", zap.Uint8("val", t))

		if t == 0 {
			id, e := d.ReadChecksum256()
			if err != nil {
				err = fmt.Errorf("decode: TransactionWithID failed to read id: %s", e)
				return
			}

			trx := TransactionWithID{ID: id}
			rv.Set(reflect.ValueOf(trx))
			return nil

		} else {
			packedTrx := &PackedTransaction{}
			if err := d.Decode(packedTrx); err != nil {
				return err
			}
			trx := TransactionWithID{Packed: packedTrx}
			rv.Set(reflect.ValueOf(trx))
			return nil
		}

	case **OptionalProducerSchedule:
		isPresent, e := d.ReadByte()
		if e != nil {
			err = fmt.Errorf("decode: OptionalProducerSchedule isPresent, %s", e)
			return
		}

		if isPresent == 0 {
			decoderLog.Debug("skipping optional OptionalProducerSchedule")
			*realV = nil
			return
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
		decoderLog.Debug("reading array")
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
		decoderLog.Debug("reading slice", zap.Uint64("len", l), typeField("type", v))
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

	for i := 0; i < l; i++ {

		if tag := t.Field(i).Tag.Get("eos"); tag == "-" {
			continue
		}

		typeField := t.Field(i)
		if v := rv.Field(i); v.CanSet() && typeField.Name != "_" {
			iface := v.Addr().Interface()
			decoderLog.Debug("field", zap.String("name", typeField.Name))
			if err = d.Decode(iface); err != nil {
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
	decoderLog.Debug("read uvarint64", zap.Uint64("val", l))
	d.pos += read
	return l, nil
}
func (d *Decoder) ReadVarint64() (out int64, err error) {
	l, read := binary.Varint(d.data[d.pos:])
	if read <= 0 {
		return l, ErrVarIntBufferSize
	}
	decoderLog.Debug("read varint", zap.Int64("val", l))
	d.pos += read
	return l, nil
}

func (d *Decoder) ReadVarint32() (out int32, err error) {
	n, err := d.ReadVarint64()
	if err != nil {
		return out, err
	}
	out = int32(n)
	decoderLog.Debug("read varint32", zap.Int32("val", out))
	return
}
func (d *Decoder) ReadUvarint32() (out uint32, err error) {

	n, err := d.ReadUvarint64()
	if err != nil {
		return out, err
	}
	out = uint32(n)
	decoderLog.Debug("read uvarint32", zap.Uint32("val", out))
	return
}

func (d *Decoder) ReadByteArray() (out []byte, err error) {

	l, err := d.ReadUvarint64()
	if err != nil {
		return nil, err
	}

	if len(d.data) < d.pos+int(l) {
		return nil, fmt.Errorf("byte array: varlen=%d, missing %d bytes", l, d.pos+int(l)-len(d.data))
	}

	out = d.data[d.pos : d.pos+int(l)]
	d.pos += int(l)
	decoderLog.Debug("read byte array", zap.Stringer("hex", HexBytes(out)))
	return
}

func (d *Decoder) ReadByte() (out byte, err error) {
	if d.remaining() < TypeSize.Byte {
		err = fmt.Errorf("byte required [1] byte, remaining [%d]", d.remaining())
		return
	}

	out = d.data[d.pos]
	d.pos++
	decoderLog.Debug("read byte", zap.Uint8("byte", out))
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
	decoderLog.Debug("read bool", zap.Bool("val", out))
	return

}

func (d *Decoder) ReadUInt8() (out uint8, err error) {
	out, err = d.ReadByte()
	return
}
func (d *Decoder) ReadInt8() (out int8, err error) {
	b, err := d.ReadByte()
	out = int8(b)
	decoderLog.Debug("read int8", zap.Int8("val", out))
	return
}

func (d *Decoder) ReadUint16() (out uint16, err error) {
	if d.remaining() < TypeSize.UInt16 {
		err = fmt.Errorf("uint16 required [%d] bytes, remaining [%d]", TypeSize.UInt16, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint16(d.data[d.pos:])
	d.pos += TypeSize.UInt16
	decoderLog.Debug("read uint16", zap.Uint16("val", out))
	return
}

func (d *Decoder) ReadInt16() (out int16, err error) {
	n, err := d.ReadUint16()
	out = int16(n)
	decoderLog.Debug("read int16", zap.Int16("val", out))
	return
}
func (d *Decoder) ReadInt64() (out int64, err error) {
	n, err := d.ReadUint64()
	out = int64(n)
	decoderLog.Debug("read int64", zap.Int64("val", out))
	return
}

func (d *Decoder) ReadUint32() (out uint32, err error) {
	if d.remaining() < TypeSize.UInt32 {
		err = fmt.Errorf("uint32 required [%d] bytes, remaining [%d]", TypeSize.UInt32, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint32(d.data[d.pos:])
	d.pos += TypeSize.UInt32
	decoderLog.Debug("read uint32", zap.Uint32("val", out))
	return
}
func (d *Decoder) ReadInt32() (out int32, err error) {
	n, err := d.ReadUint32()
	out = int32(n)
	decoderLog.Debug("read int32", zap.Int32("val", out))
	return
}

func (d *Decoder) ReadUint64() (out uint64, err error) {
	if d.remaining() < TypeSize.UInt64 {
		err = fmt.Errorf("uint64 required [%d] bytes, remaining [%d]", TypeSize.UInt64, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.UInt64]
	out = binary.LittleEndian.Uint64(data)
	d.pos += TypeSize.UInt64
	decoderLog.Debug("read uint64", zap.Uint64("val", out), zap.Stringer("hex", HexBytes(data)))
	return
}

func (d *Decoder) ReadUint128(typeName string) (out Uint128, err error) {
	if d.remaining() < TypeSize.UInt128 {
		err = fmt.Errorf("%s required [%d] bytes, remaining [%d]", typeName, TypeSize.UInt128, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.UInt128]
	out.Lo = binary.LittleEndian.Uint64(data)
	out.Hi = binary.LittleEndian.Uint64(data[8:])

	d.pos += TypeSize.UInt128
	decoderLog.Debug("read uint128", zap.Stringer("hex", out), zap.Uint64("lo", out.Lo), zap.Uint64("lo", out.Lo))
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
	decoderLog.Debug("read float32", zap.Float32("val", out))
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
	decoderLog.Debug("read Float64", zap.Float64("val", float64(out)))
	return
}

func (d *Decoder) ReadString() (out string, err error) {
	data, err := d.ReadByteArray()
	out = string(data)
	decoderLog.Debug("read string", zap.String("val", out))
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
	decoderLog.Debug("read checksum160", zap.Stringer("hex", HexBytes(out)))
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
	decoderLog.Debug("read checksum256", zap.Stringer("hex", HexBytes(out)))
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
	decoderLog.Debug("read checksum512", zap.Stringer("hex", HexBytes(out)))
	return
}

func (d *Decoder) ReadPublicKey() (out ecc.PublicKey, err error) {

	if d.remaining() < TypeSize.PublicKey {
		err = fmt.Errorf("publicKey required [%d] bytes, remaining [%d]", TypeSize.PublicKey, d.remaining())
		return
	}
	keyContent := make([]byte, 34)
	copy(keyContent, d.data[d.pos:d.pos+TypeSize.PublicKey])

	out, err = ecc.NewPublicKeyFromData(keyContent)
	if err != nil {
		err = fmt.Errorf("publicKey: key from data: %s", err)
	}

	d.pos += TypeSize.PublicKey
	decoderLog.Debug("read public key", zap.Stringer("pubkey", out))
	return
}

func (d *Decoder) ReadSignature() (out ecc.Signature, err error) {
	if d.remaining() < TypeSize.Signature {
		err = fmt.Errorf("signature required [%d] bytes, remaining [%d]", TypeSize.Signature, d.remaining())
		return
	}

	sigContent := make([]byte, 66)
	copy(sigContent, d.data[d.pos:d.pos+TypeSize.Signature])

	out, err = ecc.NewSignatureFromData(sigContent)
	if err != nil {
		return out, fmt.Errorf("new signature: %s", err)
	}

	d.pos += TypeSize.Signature
	decoderLog.Debug("read signature", zap.Stringer("sig", out))
	return
}

func (d *Decoder) ReadTstamp() (out Tstamp, err error) {
	if d.remaining() < TypeSize.Tstamp {
		err = fmt.Errorf("tstamp required [%d] bytes, remaining [%d]", TypeSize.Tstamp, d.remaining())
		return
	}

	unixNano, err := d.ReadUint64()
	out.Time = time.Unix(0, int64(unixNano))
	decoderLog.Debug("read tstamp", zap.Time("time", out.Time))
	return
}

func (d *Decoder) ReadBlockTimestamp() (out BlockTimestamp, err error) {
	if d.remaining() < TypeSize.BlockTimestamp {
		err = fmt.Errorf("blockTimestamp required [%d] bytes, remaining [%d]", TypeSize.BlockTimestamp, d.remaining())
		return
	}
	n, err := d.ReadUint32()
	out.Time = time.Unix(int64(n)+946684800, 0)
	decoderLog.Debug("read block timestamp", zap.Time("time", out.Time))
	return
}

func (d *Decoder) ReadTimePoint() (out TimePoint, err error) {
	n, err := d.ReadUint64()
	out = TimePoint(n)
	decoderLog.Debug("read TimePoint", zap.Uint64("us", uint64(out)))
	return

}
func (d *Decoder) ReadTimePointSec() (out TimePointSec, err error) {
	n, err := d.ReadUint32()
	out = TimePointSec(n)
	decoderLog.Debug("read TimePointSec", zap.Uint32("secs", uint32(out)))
	return

}

func (d *Decoder) ReadJSONTime() (jsonTime JSONTime, err error) {
	n, err := d.ReadUint32()
	jsonTime = JSONTime{time.Unix(int64(n), 0).UTC()}
	decoderLog.Debug("read json time", zap.Time("time", jsonTime.Time))
	return
}

func (d *Decoder) ReadName() (out Name, err error) {
	n, err := d.ReadUint64()
	out = Name(NameToString(n))
	decoderLog.Debug("read name", zap.String("name", string(out)))
	return
}

func (d *Decoder) ReadCurrencyName() (out CurrencyName, err error) {
	data := d.data[d.pos : d.pos+TypeSize.CurrencyName]
	d.pos += TypeSize.CurrencyName
	out = CurrencyName(strings.TrimRight(string(data), "\x00"))
	decoderLog.Debug("read currency name", zap.String("name", string(out)))
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
	decoderLog.Debug("read asset", zap.Stringer("value", out))
	return
}

func (d *Decoder) ReadExtendedAsset() (out ExtendedAsset, err error) {
	asset, err := d.ReadAsset()
	if err != nil {
		return out, fmt.Errorf("read extended asset: read asset: %s", err)
	}

	contract, err := d.ReadName()
	if err != nil {
		return out, fmt.Errorf("read extended asset: read name: %s", err)
	}

	extendedAsset := ExtendedAsset{
		Asset:    asset,
		Contract: AccountName(contract),
	}

	decoderLog.Debug("read extended asset")

	return extendedAsset, err
}

func (d *Decoder) ReadSymbol() (out *Symbol, err error) {

	precision, err := d.ReadUInt8()
	if err != nil {
		return out, fmt.Errorf("read symbol: read precision: %s", err)
	}
	symbol, err := d.ReadString()
	if err != nil {
		return out, fmt.Errorf("read symbol: read symbol: %s", err)
	}

	out = &Symbol{
		Precision: precision,
		Symbol:    symbol,
	}
	decoderLog.Debug("read symbol", zap.String("symbol", symbol), zap.Uint8("precision", precision))
	return
}

func (d *Decoder) ReadSymbolCode() (out SymbolCode, err error) {

	n, err := d.ReadUint64()
	out = SymbolCode(n)

	decoderLog.Debug("read symbol code")
	return
}

func (d *Decoder) ReadActionData(action *Action) (err error) {
	actionMap := RegisteredActions[action.Account]

	var decodeInto reflect.Type
	if actionMap != nil {

		objType := actionMap[action.Name]
		if objType != nil {
			decoderLog.Debug("read object", zap.String("type", objType.Name()))
			decodeInto = objType
		}
	}
	if decodeInto == nil {
		return
	}

	decoderLog.Debug("reflect type", zap.String("type", decodeInto.Name()))
	obj := reflect.New(decodeInto)
	iface := obj.Interface()
	decoderLog.Debug("reflect object", typeField("type", iface), zap.Reflect("obj", obj))
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
		err = fmt.Errorf("p2p envelope length: %s", err)
		return
	}
	out.Length = l
	b, err := d.ReadByte()
	if err != nil {
		err = fmt.Errorf("p2p envelope type: %s", err)
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
