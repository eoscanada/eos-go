package eos

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"time"

	"errors"
	"reflect"

	"encoding/hex"

	"strings"

	"io/ioutil"

	"github.com/eoscanada/eos-go/ecc"
)

var TypeSize = struct {
	Byte           int
	Int8           int
	UInt8          int
	UInt16         int
	Int16          int
	UInt32         int
	UInt64         int
	Float32        int
	Float64        int
	SHA256Bytes    int
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
	Float32:        4,
	Float64:        8,
	SHA256Bytes:    32,
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

	Logger.Decoder.Print(fmt.Sprintf("Decode type [%T]", v))
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
		Logger.Decoder.Print(fmt.Sprintf("readName [%s]", name))
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
	case *int64:
		var n int64
		n, err = d.ReadInt64()
		rv.SetInt(int64(n))
		return
	case *JSONInt64:
		var n int64
		n, err = d.ReadInt64()
		rv.SetInt(int64(n))
		return
	case *JSONFloat64:
		var n float64
		n, err = d.readFloat64()
		rv.SetFloat(n)
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
	case *SHA256Bytes:
		var s SHA256Bytes
		s, err = d.ReadSHA256Bytes()
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

		Logger.Decoder.Print(fmt.Sprintf("Type byte value : %d", t))

		if t == 0 {
			id, e := d.ReadSHA256Bytes()
			if err != nil {
				err = fmt.Errorf("decode: TransactionWithID failed to read id: %s", e)
				return
			}

			trx := TransactionWithID{ID: id}
			rv.Set(reflect.ValueOf(trx))
			return nil

		} else {
			packedTrx := &PackedTransaction{}
			d.Decode(packedTrx)
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
			Logger.Decoder.Print("Skipping optional OptionalProducerSchedule")
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
		Logger.Decoder.Print("Array")
		len := t.Len()
		for i := 0; i < int(len); i++ {
			if err = d.Decode(rv.Index(i).Addr().Interface()); err != nil {
				return
			}
		}
		return

	case reflect.Slice:
		Logger.Decoder.Print("Reading Slice length ")
		var l uint64
		if l, err = d.ReadUvarint64(); err != nil {
			return
		}
		Logger.Decoder.Print(fmt.Sprintf("Slice [%T] of length: %d", v, l))
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

	case reflect.Map:
		//fmt.Logger.Decoder.Print("Map")
		var l uint64
		if l, err = d.ReadUvarint64(); err != nil {
			return
		}
		kt := t.Key()
		vt := t.Elem()
		rv.Set(reflect.MakeMap(t))
		for i := 0; i < int(l); i++ {
			kv := reflect.Indirect(reflect.New(kt))
			if err = d.Decode(kv.Addr().Interface()); err != nil {
				return
			}
			vv := reflect.Indirect(reflect.New(vt))
			if err = d.Decode(vv.Addr().Interface()); err != nil {
				return
			}
			rv.SetMapIndex(kv, vv)
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

		if v := rv.Field(i); v.CanSet() && t.Field(i).Name != "_" {
			iface := v.Addr().Interface()
			Logger.Decoder.Print(fmt.Sprintf("Field name: %s", t.Field(i).Name))
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
		Logger.Decoder.Print(fmt.Sprintf("readUvarint [%d]", l))
		return l, ErrVarIntBufferSize
	}

	d.pos += read
	Logger.Decoder.Print(fmt.Sprintf("readUvarint [%d]", l))
	return l, nil
}
func (d *Decoder) ReadVarint64() (out int64, err error) {

	l, read := binary.Varint(d.data[d.pos:])
	if read <= 0 {
		Logger.Decoder.Print(fmt.Sprintf("readVarint [%d]", l))
		return l, ErrVarIntBufferSize
	}

	d.pos += read
	Logger.Decoder.Print(fmt.Sprintf("readVarint [%d]", l))
	return l, nil
}

func (d *Decoder) ReadVarint32() (out int32, err error) {

	n, err := d.ReadVarint64()
	if err != nil {
		return out, err
	}
	out = int32(n)
	Logger.Decoder.Print(fmt.Sprintf("readVarint32 [%d]", out))
	return
}
func (d *Decoder) ReadUvarint32() (out uint32, err error) {

	n, err := d.ReadUvarint64()
	if err != nil {
		return out, err
	}
	out = uint32(n)
	Logger.Decoder.Print(fmt.Sprintf("readUvarint32 [%d]", out))
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

	Logger.Decoder.Print(fmt.Sprintf("readByteArray [%s]", hex.EncodeToString(out)))
	return
}

func (d *Decoder) ReadByte() (out byte, err error) {

	if d.remaining() < TypeSize.Byte {
		err = fmt.Errorf("byte required [1] byte, remaining [%d]", d.remaining())
		return
	}

	out = d.data[d.pos]
	d.pos++
	Logger.Decoder.Print(fmt.Sprintf("readByte [%d]", out))
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
	return

}

func (d *Decoder) ReadUInt8() (out uint8, err error) {
	out, err = d.ReadByte()
	Logger.Decoder.Print(fmt.Sprintf("readUint8 [%d]", out))
	return
}
func (d *Decoder) ReadInt8() (out int8, err error) {
	b, err := d.ReadByte()
	out = int8(b)
	Logger.Decoder.Print(fmt.Sprintf("readInt8 [%d]", out))
	return
}

func (d *Decoder) ReadUint16() (out uint16, err error) {
	if d.remaining() < TypeSize.UInt16 {
		err = fmt.Errorf("uint16 required [%d] bytes, remaining [%d]", TypeSize.UInt16, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint16(d.data[d.pos:])
	d.pos += TypeSize.UInt16
	Logger.Decoder.Print(fmt.Sprintf("readUint16 [%d]", out))
	return
}

func (d *Decoder) ReadInt16() (out int16, err error) {
	n, err := d.ReadUint16()
	out = int16(n)
	return
}
func (d *Decoder) ReadInt64() (out int64, err error) {
	n, err := d.ReadUint64()
	out = int64(n)
	return
}

func (d *Decoder) ReadUint32() (out uint32, err error) {
	if d.remaining() < TypeSize.UInt32 {
		err = fmt.Errorf("uint32 required [%d] bytes, remaining [%d]", TypeSize.UInt32, d.remaining())
		return
	}

	out = binary.LittleEndian.Uint32(d.data[d.pos:])
	d.pos += TypeSize.UInt32
	Logger.Decoder.Print(fmt.Sprintf("readUint32 [%d]", out))
	return
}
func (d *Decoder) ReadInt32() (out int32, err error) {
	n, err := d.ReadUint32()
	out = int32(n)
	Logger.Decoder.Print(fmt.Sprintf("readInt32 [%d]", out))
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
	Logger.Decoder.Print(fmt.Sprintf("readUint64 [%d] [%s]", out, hex.EncodeToString(data)))
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
	Logger.Decoder.Print(fmt.Sprintf("readFloat32 [%f]", out))
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
	Logger.Decoder.Print(fmt.Sprintf("readFloat64 [%f]", out))
	return
}

func (d *Decoder) readFloat64() (out float64, err error) {
	if d.remaining() < TypeSize.Float64 {
		err = fmt.Errorf("uint64 required [%d] bytes, remaining [%d]", TypeSize.UInt64, d.remaining())
		return
	}

	data := d.data[d.pos : d.pos+TypeSize.Float64]
	out = math.Float64frombits(binary.LittleEndian.Uint64(data))
	d.pos += TypeSize.Float64
	println(fmt.Sprintf("readFloat64 [%f] [%s]", out, hex.EncodeToString(data)))
	return
}

func (d *Decoder) ReadString() (out string, err error) {
	data, err := d.ReadByteArray()
	out = string(data)
	Logger.Decoder.Print(fmt.Sprintf("readString [%s]", out))
	return
}

func (d *Decoder) ReadSHA256Bytes() (out SHA256Bytes, err error) {

	if d.remaining() < TypeSize.SHA256Bytes {
		err = fmt.Errorf("sha256 required [%d] bytes, remaining [%d]", TypeSize.SHA256Bytes, d.remaining())
		return
	}

	out = SHA256Bytes(d.data[d.pos : d.pos+TypeSize.SHA256Bytes])
	d.pos += TypeSize.SHA256Bytes
	Logger.Decoder.Print(fmt.Sprintf("readSHA256Bytes [%s]", hex.EncodeToString(out)))
	return
}

func (d *Decoder) ReadChecksum160() (out Checksum160, err error) {

	if d.remaining() < TypeSize.Checksum160 {
		err = fmt.Errorf("checksum 160 required [%d] bytes, remaining [%d]", TypeSize.Checksum160, d.remaining())
		return
	}

	out = d.data[d.pos : d.pos+TypeSize.Checksum160]
	d.pos += TypeSize.Checksum160
	Logger.Decoder.Print(fmt.Sprintf("ReadChecksum160Bytes [%s]", hex.EncodeToString(out)))
	return
}

func (d *Decoder) ReadChecksum256() (out Checksum256, err error) {

	if d.remaining() < TypeSize.Checksum256 {
		err = fmt.Errorf("checksum 256 required [%d] bytes, remaining [%d]", TypeSize.Checksum256, d.remaining())
		return
	}

	out = d.data[d.pos : d.pos+TypeSize.Checksum256]
	d.pos += TypeSize.Checksum256
	Logger.Decoder.Print(fmt.Sprintf("ReadChecksum256Bytes [%s]", hex.EncodeToString(out)))
	return
}

func (d *Decoder) ReadChecksum512() (out Checksum512, err error) {

	if d.remaining() < TypeSize.Checksum512 {
		err = fmt.Errorf("checksum 512 required [%d] bytes, remaining [%d]", TypeSize.Checksum512, d.remaining())
		return
	}

	out = d.data[d.pos : d.pos+TypeSize.Checksum512]
	d.pos += TypeSize.Checksum512
	Logger.Decoder.Print(fmt.Sprintf("ReadChecksum512Bytes [%s]", hex.EncodeToString(out)))
	return
}

func (d *Decoder) ReadPublicKey() (out ecc.PublicKey, err error) {

	if d.remaining() < TypeSize.PublicKey {
		err = fmt.Errorf("publicKey required [%d] bytes, remaining [%d]", TypeSize.PublicKey, d.remaining())
		return
	}
	out = ecc.PublicKey{
		Curve:   ecc.CurveID(d.data[d.pos]),                 // 1 byte
		Content: d.data[d.pos+1 : d.pos+TypeSize.PublicKey], // 33 bytes
	}
	d.pos += TypeSize.PublicKey
	Logger.Decoder.Print(fmt.Sprintf("readPublicKey [curve=%d, content=%s]", out.Curve, hex.EncodeToString(out.Content)))
	return
}

func (d *Decoder) ReadSignature() (out ecc.Signature, err error) {
	if d.remaining() < TypeSize.Signature {
		err = fmt.Errorf("signature required [%d] bytes, remaining [%d]", TypeSize.Signature, d.remaining())
		return
	}
	out = ecc.Signature{
		Curve:   ecc.CurveID(d.data[d.pos]),                 // 1 byte
		Content: d.data[d.pos+1 : d.pos+TypeSize.Signature], // 65 bytes
	}
	d.pos += TypeSize.Signature
	Logger.Decoder.Print(fmt.Sprintf("readSignature [curve=%d, content=%s]", out.Curve, hex.EncodeToString(out.Content)))
	return
}

func (d *Decoder) ReadTstamp() (out Tstamp, err error) {

	if d.remaining() < TypeSize.Tstamp {
		err = fmt.Errorf("tstamp required [%d] bytes, remaining [%d]", TypeSize.Tstamp, d.remaining())
		return
	}

	unixNano, err := d.ReadUint64()
	out.Time = time.Unix(0, int64(unixNano))
	Logger.Decoder.Print(fmt.Sprintf("readTstamp [%s]", out))
	return
}

func (d *Decoder) ReadBlockTimestamp() (out BlockTimestamp, err error) {
	if d.remaining() < TypeSize.BlockTimestamp {
		err = fmt.Errorf("blockTimestamp required [%d] bytes, remaining [%d]", TypeSize.BlockTimestamp, d.remaining())
		return
	}
	n, err := d.ReadUint32()
	out.Time = time.Unix(int64(n)+946684800, 0)
	return
}

func (d *Decoder) ReadTimePoint() (out TimePoint, err error) {
	n, err := d.ReadUint64()
	out = TimePoint(n)
	Logger.Decoder.Print(fmt.Sprintf("ReadTimePointSec [%d]", out))
	return

}
func (d *Decoder) ReadTimePointSec() (out TimePointSec, err error) {
	n, err := d.ReadUint32()
	out = TimePointSec(n)
	Logger.Decoder.Print(fmt.Sprintf("ReadTimePointSec [%d]", out))
	return

}

func (d *Decoder) ReadJSONTime() (jsonTime JSONTime, err error) {
	n, err := d.ReadUint32()
	jsonTime = JSONTime{time.Unix(int64(n), 0).UTC()}
	Logger.Decoder.Print("readJSONTime: ", jsonTime)
	return
}

func (d *Decoder) ReadName() (out Name, err error) {

	n, err := d.ReadUint64()
	out = Name(NameToString(n))
	Logger.Decoder.Print(fmt.Sprintf("readName [%s]", out))
	return
}

func (d *Decoder) ReadCurrencyName() (out CurrencyName, err error) {

	data := d.data[d.pos : d.pos+TypeSize.CurrencyName]
	d.pos += TypeSize.CurrencyName

	out = CurrencyName(strings.TrimRight(string(data), "\x00"))
	return
}

func (d *Decoder) ReadAsset() (out Asset, err error) {

	amount, err := d.ReadInt64()
	precision, err := d.ReadByte()
	if err != nil {
		return out, fmt.Errorf("readSymbol precision, %s", err)
	}

	data := d.data[d.pos : d.pos+7]
	d.pos += 7

	out = Asset{}
	out.Amount = amount
	out.Precision = precision
	out.Symbol.Symbol = strings.TrimRight(string(data), "\x00")
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
	return
}

func (d *Decoder) ReadSymbolCode() (out SymbolCode, err error) {

	n, err := d.ReadUint64()
	out = SymbolCode(n)
	Logger.Decoder.Print(fmt.Sprintf("ReadSymbolCode [%d]", out))

	return
}

func (d *Decoder) ReadActionData(action *Action) (err error) {

	actionMap := RegisteredActions[action.Account]

	var decodeInto reflect.Type
	if actionMap != nil {
		objType := actionMap[action.Name]
		Logger.Decoder.Print("object type :", objType)
		if objType != nil {
			decodeInto = objType
		}
	}
	if decodeInto == nil {
		return
	}

	Logger.Decoder.Print("Reflect type :", decodeInto)
	obj := reflect.New(decodeInto)
	Logger.Decoder.Print("obj :", obj)
	err = UnmarshalBinary(action.ActionData.HexData, obj.Interface())
	if err != nil {
		return fmt.Errorf("decoding Action [%s], %s", obj.Type().Name(), err)
	}

	Logger.Decoder.Print("Object type :", obj.Interface())
	action.ActionData.Data = obj.Interface()

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
	payload := d.data[d.pos : d.pos+int(payloadLength)]
	d.pos += int(out.Length)

	out.Payload = payload
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
