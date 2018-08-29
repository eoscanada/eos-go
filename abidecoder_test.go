package eos

import (
	"testing"

	"fmt"

	"bytes"

	"time"

	"strings"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/stretchr/testify/assert"
)

var abiString = `
{
	"version": "eosio::abi/1.0",
	"types": [{
		"new_type_name": "new.type.name.1",
		"type": "name"
	}],
	"structs": [
	{
		"name": "struct.name.1",
		"base": "struct.name.2",
		"fields": [
			{"name":"struct.1.field.1", "type":"new.type.name.1"},
			{"name":"struct.1.field.2", "type":"struct.name.3"}
		]
    },{
		"name": "struct.name.2",
		"base": "",
		"fields": [
			{"name":"struct.2.field.1", "type":"string"}
		]
    },{
		"name": "struct.name.3",
		"base": "",
		"fields": [
			{"name":"struct.3.field.1", "type":"string"}
		]
    }
	],
   "actions": [{
		"name": "action.name.1",
		"type": "struct.name.1",
		"ricardian_contract": ""
   }]
}
`

func TestABI_Decode(t *testing.T) {

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1 string
		F1  Name
		F2  string
	}{
		BF1: "value.struct.2.field.1",
		F1:  Name("eoscanadacom"),
		F2:  "value.struct.3.field.1",
	}

	var b bytes.Buffer
	encoder := NewEncoder(&b)
	err := encoder.Encode(mockData)
	assert.NoError(t, err)

	decoder := NewABIDecoder(b.Bytes(), abiReader)
	result := map[string]interface{}{}
	err = decoder.Decode(result, "action.name.1")
	assert.NoError(t, err)

	assert.Equal(t, Name("eoscanadacom"), result["struct.1.field.1"])
	assert.Equal(t, "value.struct.2.field.1", result["struct.2.field.1"])
	assert.Equal(t, "value.struct.3.field.1", result["struct.3.field.1"])

}

func TestABI_DecodeMissingData(t *testing.T) {

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1 string
		F1  Name
		//F2  string
	}{
		BF1: "value.struct.2.field.1",
		F1:  Name("eoscanadacom"),
		//F2:  "value.struct.3.field.1",
	}

	var b bytes.Buffer
	encoder := NewEncoder(&b)
	err := encoder.Encode(mockData)
	assert.NoError(t, err)

	decoder := NewABIDecoder(b.Bytes(), abiReader)
	result := map[string]interface{}{}
	err = decoder.Decode(result, "action.name.1")
	assert.Equal(t, fmt.Errorf("decode field [struct.3.field.1] of type[string]: varint: invalid buffer size"), err)

}

func TestABI_DecodeMissingAction(t *testing.T) {

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1 string
		F1  Name
	}{
		BF1: "value.base.field.1",
		F1:  Name("eoscanadacom"),
	}

	var b bytes.Buffer
	encoder := NewEncoder(&b)
	err := encoder.Encode(mockData)
	assert.NoError(t, err)

	decoder := NewABIDecoder(b.Bytes(), abiReader)
	result := map[string]interface{}{}
	err = decoder.Decode(result, "bad.action.name")
	assert.Equal(t, fmt.Errorf("action bad.action.name not found in abi"), err)
}

func TestABI_DecodeBadABI(t *testing.T) {

	abiReader := strings.NewReader("{")

	var b bytes.Buffer

	decoder := NewABIDecoder(b.Bytes(), abiReader)
	result := map[string]interface{}{}
	err := decoder.Decode(result, "bad.action.name")
	assert.Equal(t, fmt.Errorf("read abi: unexpected EOF"), err)
}

func TestABI_decode(t *testing.T) {

	abi := &ABI{
		Structs: []StructDef{
			{
				Name: "struct.base.1",
				Fields: []FieldDef{
					{Name: "base.field.1", Type: "string"},
				},
			},
			{
				Name: "struct.1",
				Base: "struct.base.1",
				Fields: []FieldDef{
					{Name: "field.1", Type: "string"},
				},
			},
		},
	}

	s := struct {
		BF1 string
		F1  string
	}{
		BF1: "value.base.field.1",
		F1:  "value.field.1",
	}

	var b bytes.Buffer
	encoder := NewEncoder(&b)
	err := encoder.Encode(s)
	assert.NoError(t, err)

	decoder := NewABIDecoder(b.Bytes(), nil)
	decoder.abi = abi
	result := make(Result)
	err = decoder.decode("struct.1", result)
	assert.NoError(t, err)

	assert.Equal(t, "value.field.1", result["field.1"])
	assert.Equal(t, "value.base.field.1", result["base.field.1"])

}

func TestABI_decodeStructNotFound(t *testing.T) {

	abi := &ABI{
		Structs: []StructDef{
			{
				Name: "struct.1",
				Base: "struct.base.1",
				Fields: []FieldDef{
					{Name: "field.1", Type: "string"},
				},
			},
		},
	}

	s := struct{}{}

	var b bytes.Buffer
	encoder := NewEncoder(&b)
	err := encoder.Encode(s)
	assert.NoError(t, err)

	decoder := NewABIDecoder(b.Bytes(), nil)
	decoder.abi = abi
	result := make(Result)
	err = decoder.decode("struct.1", result)
	assert.Equal(t, fmt.Errorf("decode base [struct.1]: structure [struct.base.1] not found in abi"), err)
}

func TestABI_decodeStructBaseNotFound(t *testing.T) {

	abi := &ABI{
		Structs: []StructDef{},
	}

	s := struct{}{}

	var b bytes.Buffer
	encoder := NewEncoder(&b)
	err := encoder.Encode(s)
	assert.NoError(t, err)

	decoder := NewABIDecoder(b.Bytes(), nil)
	decoder.abi = abi
	result := make(Result)
	err = decoder.decode("struct.1", result)
	assert.Equal(t, fmt.Errorf("structure [struct.1] not found in abi"), err)
}

func TestABI_decodeFields(t *testing.T) {

	types := []ABIType{
		{NewTypeName: "action.type.1", Type: "name"},
	}
	fields := []FieldDef{
		{Name: "F1", Type: "int8"},
		{Name: "F2", Type: "action.type.1"},
	}
	abi := &ABI{
		Types: types,
		Structs: []StructDef{
			{Fields: fields},
		},
	}

	s := struct {
		F1 int8
		F2 Name
	}{
		F1: int8(10),
		F2: Name("action.name.1"),
	}

	var b bytes.Buffer
	encoder := NewEncoder(&b)
	err := encoder.Encode(s)
	assert.NoError(t, err)

	decoder := NewABIDecoder(b.Bytes(), nil)
	decoder.abi = abi
	result := make(Result)
	err = decoder.decodeFields(fields, result)
	assert.NoError(t, err)

	assert.Equal(t, int8(10), result["F1"])
	assert.Equal(t, Name("action.name.1"), result["F2"])
}

func TestABI_decodeFieldsErr(t *testing.T) {

	types := []ABIType{}
	fields := []FieldDef{
		{
			Name: "field.with.bad.type.1",
			Type: "bad.type.1",
		},
	}

	s := struct{}{}

	abi := &ABI{
		Types: types,
		Structs: []StructDef{
			{Fields: fields},
		},
	}

	var b bytes.Buffer
	encoder := NewEncoder(&b)
	err := encoder.Encode(s)
	assert.NoError(t, err)

	decoder := NewABIDecoder(b.Bytes(), nil)
	decoder.abi = abi
	result := make(Result)

	err = decoder.decodeFields(fields, result)
	assert.Equal(t, fmt.Errorf("decode field [field.with.bad.type.1] of type[bad.type.1]: read field [field.with.bad.type.1] of type [bad.type.1]: unknown type"), err)

}

func TestABI_Read(t *testing.T) {

	bt := BlockTimestamp{
		Time: time.Unix(time.Now().Unix(), 0),
	}

	testCases := []map[string]interface{}{
		{"typeName": "int8", "value": int8(1), "encode": int8(1), "expectedError": nil},
		{"typeName": "uint8", "value": uint8(1), "encode": uint8(1), "expectedError": nil},
		{"typeName": "int16", "value": int16(1), "encode": int16(1), "expectedError": nil},
		{"typeName": "uint16", "value": uint16(1), "encode": uint16(1), "expectedError": nil},
		{"typeName": "int32", "value": int32(1), "encode": int32(1), "expectedError": nil},
		{"typeName": "uint32", "value": uint32(1), "encode": uint32(1), "expectedError": nil},
		{"typeName": "int64", "value": int64(1), "encode": int64(1), "expectedError": nil},
		{"typeName": "uint64", "value": uint64(1), "encode": uint64(1), "expectedError": nil},
		{"typeName": "int128", "value": int64(1), "encode": int64(1), "expectedError": fmt.Errorf("read field: int128 support not implemented")},
		{"typeName": "uint128", "value": uint64(1), "encode": uint64(1), "expectedError": fmt.Errorf("read field: uint128 support not implemented")},
		{"typeName": "varint32", "value": int64(1), "encode": Varuint32(1), "expectedError": nil},
		{"typeName": "varuint32", "value": uint64(1), "encode": Varuint32(1), "expectedError": nil},
		{"typeName": "float32", "value": float32(1), "encode": float32(1), "expectedError": nil},
		{"typeName": "float64", "value": float64(1), "encode": float64(1), "expectedError": nil},
		{"typeName": "float128", "value": uint64(1), "encode": uint64(1), "expectedError": fmt.Errorf("read field: float128 support not implemented")},
		{"typeName": "bool", "value": true, "encode": true, "expectedError": nil},
		{"typeName": "bool", "value": false, "encode": false, "expectedError": nil},
		{"typeName": "time_point", "value": TimePoint(1), "encode": TimePoint(1), "expectedError": nil},
		{"typeName": "time_point_sec", "value": TimePointSec(1), "encode": TimePointSec(1), "expectedError": nil},
		{"typeName": "block_timestamp_type", "value": bt, "encode": bt, "expectedError": nil},
		{"typeName": "name", "value": Name("eoscanadacom"), "encode": Name("eoscanadacom"), "expectedError": nil},
		{"typeName": "bytes", "value": []byte("this.is.a.test"), "encode": []byte("this.is.a.test"), "expectedError": nil},
		{"typeName": "string", "value": "this.is.a.test", "encode": "this.is.a.test", "expectedError": nil},
		{"typeName": "checksum160", "value": Checksum160(make([]byte, TypeSize.Checksum160)), "encode": Checksum160(make([]byte, TypeSize.Checksum160)), "expectedError": nil},
		{"typeName": "checksum256", "value": Checksum256(make([]byte, TypeSize.Checksum256)), "encode": Checksum256(make([]byte, TypeSize.Checksum256)), "expectedError": nil},
		{"typeName": "checksum512", "value": Checksum512(make([]byte, TypeSize.Checksum512)), "encode": Checksum512(make([]byte, TypeSize.Checksum512)), "expectedError": nil},
		{"typeName": "public_key", "value": ecc.PublicKey{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 33)}, "encode": ecc.PublicKey{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 33)}, "expectedError": nil},
		{"typeName": "signature", "value": ecc.Signature{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 65)}, "encode": ecc.Signature{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 65)}, "expectedError": nil},
		{"typeName": "symbol", "value": &Symbol{Precision: 8, Symbol: "symbol.1"}, "encode": Symbol{Precision: 8, Symbol: "symbol.1"}, "expectedError": nil},
		{"typeName": "symbol_code", "value": SymbolCode(0), "encode": SymbolCode(0), "expectedError": nil},
		{"typeName": "asset", "value": Asset{Amount: 10, Symbol: EOSSymbol}, "encode": Asset{Amount: 10, Symbol: EOSSymbol}, "expectedError": nil},
		{"typeName": "extended_asset", "value": ExtendedAsset{Asset: Asset{Amount: 10, Symbol: EOSSymbol}, Contract: "eoscanadacom"}, "encode": ExtendedAsset{Asset: Asset{Amount: 10, Symbol: EOSSymbol}, Contract: "eoscanadacom"}, "expectedError": nil},
		{"typeName": "bad.type.1", "value": nil, "encode": nil, "expectedError": fmt.Errorf("read field [testedField] of type [bad.type.1]: unknown type")},
	}

	for _, c := range testCases {

		var b bytes.Buffer
		encoder := NewEncoder(&b)
		err := encoder.Encode(c["encode"])

		assert.NoError(t, err, fmt.Sprintf("encoding value %s, of type %s", c["value"], c["typeName"]))

		decoder := NewABIDecoder(b.Bytes(), nil)
		result := make(Result)
		err = decoder.read("testedField", c["typeName"].(string), result)

		assert.Equal(t, err, c["expectedError"])

		if c["expectedError"] == nil {
			assert.Equal(t, c["value"], result["testedField"])
		}

	}

}
