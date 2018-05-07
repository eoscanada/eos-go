package eos

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"

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

func (e *Encoder) writeName(name Name) (err error) {

	val, er := StringToName(string(name))
	if er != nil {
		err = fmt.Errorf("encode, name, %s", er)
		return
	}
	err = e.writeUint64(val)
	return

}
func (e *Encoder) Encode(v interface{}) (err error) {
	switch cv := v.(type) {
	case Name:
		err = e.writeName(cv)
		return
	case AccountName:
		name := Name(cv)
		err = e.writeName(name)
		return
	case PermissionName:
		name := Name(cv)
		err = e.writeName(name)
		return
	case ActionName:
		name := Name(cv)
		err = e.writeName(name)
		return
	case TableName:
		name := Name(cv)
		err = e.writeName(name)
		return
	case ScopeName:
		name := Name(cv)
		err = e.writeName(name)
		return
	case string:
		err = e.writeString(cv)
		return
	case byte:
		err = e.writeByte(cv)
		return
	case int16:
		err = e.writeInt16(cv)
		return
	case uint16:
		err = e.writeUint16(cv)
		return
	case uint32:
		err = e.writeUint32(cv)
		return
	case uint64:
		err = e.writeUint64(cv)
		return
	case Varuint32:
		err = e.writeUVarInt(int(cv))
		return
	case JSONTime:
		err = e.writeJsonTime(cv)
		return
	case HexBytes:
		err = e.writeByteArray(cv)
		return
	case []byte:
		err = e.writeByteArray(cv)
		return
	case SHA256Bytes:
		err = e.writeSHA256Bytes(cv)
		return
	case ecc.PublicKey:
		err = e.writePublicKey(cv)
		return
	case ecc.Signature:
		err = e.writeSignature(cv)
		return
	case Tstamp:
		err = e.writeTstamp(cv)
		return
	case BlockTimestamp:
		err = e.writeBlockTimestamp(cv)
		return
	case CurrencyName:
		err = e.writeCurrencyName(cv)
		return
	case Asset:
		err = e.writeAsset(cv)
		return
	case ActionData:
		err = e.writeActionData(cv)
		return
	case *P2PMessageEnvelope:
		err = e.writeBlockP2PMessageEnvelope(*cv)
		return
	default:

		rv := reflect.Indirect(reflect.ValueOf(v))
		t := rv.Type()

		switch t.Kind() {

		case reflect.Array:
			l := t.Len()
			prefix = append(prefix, "     ")
			println(fmt.Sprintf("Encode: array [%T] of length: %d", v, l))

			for i := 0; i < l; i++ {
				if err = e.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}
			prefix = prefix[:len(prefix)-1]
		case reflect.Slice:
			l := rv.Len()
			e.writeUVarInt(l)
			prefix = append(prefix, "     ")
			println(fmt.Sprintf("Encode: slice [%T] of length: %d", v, l))

			for i := 0; i < l; i++ {
				if err = e.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}
			prefix = prefix[:len(prefix)-1]
		case reflect.Struct:
			l := rv.NumField()
			println(fmt.Sprintf("Encode: struct [%T] with %d field.", v, l))
			prefix = append(prefix, "     ")

			n := 0
			for i := 0; i < l; i++ {
				field := t.Field(i)
				println(fmt.Sprintf("field -> %s", field.Name))

				if tag := field.Tag.Get("eos"); tag == "-" {
					continue
				}
				if v := rv.Field(i); t.Field(i).Name != "_" {
					if v.CanInterface() {
						iface := v.Interface()
						if iface != nil {
							if err = e.Encode(iface); err != nil {
								return
							}
						}
					}
					n++
				}
			}
			prefix = prefix[:len(prefix)-1]

		case reflect.Map:
			l := rv.Len()
			e.writeUVarInt(l)
			println(fmt.Sprintf("Map [%T] of length: %d", v, l))
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
			return errors.New("binary: unsupported type " + t.String())
		}
	}

	return
}

func (e *Encoder) toWriter(bytes []byte) (err error) {

	e.count += len(bytes)
	println(fmt.Sprintf("    Appending : [%s] pos [%d]", hex.EncodeToString(bytes), e.count))
	_, err = e.output.Write(bytes)
	return
}

func (e *Encoder) writeByteArray(b []byte) error {
	e.writeUVarInt(len(b))
	return e.toWriter(b)
}

func (e *Encoder) writeUVarInt(v int) (err error) {
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return e.toWriter(buf[:l])
}

func (e *Encoder) writeByte(b byte) (err error) {
	return e.toWriter([]byte{b})
}

func (e *Encoder) writeUint16(i uint16) (err error) {
	buf := make([]byte, TypeSize.UInt16)
	binary.LittleEndian.PutUint16(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeInt16(i int16) (err error) {
	return e.writeUint16(uint16(i))
}

func (e *Encoder) writeUint32(i uint32) (err error) {
	buf := make([]byte, TypeSize.UInt32)
	binary.LittleEndian.PutUint32(buf, i)
	return e.toWriter(buf)

}

func (e *Encoder) writeUint64(i uint64) (err error) {
	buf := make([]byte, TypeSize.UInt64)
	binary.LittleEndian.PutUint64(buf, i)
	return e.toWriter(buf)

}

func (e *Encoder) writeString(s string) (err error) {
	return e.writeByteArray([]byte(s))
}

func (e *Encoder) writeSHA256Bytes(s SHA256Bytes) error {
	if len(s) == 0 {
		return e.toWriter(bytes.Repeat([]byte{0}, TypeSize.SHA256Bytes))
	}
	return e.toWriter(s)
}

func (e *Encoder) writePublicKey(pk ecc.PublicKey) (err error) {
	if len(pk) == 0 {
		return e.toWriter(bytes.Repeat([]byte{0}, TypeSize.PublicKey))
	}

	return e.toWriter(append(bytes.Repeat([]byte{0}, 34-len(pk)), pk...))
}

func (e *Encoder) writeSignature(s ecc.Signature) (err error) {
	if len(s) == 0 {
		return e.toWriter(bytes.Repeat([]byte{0}, TypeSize.Signature))
	}
	return e.toWriter(s)
}

func (e *Encoder) writeTstamp(t Tstamp) (err error) {
	n := uint64(t.UnixNano())
	return e.writeUint64(n)
}

func (e *Encoder) writeBlockTimestamp(bt BlockTimestamp) (err error) {
	n := uint32(bt.Unix() - 946684800)
	return e.writeUint32(n)
}

func (e *Encoder) writeCurrencyName(curreny CurrencyName) (err error) {
	out := make([]byte, 7, 7)
	copy(out, []byte(curreny))

	return e.toWriter(out)
}

func (e *Encoder) writeAsset(asset Asset) (err error) {

	e.writeUint64(uint64(asset.Amount))
	e.writeByte(asset.Precision)

	symbol := make([]byte, 7, 7)

	copy(symbol[:], []byte(asset.Symbol.Symbol))
	return e.toWriter(symbol)
}

func (e *Encoder) writeJsonTime(time JSONTime) (err error) {
	return e.writeUint32(uint32(time.Unix()))
}

func (e *Encoder) writeBlockP2PMessageEnvelope(envelope P2PMessageEnvelope) (err error) {

	println("writeBlockP2PMessageEnvelope")

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
	println(fmt.Sprintf("Message length: %d", messageLen))
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

	if actionData.obj != nil {

		raw, err := MarshalBinary(actionData.obj)
		if err != nil {
			return err
		}
		e.writeByteArray(raw)

	} else {
		e.writeByteArray([]byte{})
	}
	return
}

func MarshalBinary(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := NewEncoder(buf)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}
