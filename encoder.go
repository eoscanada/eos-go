package eos

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"

	"math"

	"github.com/eoscanada/eos-go/ecc"
)

// --------------------------------------------------------------
// Encoder implements the EOS packing, similar to FC_BUFFER
// --------------------------------------------------------------
type Encoder struct {
	output io.Writer
	Order  binary.ByteOrder
	count  int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		output: w,
		Order:  binary.LittleEndian,
		count:  0,
	}
}

func (e *Encoder) writeName(name Name) error {
	val, err := StringToName(string(name))
	if err != nil {
		return fmt.Errorf("writeName: %s", err)
	}
	return e.writeUint64(val)
}

func (e *Encoder) Encode(v interface{}) (err error) {
	switch cv := v.(type) {
	case Name:
		return e.writeName(cv)
	case AccountName:
		name := Name(cv)
		return e.writeName(name)
	case PermissionName:
		name := Name(cv)
		return e.writeName(name)
	case ActionName:
		name := Name(cv)
		return e.writeName(name)
	case TableName:
		name := Name(cv)
		return e.writeName(name)
	case ScopeName:
		name := Name(cv)
		return e.writeName(name)
	case string:
		return e.writeString(cv)
	case CompressionType:
		return e.writeByte(uint8(cv))
	case TransactionStatus:
		return e.writeByte(uint8(cv))
	case IDListMode:
		return e.writeByte(byte(cv))
	case byte:
		return e.writeByte(cv)
	case int8:
		return e.writeByte(byte(cv))
	case int16:
		return e.writeInt16(cv)
	case uint16:
		return e.writeUint16(cv)
	case int32:
		return e.writeInt32(cv)
	case uint32:
		return e.writeUint32(cv)
	case uint64:
		return e.writeUint64(cv)
	case Int64:
		return e.writeUint64(uint64(cv))
	case Uint64:
		return e.writeUint64(uint64(cv))
	case int64:
		return e.writeInt64(cv)
	case float32:
		return e.writeFloat32(cv)
	case float64:
		return e.writeFloat64(cv)
	case Varint32:
		return e.writeVarInt(int(cv))
	case Uint128:
		return e.writeUint128(cv)
	case Int128:
		return e.writeUint128(Uint128(cv))
	case Float128:
		return e.writeUint128(Uint128(cv))
	case Varuint32:
		return e.writeUVarInt(int(cv))
	case bool:
		return e.writeBool(cv)
	case Bool:
		return e.writeBool(bool(cv))
	case JSONTime:
		return e.writeJSONTime(cv)
	case HexBytes:
		return e.writeByteArray(cv)
	case Checksum160:
		return e.writeChecksum160(cv)
	case Checksum256:
		return e.writeChecksum256(cv)
	case Checksum512:
		return e.writeChecksum512(cv)
	case []byte:
		return e.writeByteArray(cv)
	case ecc.PublicKey:
		return e.writePublicKey(cv)
	case ecc.Signature:
		return e.writeSignature(cv)
	case Tstamp:
		return e.writeTstamp(cv)
	case BlockTimestamp:
		return e.writeBlockTimestamp(cv)
	case CurrencyName:
		return e.writeCurrencyName(cv)
	case SymbolCode:
		return e.writeUint64(uint64(cv))
	case Asset:
		return e.writeAsset(cv)
	// case *OptionalProducerSchedule:
	// 	isPresent := cv != nil
	// 	e.writeBool(isPresent)
	// 	if isPresent {

	// 	}
	case ActionData:
		Logger.Encoder.Print("ActionData")
		return e.writeActionData(cv)
	case *ActionData:
		Logger.Encoder.Print("*ActionData")
		return e.writeActionData(*cv)
	case *Packet:
		return e.writeBlockP2PMessageEnvelope(*cv)
	case TimePoint:
		return e.writeUint64(uint64(cv))
	case TimePointSec:
		return e.writeUint64(uint64(cv))
	case nil:
	default:

		rv := reflect.Indirect(reflect.ValueOf(v))
		t := rv.Type()

		switch t.Kind() {

		case reflect.Array:
			l := t.Len()
			defer Logger.Encoder.SetPrefix(Logger.Encoder.Prefix())
			Logger.Encoder.SetPrefix(Logger.Encoder.Prefix() + "\t")
			Logger.Encoder.Printf("Encode: array [%T] of length: %d\n", v, l)

			for i := 0; i < l; i++ {
				if err = e.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}
		case reflect.Slice:
			l := rv.Len()
			if err = e.writeUVarInt(l); err != nil {
				return
			}
			defer Logger.Encoder.SetPrefix(Logger.Encoder.Prefix())
			Logger.Encoder.SetPrefix(Logger.Encoder.Prefix() + "\t")

			Logger.Encoder.Printf("Encode: slice [%T] of length: %d\n", v, l)

			for i := 0; i < l; i++ {
				if err = e.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}
		case reflect.Struct:
			l := rv.NumField()
			Logger.Encoder.Printf("Encode: struct [%T] with %d field\n", v, l)
			defer Logger.Encoder.SetPrefix(Logger.Encoder.Prefix())
			Logger.Encoder.SetPrefix(Logger.Encoder.Prefix() + "\t")

			for i := 0; i < l; i++ {
				field := t.Field(i)
				Logger.Encoder.Printf("field -> %s\n", field.Name)

				tag := field.Tag.Get("eos")
				if tag == "-" {
					continue
				}

				if v := rv.Field(i); t.Field(i).Name != "_" {
					if v.CanInterface() {
						isPresent := true
						if tag == "optional" {
							isPresent = !v.IsNil()
							e.writeBool(isPresent)
						}

						if isPresent {
							if err = e.Encode(v.Interface()); err != nil {
								return
							}
						}
					}
				}
			}

		case reflect.Map:
			l := rv.Len()
			if err = e.writeUVarInt(l); err != nil {
				return
			}
			Logger.Encoder.Printf("Map [%T] of length: %d\n", v, l)
			for _, key := range rv.MapKeys() {
				value := rv.MapIndex(key)
				if err = e.Encode(key.Interface()); err != nil {
					return err
				}
				if err = e.Encode(value.Interface()); err != nil {
					return err
				}
			}
		default:
			return errors.New("Encode: unsupported type " + t.String())
		}
	}

	return
}

func (e *Encoder) toWriter(bytes []byte) (err error) {

	e.count += len(bytes)
	Logger.Encoder.Printf("    Appending : [%s] at pos [%d]\n", hex.EncodeToString(bytes), e.count)
	_, err = e.output.Write(bytes)
	return
}

func (e *Encoder) writeByteArray(b []byte) error {
	Logger.Encoder.Printf("Writing byte array of len [%d]\n", len(b))
	if err := e.writeUVarInt(len(b)); err != nil {
		return err
	}
	return e.toWriter(b)
}

func (e *Encoder) writeUVarInt(v int) (err error) {
	Logger.Encoder.Printf("Writing uvarint [%d]\n", v)
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return e.toWriter(buf[:l])
}

func (e *Encoder) writeVarInt(v int) (err error) {
	Logger.Encoder.Printf("Writing varint [%d]\n", v)
	buf := make([]byte, 8)
	l := binary.PutVarint(buf, int64(v))
	return e.toWriter(buf[:l])
}

func (e *Encoder) writeByte(b byte) (err error) {
	Logger.Encoder.Printf("Writing byte [%d]\n", b)
	return e.toWriter([]byte{b})
}

func (e *Encoder) writeBool(b bool) (err error) {
	Logger.Encoder.Printf("Writing byte [%t]\n", b)
	var out byte
	if b {
		out = 1
	}
	return e.writeByte(out)
}

func (e *Encoder) writeUint16(i uint16) (err error) {
	Logger.Encoder.Printf("Writing uint16 [%d]\n", i)
	buf := make([]byte, TypeSize.UInt16)
	binary.LittleEndian.PutUint16(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeInt16(i int16) (err error) {
	Logger.Encoder.Printf("Writing int16 [%d]\n", i)
	return e.writeUint16(uint16(i))
}

func (e *Encoder) writeInt32(i int32) (err error) {
	Logger.Encoder.Printf("Writing int32 [%d]\n", i)
	return e.writeUint32(uint32(i))
}

func (e *Encoder) writeUint32(i uint32) (err error) {
	Logger.Encoder.Printf("Writing unint32 [%d]\n", i)
	buf := make([]byte, TypeSize.UInt32)
	binary.LittleEndian.PutUint32(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeInt64(i int64) (err error) {
	Logger.Encoder.Printf("Writing int64 [%d]\n", i)
	return e.writeUint64(uint64(i))
}

func (e *Encoder) writeUint64(i uint64) (err error) {
	Logger.Encoder.Printf("Writing uint64 [%d]\n", i)
	buf := make([]byte, TypeSize.UInt64)
	binary.LittleEndian.PutUint64(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeUint128(i Uint128) (err error) {
	Logger.Encoder.Printf("Writing uint128 [%d]\n", i)
	buf := make([]byte, TypeSize.UInt128)
	binary.LittleEndian.PutUint64(buf, i.Lo)
	binary.LittleEndian.PutUint64(buf[TypeSize.UInt64:], i.Hi)
	return e.toWriter(buf)
}

func (e *Encoder) writeFloat32(f float32) (err error) {
	Logger.Encoder.Printf("Writing float32 [%f]\n", f)
	i := math.Float32bits(f)
	buf := make([]byte, TypeSize.UInt32)
	binary.LittleEndian.PutUint32(buf, i)

	return e.toWriter(buf)
}
func (e *Encoder) writeFloat64(f float64) (err error) {
	Logger.Encoder.Printf("Writing float64 [%f]\n", f)
	i := math.Float64bits(f)
	buf := make([]byte, TypeSize.UInt64)
	binary.LittleEndian.PutUint64(buf, i)

	return e.toWriter(buf)
}

func (e *Encoder) writeString(s string) (err error) {
	Logger.Encoder.Printf("Writing string [%s]\n", s)
	return e.writeByteArray([]byte(s))
}

func (e *Encoder) writeChecksum160(checksum Checksum160) error {
	Logger.Encoder.Printf("Writing checksum160 [%s]\n", hex.EncodeToString(checksum))
	if len(checksum) == 0 {
		return e.toWriter(bytes.Repeat([]byte{0}, TypeSize.Checksum160))
	}
	return e.toWriter(checksum)
}

func (e *Encoder) writeChecksum256(checksum Checksum256) error {
	Logger.Encoder.Printf("Writing checksum256 [%s]\n", hex.EncodeToString(checksum))
	if len(checksum) == 0 {
		return e.toWriter(bytes.Repeat([]byte{0}, TypeSize.Checksum256))
	}
	return e.toWriter(checksum)
}

func (e *Encoder) writeChecksum512(checksum Checksum512) error {
	Logger.Encoder.Printf("Writing checksum512 [%s]\n", hex.EncodeToString(checksum))
	if len(checksum) == 0 {
		return e.toWriter(bytes.Repeat([]byte{0}, TypeSize.Checksum512))
	}
	return e.toWriter(checksum)
}

func (e *Encoder) writePublicKey(pk ecc.PublicKey) (err error) {
	Logger.Encoder.Printf("Writing public key [%s]\n", pk.String())
	if len(pk.Content) != 33 {
		return fmt.Errorf("public key %q should be 33 bytes, was %d", hex.EncodeToString(pk.Content), len(pk.Content))
	}

	if err = e.writeByte(byte(pk.Curve)); err != nil {
		return err
	}

	return e.toWriter(pk.Content)
}

func (e *Encoder) writeSignature(s ecc.Signature) (err error) {
	Logger.Encoder.Printf("Writing signature [%s]\n", s.String())
	if len(s.Content) != 65 {
		return fmt.Errorf("signature should be 65 bytes, was %d", len(s.Content))
	}

	if err = e.writeByte(byte(s.Curve)); err != nil {
		return
	}

	return e.toWriter(s.Content) // should write 65 bytes
}

func (e *Encoder) writeTstamp(t Tstamp) (err error) {
	Logger.Encoder.Printf("Writing tstamp [%s]\n", t)
	n := uint64(t.UnixNano())
	return e.writeUint64(n)
}

func (e *Encoder) writeBlockTimestamp(bt BlockTimestamp) (err error) {
	Logger.Encoder.Printf("Writing block time stamp [%s]\n", bt)
	n := uint32(bt.Unix() - 946684800)
	return e.writeUint32(n)
}

func (e *Encoder) writeCurrencyName(currecy CurrencyName) (err error) {
	Logger.Encoder.Printf("Writing currency stamp [%s]\n", currecy)
	out := make([]byte, 7, 7)
	copy(out, []byte(currecy))

	return e.toWriter(out)
}

func (e *Encoder) writeAsset(asset Asset) (err error) {

	Logger.Encoder.Printf("Writing asset [%s]\n", asset)
	e.writeUint64(uint64(asset.Amount))
	e.writeByte(asset.Precision)

	symbol := make([]byte, 7, 7)

	copy(symbol[:], []byte(asset.Symbol.Symbol))
	return e.toWriter(symbol)
}

func (e *Encoder) writeJSONTime(time JSONTime) (err error) {
	Logger.Encoder.Printf("Writing json time [%s]\n", time)
	return e.writeUint32(uint32(time.Unix()))
}

func (e *Encoder) writeBlockP2PMessageEnvelope(envelope Packet) (err error) {

	Logger.Encoder.Print("writeBlockP2PMessageEnvelope")

	if envelope.P2PMessage != nil {
		buf := new(bytes.Buffer)
		subEncoder := NewEncoder(buf)
		err = subEncoder.Encode(envelope.P2PMessage)
		if err != nil {
			err = fmt.Errorf("p2p message, %s", err)
			return
		}
		envelope.Payload = buf.Bytes()
	}

	messageLen := uint32(len(envelope.Payload) + 1)
	Logger.Encoder.Printf("Message length: %d\n", messageLen)
	err = e.writeUint32(messageLen)
	if err == nil {
		err = e.writeByte(byte(envelope.Type))

		if err == nil {
			return e.toWriter(envelope.Payload)
		}
	}
	return
}

func (e *Encoder) writeActionData(actionData ActionData) (err error) {
	if actionData.Data != nil {
		//if reflect.TypeOf(actionData.Data) == reflect.TypeOf(&ActionData{}) {
		//	log.Fatal("pas cool")
		//}

		Logger.Encoder.Printf("entering action data, %T\n", actionData)
		var d interface{}
		d = actionData.Data
		if reflect.TypeOf(d).Kind() == reflect.Ptr {
			d = reflect.ValueOf(actionData.Data).Elem().Interface()
		}

		if reflect.TypeOf(d).Kind() == reflect.String { //todo : this is a very bad ack ......

			data, err := hex.DecodeString(d.(string))
			if err != nil {
				return fmt.Errorf("ack, %s", err)
			}
			e.writeByteArray(data)
			return nil

		}

		Logger.Encoder.Printf("encoding action data, %T\n", d)
		raw, err := MarshalBinary(d)
		if err != nil {
			return err
		}
		Logger.Encoder.Printf("writing action data, %T\n", d)
		return e.writeByteArray(raw)
	}

	return e.writeByteArray(actionData.HexData)
}

func MarshalBinary(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := NewEncoder(buf)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}
