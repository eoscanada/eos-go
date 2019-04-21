package eos

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"bytes"

	"time"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoder_Remaining(t *testing.T) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, 1)
	binary.LittleEndian.PutUint16(b[2:], 2)

	d := NewDecoder(b)

	n, err := d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(1), n)
	assert.Equal(t, 2, d.remaining())

	n, err = d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(2), n)
	assert.Equal(t, 0, d.remaining())

}

func TestDecoder_Byte(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeByte(0)
	enc.writeByte(1)

	d := NewDecoder(buf.Bytes())

	n, err := d.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(0), n)
	assert.Equal(t, 1, d.remaining())

	n, err = d.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(1), n)
	assert.Equal(t, 0, d.remaining())

}

func TestDecoder_ByteArray(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeByteArray([]byte{1, 2, 3})
	enc.writeByteArray([]byte{4, 5, 6})

	d := NewDecoder(buf.Bytes())

	data, err := d.ReadByteArray()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, data)
	assert.Equal(t, 4, d.remaining())

	data, err = d.ReadByteArray()
	assert.Equal(t, []byte{4, 5, 6}, data)
	assert.Equal(t, 0, d.remaining())

}

func TestDecoder_ByteArray_MissingData(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUVarInt(10)

	d := NewDecoder(buf.Bytes())

	_, err := d.ReadByteArray()
	assert.EqualError(t, err, "byte array: varlen=10, missing 10 bytes")

}

func TestDecoder_ByteArrayDataTooSmall(t *testing.T) {

	buf := new(bytes.Buffer)

	//to smalls
	d := NewDecoder(buf.Bytes())
	_, err := d.ReadByteArray()
	assert.Equal(t, ErrVarIntBufferSize, err)
}

func TestDecoder_Uint16(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUint16(uint16(99))
	enc.writeUint16(uint16(100))

	d := NewDecoder(buf.Bytes())

	n, err := d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(99), n)
	assert.Equal(t, 2, d.remaining())

	n, err = d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(100), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_int16(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeInt16(int16(-99))
	enc.writeInt16(int16(100))

	d := NewDecoder(buf.Bytes())

	n, err := d.ReadInt16()
	assert.NoError(t, err)
	assert.Equal(t, int16(-99), n)
	assert.Equal(t, 2, d.remaining())

	n, err = d.ReadInt16()
	assert.NoError(t, err)
	assert.Equal(t, int16(100), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Uint32(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUint32(uint32(342))
	enc.writeUint32(uint32(100))

	d := NewDecoder(buf.Bytes())

	n, err := d.ReadUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(342), n)
	assert.Equal(t, 4, d.remaining())

	n, err = d.ReadUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(100), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Int32(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeInt32(int32(-342))
	enc.writeInt32(int32(100))

	d := NewDecoder(buf.Bytes())

	n, err := d.ReadInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(-342), n)
	assert.Equal(t, 4, d.remaining())

	n, err = d.ReadInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(100), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Uint64(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUint64(uint64(99))
	enc.writeUint64(uint64(100))

	d := NewDecoder(buf.Bytes())

	n, err := d.ReadUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(99), n)
	assert.Equal(t, 8, d.remaining())

	n, err = d.ReadUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_string(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeString("123")
	enc.writeString("")
	enc.writeString("abc")

	d := NewDecoder(buf.Bytes())

	s, err := d.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "123", s)
	assert.Equal(t, 5, d.remaining())

	s, err = d.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "", s)
	assert.Equal(t, 4, d.remaining())

	s, err = d.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "abc", s)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Checksum256(t *testing.T) {

	s := Checksum256(bytes.Repeat([]byte{1}, 32))

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeChecksum256(s)

	d := NewDecoder(buf.Bytes())

	rs, err := d.ReadChecksum256()
	assert.NoError(t, err)

	assert.Equal(t, s, rs)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Empty_Checksum256(t *testing.T) {

	s := Checksum256([]byte{})

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeChecksum256(s)

	d := NewDecoder(buf.Bytes())

	s, err := d.ReadChecksum256()
	assert.NoError(t, err)
	assert.Equal(t, s, Checksum256(bytes.Repeat([]byte{0}, 32)))
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_PublicKey(t *testing.T) {

	pk := ecc.MustNewPublicKey("EOS1111111111111111111111111111111114T1Anm")

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	assert.NoError(t, enc.writePublicKey(pk))

	d := NewDecoder(buf.Bytes())

	rpk, err := d.ReadPublicKey()
	assert.NoError(t, err)

	assert.Equal(t, pk, rpk)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_PublicKey_K1(t *testing.T) {
	pk := ecc.MustNewPublicKey("PUB_K1_1111111111111111111111111111111114T1Anm")

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	assert.NoError(t, enc.writePublicKey(pk))

	d := NewDecoder(buf.Bytes())

	rpk, err := d.ReadPublicKey()
	assert.NoError(t, err)

	assert.Equal(t, pk, rpk)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_PublicKey_R1(t *testing.T) {
	pk := ecc.MustNewPublicKey("PUB_R1_81x8BXgDQGTWmcAaavfCDcVTTyzz1BeBYbje9yJomVMCJZbz86")

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	assert.NoError(t, enc.writePublicKey(pk))

	d := NewDecoder(buf.Bytes())

	rpk, err := d.ReadPublicKey()

	require.NoError(t, err)

	assert.Equal(t, pk, rpk)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Empty_PublicKey(t *testing.T) {

	pk := ecc.PublicKey{Curve: ecc.CurveK1, Content: []byte{}}

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	assert.Error(t, enc.writePublicKey(pk))
}

func TestDecoder_Signature(t *testing.T) {

	sig := ecc.MustNewSignatureFromData(bytes.Repeat([]byte{0}, 66))

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeSignature(sig)

	d := NewDecoder(buf.Bytes())

	rsig, err := d.ReadSignature()
	assert.NoError(t, err)
	assert.Equal(t, sig, rsig)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Empty_Signature(t *testing.T) {

	sig := ecc.Signature{Content: []byte{}}

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	assert.Error(t, enc.writeSignature(sig))
}

func TestDecoder_Tstamp(t *testing.T) {

	ts := Tstamp{
		time.Unix(0, time.Now().UnixNano()),
	}

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeTstamp(ts)

	d := NewDecoder(buf.Bytes())

	rts, err := d.ReadTstamp()
	assert.NoError(t, err)
	assert.Equal(t, ts, rts)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_BlockTimestamp(t *testing.T) {
	// Represents block timestamp at slot 1, which is 500 millisecons pass
	// the block epoch which is
	ts := BlockTimestamp{
		time.Unix(0, 500*1000*1000+946684800000*1000*1000),
	}

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeBlockTimestamp(ts)

	// This represents slot 1 in big endian uint32 encoding
	assert.Equal(t, "01000000", hex.EncodeToString(buf.Bytes()))

	d := NewDecoder(buf.Bytes())

	rbt, err := d.ReadBlockTimestamp()
	assert.NoError(t, err)
	assert.Equal(t, ts, rbt)
	assert.Equal(t, 0, d.remaining())
}

type EncodeTestStruct struct {
	F1 string
	F2 int16
	F3 uint16
	F4 uint32
	F5 Checksum256
	F6 []string
	F7 [2]string
	//	F8  map[string]string
	F9  ecc.PublicKey
	F10 ecc.Signature
	F11 byte
	F12 uint64
	F13 []byte
	F14 Tstamp
	F15 BlockTimestamp
	F16 Varuint32
	F17 bool
	F18 Asset
}

func TestDecoder_Encode(t *testing.T) {
	//EnableDecoderLogging()
	//EnableEncoderLogging()

	now := time.Date(2018, time.September, 26, 1, 2, 3, 4, time.UTC)
	tstamp := Tstamp{Time: time.Unix(0, now.UnixNano())}
	blockts := BlockTimestamp{time.Unix(now.Unix(), 0)}
	s := &EncodeTestStruct{
		F1: "abc",
		F2: -75,
		F3: 99,
		F4: 999,
		F5: bytes.Repeat([]byte{0}, 32),
		F6: []string{"def", "789"},
		F7: [2]string{"foo", "bar"},
		// maps don't serialize deterministically.. we no want that.
		//		F8:  map[string]string{"foo": "bar", "hello": "you"},
		F9:  ecc.MustNewPublicKey("EOS1111111111111111111111111111111114T1Anm"),
		F10: ecc.Signature{Curve: ecc.CurveK1, Content: make([]byte, 65)},
		F11: byte(1),
		F12: uint64(87),
		F13: []byte{1, 2, 3, 4, 5},
		F14: tstamp,
		F15: blockts,
		F16: Varuint32(999),
		F17: true,
		F18: NewEOSAsset(100000),
	}

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	assert.NoError(t, enc.Encode(s))

	assert.Equal(t, "03616263b5ff6300e7030000000000000000000000000000000000000000000000000000000000000000000002036465660337383903666f6f036261720000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001570000000000000005010203040504ae0f517acd5715162e7b46e70701a08601000000000004454f5300000000", hex.EncodeToString(buf.Bytes()))

	decoder := NewDecoder(buf.Bytes())
	assert.NoError(t, decoder.Decode(s))

	assert.Equal(t, "abc", s.F1)
	assert.Equal(t, int16(-75), s.F2)
	assert.Equal(t, uint16(99), s.F3)
	assert.Equal(t, uint32(999), s.F4)
	assert.Equal(t, Checksum256(bytes.Repeat([]byte{0}, 32)), s.F5)
	assert.Equal(t, []string{"def", "789"}, s.F6)
	assert.Equal(t, [2]string{"foo", "bar"}, s.F7)
	//	assert.Equal(t, map[string]string{"foo": "bar", "hello": "you"}, s.F8)
	assert.Equal(t, ecc.MustNewPublicKeyFromData(bytes.Repeat([]byte{0}, 34)), s.F9)
	assert.Equal(t, ecc.MustNewSignatureFromData(bytes.Repeat([]byte{0}, 66)), s.F10)
	assert.Equal(t, byte(1), s.F11)
	assert.Equal(t, uint64(87), s.F12)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, s.F13)
	assert.Equal(t, tstamp, s.F14)
	assert.Equal(t, blockts, s.F15)
	assert.Equal(t, Varuint32(999), s.F16)
	assert.Equal(t, true, s.F17)
	assert.Equal(t, Int64(100000), s.F18.Amount)
	assert.Equal(t, uint8(4), s.F18.Precision)
	assert.Equal(t, "EOS", s.F18.Symbol.Symbol)

}

func TestDecoder_Decode_No_Ptr(t *testing.T) {
	decoder := NewDecoder([]byte{})
	err := decoder.Decode(1)
	assert.EqualError(t, err, "decode, can only Decode to pointer type")
}

func TestDecoder_Decode_String_Err(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.writeUVarInt(10)

	decoder := NewDecoder(buf.Bytes())
	var s string
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "byte array: varlen=10, missing 10 bytes")
}

func TestDecoder_Decode_Array(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.Encode([3]byte{1, 2, 4})

	assert.Equal(t, []byte{1, 2, 4}, buf.Bytes())

	decoder := NewDecoder(buf.Bytes())
	var decoded [3]byte
	decoder.Decode(&decoded)
	assert.Equal(t, [3]byte{1, 2, 4}, decoded)

}

func TestDecoder_Decode_Slice_Err(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)

	decoder := NewDecoder(buf.Bytes())
	var s []string
	err := decoder.Decode(&s)
	assert.Equal(t, err, ErrVarIntBufferSize)

	enc.writeUVarInt(1)
	decoder = NewDecoder(buf.Bytes())
	err = decoder.Decode(&s)
	assert.Equal(t, err, ErrVarIntBufferSize)
}

type structWithInvalidType struct {
	F1 time.Duration
}

func TestDecoder_Decode_Struct_Err(t *testing.T) {

	s := structWithInvalidType{}
	decoder := NewDecoder([]byte{})
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "decode, unsupported type time.Duration")

}

func TestEncoder_Encode_array_error(t *testing.T) {

	decoder := NewDecoder([]byte{1})

	toDecode := [1]time.Duration{}
	err := decoder.Decode(&toDecode)

	assert.EqualError(t, err, "decode, unsupported type time.Duration")

}

func TestEncoder_Decode_array_error(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	err := enc.Encode([1]time.Duration{time.Duration(0)})
	assert.EqualError(t, err, "Encode: unsupported type time.Duration")

}

func TestEncoder_Encode_slide_error(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	err := enc.Encode([]time.Duration{time.Duration(0)})
	assert.EqualError(t, err, "Encode: unsupported type time.Duration")

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
	assert.EqualError(t, err, "Encode: unsupported type time.Duration")

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

	d := NewDecoder(buf.Bytes())
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
	assert.Equal(t, expected, buf.Bytes())

}

func TestDecoder_readUint16_missing_data(t *testing.T) {

	_, err := NewDecoder([]byte{}).ReadByte()
	assert.EqualError(t, err, "byte required [1] byte, remaining [0]")

	_, err = NewDecoder([]byte{}).ReadUint16()
	assert.EqualError(t, err, "uint16 required [2] bytes, remaining [0]")

	_, err = NewDecoder([]byte{}).ReadUint32()
	assert.EqualError(t, err, "uint32 required [4] bytes, remaining [0]")

	_, err = NewDecoder([]byte{}).ReadUint64()
	assert.EqualError(t, err, "uint64 required [8] bytes, remaining [0]")

	_, err = NewDecoder([]byte{}).ReadChecksum256()
	assert.EqualError(t, err, "checksum 256 required [32] bytes, remaining [0]")

	_, err = NewDecoder([]byte{}).ReadPublicKey()
	assert.EqualError(t, err, "publicKey required [34] bytes, remaining [0]")

	_, err = NewDecoder([]byte{}).ReadSignature()
	assert.EqualError(t, err, "signature required [66] bytes, remaining [0]")

	_, err = NewDecoder([]byte{}).ReadTstamp()
	assert.EqualError(t, err, "tstamp required [8] bytes, remaining [0]")

	_, err = NewDecoder([]byte{}).ReadBlockTimestamp()
	assert.EqualError(t, err, "blockTimestamp required [4] bytes, remaining [0]")
}

func TestDecoder_SignedBlock_Full(t *testing.T) {
	dataHex := "1b146b480000000000ea305500000000000140215a6edeea1e697207b5a917d83edf56a963d03e3d5d8d8e1ddb0900000000000000000000000000000000000000000000000000000000000000006a46611d7b15f71ff42de916e19f8ed1011096178f81d9b17987637a545b152100000000000100002001000000000000000000000000000000000000000000000000000000000000fe001f5e6962745bb4fb84dec1e7779b7e3b58c5fe20ed39c41cac1bdefe3e71568bf67bdfb05e4df40087c9ecbc6d9d0f6bcfcddfdc1d00dc1d035c665936307535ab0001000020fe00000000000000000000000000000000000000000000000000000000000004"
	data, err := hex.DecodeString(dataHex)
	require.NoError(t, err)

	var signedBlock *SignedBlock
	err = UnmarshalBinary(data, &signedBlock)
	require.NoError(t, err)

	expectedTimestamp, _ := time.Parse("2006-01-02T15:04:05.999-0700", "2019-04-01T22:48:45.500-0400")

	// TODO: I'm not quite sure about endiannes of this value, I would have though it should have been
	//       equal to `fe00000000000000000000000000000000000000000000000000000000000001` as in the
	//       nodeos binary data, which is usually all big endian, it is written as
	//       `01000000000000000000000000000000000000000000000000000000000000fe`.
	//
	//       Our current real data has only 0s so it's impossible to tell right endiannes. We would
	//       need to craft a block with some data in it or search `nodeos` to validate how a `vector<char>`
	//       is written to binary.
	//
	//       Same reasoning apply to both []*Extension fields
	expectedHeaderExtension, _ := hex.DecodeString("01000000000000000000000000000000000000000000000000000000000000fe")
	expectedBlockExtension, _ := hex.DecodeString("fe00000000000000000000000000000000000000000000000000000000000004")

	assert.Equal(t, BlockTimestamp{expectedTimestamp.Local()}, signedBlock.Timestamp)
	assert.Equal(t, AccountName("eosio"), signedBlock.Producer)
	assert.Equal(t, uint16(0), signedBlock.Confirmed)
	assert.Equal(t, "0000000140215a6edeea1e697207b5a917d83edf56a963d03e3d5d8d8e1ddb09", signedBlock.Previous.String())
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", signedBlock.TransactionMRoot.String())
	assert.Equal(t, "6a46611d7b15f71ff42de916e19f8ed1011096178f81d9b17987637a545b1521", signedBlock.ActionMRoot.String())
	assert.Equal(t, uint32(0), signedBlock.ScheduleVersion)
	assert.Equal(t, (*OptionalProducerSchedule)(nil), signedBlock.NewProducers)
	assert.Equal(t, []*Extension{&Extension{uint16(0), expectedHeaderExtension}}, signedBlock.HeaderExtensions)
	assert.Equal(t, "SIG_K1_K7cBDNuka9kLUNAGaCm4FpNTdJwVKY3rP3v2esU8RGv1KXNNDEEdrWBAJSH3cPB8t1478e4RmhjkP48Sbuaqkf6Z5iDZKW", signedBlock.ProducerSignature.String())
	assert.Equal(t, []TransactionReceipt{}, signedBlock.Transactions)
	assert.Equal(t, []*Extension{&Extension{uint16(0), expectedBlockExtension}}, signedBlock.BlockExtensions)
}
