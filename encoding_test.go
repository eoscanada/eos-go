package eos

import (
	"encoding/binary"
	"testing"

	"bytes"

	"time"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/stretchr/testify/assert"
)

func TestDecoder_Remaining(t *testing.T) {

	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, 1)
	binary.LittleEndian.PutUint16(b[2:], 2)

	d := NewDecoder(b)

	assert.Equal(t, uint16(1), d.readUint16())
	assert.Equal(t, 2, d.remaining())
	assert.Equal(t, uint16(2), d.readUint16())
	assert.Equal(t, 0, d.remaining())

}

func TestDecoder_Byte(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeByte(0)
	enc.writeByte(1)

	d := NewDecoder(enc.data)

	assert.Equal(t, byte(0), d.readByte())
	assert.Equal(t, 1, d.remaining())
	assert.Equal(t, byte(1), d.readByte())
	assert.Equal(t, 0, d.remaining())

}

func TestDecoder_ByteArray(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeByteArray([]byte{1, 2, 3})
	enc.writeByteArray([]byte{4, 5, 6})

	d := NewDecoder(enc.data)

	data, err := d.readByteArray()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, data)
	assert.Equal(t, 4, d.remaining())

	data, err = d.readByteArray()
	assert.Equal(t, []byte{4, 5, 6}, data)
	assert.Equal(t, 0, d.remaining())

}

func TestDecoder_ByteArray_MissingData(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUVarInt(10)

	d := NewDecoder(enc.data)

	_, err := d.readByteArray()
	assert.EqualError(t, err, "byte array: varlen=10, missing 10 bytes")

}

func TestDecoder_ByteArrayDataTooSmall(t *testing.T) {

	buf := new(bytes.Buffer)

	//to smalls
	d := NewDecoder(buf.Bytes())
	_, err := d.readByteArray()
	assert.Equal(t, VarIntBufferSizeError, err)

}
func TestDecoder_Uint16(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUint16(uint16(99))
	enc.writeUint16(uint16(100))

	d := NewDecoder(enc.data)

	assert.Equal(t, uint16(99), d.readUint16())
	assert.Equal(t, 2, d.remaining())
	assert.Equal(t, uint16(100), d.readUint16())
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_int16(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeInt16(int16(-99))
	enc.writeInt16(int16(100))

	d := NewDecoder(enc.data)

	assert.Equal(t, int16(-99), d.readInt16())
	assert.Equal(t, 2, d.remaining())
	assert.Equal(t, int16(100), d.readInt16())
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Uint32(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUint32(uint32(99))
	enc.writeUint32(uint32(100))

	d := NewDecoder(enc.data)

	assert.Equal(t, uint32(99), d.readUint32())
	assert.Equal(t, 4, d.remaining())
	assert.Equal(t, uint32(100), d.readUint32())
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Uint64(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUint64(uint64(99))
	enc.writeUint64(uint64(100))

	d := NewDecoder(enc.data)

	assert.Equal(t, uint64(99), d.readUint64())
	assert.Equal(t, 8, d.remaining())
	assert.Equal(t, uint64(100), d.readUint64())
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_string(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeString("123")
	enc.writeString("")
	enc.writeString("abc")

	d := NewDecoder(enc.data)

	s, err := d.readString()
	assert.NoError(t, err)
	assert.Equal(t, "123", s)
	assert.Equal(t, 5, d.remaining())

	s, err = d.readString()
	assert.NoError(t, err)
	assert.Equal(t, "", s)
	assert.Equal(t, 4, d.remaining())

	s, err = d.readString()
	assert.NoError(t, err)
	assert.Equal(t, "abc", s)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_SHA256Bytes(t *testing.T) {

	s := SHA256Bytes(bytes.Repeat([]byte{1}, 32))

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeSHA256Bytes(s)

	d := NewDecoder(enc.data)

	assert.Equal(t, s, d.readSHA256Bytes())
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Empty_SHA256Bytes(t *testing.T) {

	s := SHA256Bytes([]byte{})

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeSHA256Bytes(s)

	d := NewDecoder(enc.data)

	s = d.readSHA256Bytes()
	assert.Equal(t, s, SHA256Bytes(bytes.Repeat([]byte{0}, 32)))
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_PublicKey(t *testing.T) {

	pk := ecc.PublicKey(bytes.Repeat([]byte{1}, 34))

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writePublicKey(pk)

	d := NewDecoder(enc.data)

	assert.Equal(t, pk, d.readPublicKey())
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Empty_PublicKey(t *testing.T) {

	pk := ecc.PublicKey([]byte{})

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writePublicKey(pk)

	d := NewDecoder(enc.data)

	pk = d.readPublicKey()
	assert.Equal(t, pk, ecc.PublicKey(bytes.Repeat([]byte{0}, 34)))
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Signature(t *testing.T) {

	sig := ecc.Signature(bytes.Repeat([]byte{1}, 66))

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeSignature(sig)

	d := NewDecoder(enc.data)

	assert.Equal(t, sig, d.readSignature())
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Empty_Signature(t *testing.T) {

	sig := ecc.Signature([]byte{})

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeSignature(sig)

	d := NewDecoder(enc.data)

	sig = d.readSignature()
	assert.Equal(t, sig, ecc.Signature(bytes.Repeat([]byte{0}, 66)))
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Tstamp(t *testing.T) {

	ts := Tstamp{
		time.Unix(0, time.Now().UnixNano()),
	}

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeTstamp(ts)

	d := NewDecoder(enc.data)

	assert.Equal(t, ts, d.readTstamp())
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_BlockTimestamp(t *testing.T) {

	ts := BlockTimestamp{
		time.Unix(time.Now().Unix(), 0),
	}

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeBlockTimestamp(ts)

	d := NewDecoder(enc.data)

	assert.Equal(t, ts, d.readBlockTimestamp())
	assert.Equal(t, 0, d.remaining())
}

type EncodeTestStruct struct {
	F1  string
	F2  int16
	F3  uint16
	F4  uint32
	F5  SHA256Bytes
	F6  []string
	F7  [2]string
	F8  map[string]string
	F9  ecc.PublicKey
	F10 ecc.Signature
	F11 byte
	F12 uint64
	F13 []byte
	F14 Tstamp
	F15 BlockTimestamp
	F16 Varuint32
}

func TestDecoder_Encode(t *testing.T) {

	tstamp := Tstamp{Time: time.Unix(0, time.Now().UnixNano())}
	blockts := BlockTimestamp{time.Unix(time.Now().Unix(), 0)}
	s := &EncodeTestStruct{
		F1:  "abc",
		F2:  -75,
		F3:  99,
		F4:  999,
		F5:  bytes.Repeat([]byte{0}, 32),
		F6:  []string{"def", "789"},
		F7:  [2]string{"foo", "bar"},
		F8:  map[string]string{"foo": "bar", "hello": "you"},
		F9:  ecc.PublicKey(bytes.Repeat([]byte{0}, 34)),
		F10: ecc.Signature(bytes.Repeat([]byte{0}, 66)),
		F11: byte(1),
		F12: uint64(87),
		F13: []byte{1, 2, 3, 4, 5},
		F14: tstamp,
		F15: blockts,
		F16: Varuint32(999),
	}

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.Encode(s)

	decoder := NewDecoder(enc.data)
	err := decoder.Decode(s)
	assert.NoError(t, err)

	assert.Equal(t, "abc", s.F1)
	assert.Equal(t, int16(-75), s.F2)
	assert.Equal(t, uint16(99), s.F3)
	assert.Equal(t, uint32(999), s.F4)
	assert.Equal(t, SHA256Bytes(bytes.Repeat([]byte{0}, 32)), s.F5)
	assert.Equal(t, []string{"def", "789"}, s.F6)
	assert.Equal(t, [2]string{"foo", "bar"}, s.F7)
	assert.Equal(t, map[string]string{"foo": "bar", "hello": "you"}, s.F8)
	assert.Equal(t, ecc.PublicKey(bytes.Repeat([]byte{0}, 34)), s.F9)
	assert.Equal(t, ecc.Signature(bytes.Repeat([]byte{0}, 66)), s.F10)
	assert.Equal(t, byte(1), s.F11)
	assert.Equal(t, uint64(87), s.F12)
	assert.Equal(t, uint64(87), s.F12)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, s.F13)
	assert.Equal(t, tstamp, s.F14)
	assert.Equal(t, blockts, s.F15)
	assert.Equal(t, Varuint32(999), s.F16)

}

func TestDecoder_Decode_No_Ptr(t *testing.T) {
	decoder := NewDecoder([]byte{})
	err := decoder.Decode(1)
	assert.EqualError(t, err, "decode: can only Decode to pointer type")
}

func TestDecoder_Decode_String_Err(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUVarInt(10)

	decoder := NewDecoder(enc.data)
	var s string
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "byte array: varlen=10, missing 10 bytes")
}

func TestDecoder_Decode_Array_Err(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)

	decoder := NewDecoder(enc.data)
	var s [1]string
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "varint: invalide buffer size")

	enc.writeUVarInt(1)
	decoder = NewDecoder(enc.data)
	err = decoder.Decode(&s)
	assert.EqualError(t, err, "varint: invalide buffer size")

}

func TestDecoder_Decode_Slice_Err(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)

	decoder := NewDecoder(enc.data)
	var s []string
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "varint: invalide buffer size")

	enc.writeUVarInt(1)
	decoder = NewDecoder(enc.data)
	err = decoder.Decode(&s)
	assert.EqualError(t, err, "varint: invalide buffer size")

}

type structWithInvalidType struct {
	F1 time.Duration
}

func TestDecoder_Decode_Struct_Err(t *testing.T) {

	s := structWithInvalidType{}
	decoder := NewDecoder([]byte{})
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "binary: unsupported type time.Duration")

}

func TestDecoder_Decode_Map_Err(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)

	decoder := NewDecoder(enc.data)
	var m map[string]string
	err := decoder.Decode(&m)
	assert.EqualError(t, err, "varint: invalide buffer size")

	enc.writeUVarInt(1)
	decoder = NewDecoder(enc.data)
	err = decoder.Decode(&m)
	assert.EqualError(t, err, "varint: invalide buffer size")
}

func TestDecoder_Decode_Bad_Map(t *testing.T) {

	buf := new(bytes.Buffer)
	var m map[string]time.Duration
	enc := NewEncoder(buf)
	enc.writeUVarInt(1)
	enc.writeString("foo")
	enc.writeString("bar")

	decoder := NewDecoder(enc.data)
	err := decoder.Decode(&m)
	assert.EqualError(t, err, "binary: unsupported type time.Duration")

}

func TestEncoder_Encode_array_error(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	err := enc.Encode([1]time.Duration{time.Duration(0)})
	assert.EqualError(t, err, "binary: unsupported type time.Duration")

}
func TestEncoder_Encode_slide_error(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	err := enc.Encode([]time.Duration{time.Duration(0)})
	assert.EqualError(t, err, "binary: unsupported type time.Duration")

}
func TestEncoder_Encode_map_error(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	err := enc.Encode(map[string]time.Duration{"key": time.Duration(0)})
	assert.EqualError(t, err, "binary: unsupported type time.Duration")
	err = enc.Encode(map[time.Duration]string{time.Duration(0): "key"})
	assert.EqualError(t, err, "binary: unsupported type time.Duration")

}

func TestEncoder_Encode_struct_error(t *testing.T) {

	s := struct {
		F time.Duration
	}{
		F: time.Duration(0),
	}

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	err := enc.Encode(&s)
	assert.EqualError(t, err, "binary: unsupported type time.Duration")

}

type TagTestStruct struct {
	S1 string `eos:"-"`
	S2 string
}

func TestEncoder_Decode_struct_tag(t *testing.T) {
	var s TagTestStruct

	buf := new(bytes.Buffer)

	enc := NewEncoder(buf)
	enc.writeString("123")

	d := NewDecoder(enc.data)
	d.Decode(&s)
	assert.Equal(t, "", s.S1)
	assert.Equal(t, "123", s.S2)

}

func TestEncoder_Encode_struct_tag(t *testing.T) {

	s := &TagTestStruct{
		S1: "123",
		S2: "abc",
	}

	buf := new(bytes.Buffer)

	enc := NewEncoder(buf)
	enc.Encode(s)

	expected := []byte{0x3, 0x61, 0x62, 0x63}
	assert.Equal(t, expected, enc.data)

}
