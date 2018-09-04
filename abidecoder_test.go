package eos

import (
	"os"
	"testing"

	"fmt"

	"bytes"

	"time"

	"strings"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/stretchr/testify/assert"
)

func TestABI_Decode(t *testing.T) {

	Logger.Decoder.SetOutput(os.Stdout)

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1    string
		F1     Name
		F2     string
		F3FLAG byte //this a hack until we have the abi encoder
		F3     string
		F4FLAG byte //this a hack until we have the abi encoder
	}{
		BF1:    "value.struct.2.field.1",
		F1:     Name("eoscanadacom"),
		F2:     "value.struct.3.field.1",
		F3FLAG: 1,
		F3:     "value.struct.1.field.3",
		F4FLAG: 0,
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
	assert.Equal(t, "value.struct.1.field.3", result["struct.1.field.3"])

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
	assert.Equal(t, fmt.Errorf("decoding fields: decoding field [struct.3.field.1] of type [string]: read value: varint: invalid buffer size"), err)

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
	result := make(ABIMap)
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
	result := make(ABIMap)
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
	result := make(ABIMap)
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
	result := make(ABIMap)
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
	result := make(ABIMap)

	err = decoder.decodeFields(fields, result)
	assert.Equal(t, fmt.Errorf("decoding fields: decoding field [field.with.bad.type.1] of type [bad.type.1]: read value: read field of type [bad.type.1]: unknown type"), err)

}

func TestABI_Read(t *testing.T) {
	bt := BlockTimestamp{
		Time: time.Unix(time.Now().Unix(), 0),
	}

	optional := struct {
		B byte
		S string
	}{
		B: 1,
		S: "value.1",
	}
	optionalNotPresent := struct {
		B byte
		S string
	}{
		B: 0,
	}
	optionalMissingFlag := struct {
	}{}

	testCases := []map[string]interface{}{
		{"caseName": "case.1", "typeName": "int8", "value": int8(1), "encode": int8(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.2", "typeName": "uint8", "value": uint8(1), "encode": uint8(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.3", "typeName": "int16", "value": int16(1), "encode": int16(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.4", "typeName": "uint16", "value": uint16(1), "encode": uint16(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.5", "typeName": "int32", "value": int32(1), "encode": int32(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.6", "typeName": "uint32", "value": uint32(1), "encode": uint32(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.7", "typeName": "int64", "value": int64(1), "encode": int64(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.8", "typeName": "uint64", "value": uint64(1), "encode": uint64(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.9", "typeName": "int128", "value": int64(1), "encode": int64(1), "expectedError": fmt.Errorf("decoding field [testedField] of type [int128]: read value: read field: int128 support not implemented"), "isOptional": false, "isArray": false},
		{"caseName": "case.10", "typeName": "uint128", "value": uint64(1), "encode": uint64(1), "expectedError": fmt.Errorf("decoding field [testedField] of type [uint128]: read value: read field: uint128 support not implemented"), "isOptional": false, "isArray": false},
		{"caseName": "case.11", "typeName": "varint32", "value": int64(1), "encode": Varuint32(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.12", "typeName": "varuint32", "value": uint64(1), "encode": Varuint32(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.13", "typeName": "float32", "value": float32(1), "encode": float32(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.14", "typeName": "float64", "value": float64(1), "encode": float64(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.15", "typeName": "float128", "value": uint64(1), "encode": uint64(1), "expectedError": fmt.Errorf("decoding field [testedField] of type [float128]: read value: read field: float128 support not implemented"), "isOptional": false, "isArray": false},
		{"caseName": "case.16", "typeName": "bool", "value": true, "encode": true, "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.17", "typeName": "bool", "value": false, "encode": false, "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.18", "typeName": "time_point", "value": TimePoint(1), "encode": TimePoint(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.19", "typeName": "time_point_sec", "value": TimePointSec(1), "encode": TimePointSec(1), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.20", "typeName": "block_timestamp_type", "value": bt, "encode": bt, "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.21", "typeName": "name", "value": Name("eoscanadacom"), "encode": Name("eoscanadacom"), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.22", "typeName": "bytes", "value": []byte("this.is.a.test"), "encode": []byte("this.is.a.test"), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.23", "typeName": "string", "value": "this.is.a.test", "encode": "this.is.a.test", "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.24", "typeName": "checksum160", "value": Checksum160(make([]byte, TypeSize.Checksum160)), "encode": Checksum160(make([]byte, TypeSize.Checksum160)), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.25", "typeName": "checksum256", "value": Checksum256(make([]byte, TypeSize.Checksum256)), "encode": Checksum256(make([]byte, TypeSize.Checksum256)), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.26", "typeName": "checksum512", "value": Checksum512(make([]byte, TypeSize.Checksum512)), "encode": Checksum512(make([]byte, TypeSize.Checksum512)), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.27", "typeName": "public_key", "value": ecc.PublicKey{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 33)}, "encode": ecc.PublicKey{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 33)}, "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.28", "typeName": "signature", "value": ecc.Signature{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 65)}, "encode": ecc.Signature{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 65)}, "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.29", "typeName": "symbol", "value": &Symbol{Precision: 8, Symbol: "symbol.1"}, "encode": Symbol{Precision: 8, Symbol: "symbol.1"}, "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.30", "typeName": "symbol_code", "value": SymbolCode(0), "encode": SymbolCode(0), "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.31", "typeName": "asset", "value": Asset{Amount: 10, Symbol: EOSSymbol}, "encode": Asset{Amount: 10, Symbol: EOSSymbol}, "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.32", "typeName": "extended_asset", "value": ExtendedAsset{Asset: Asset{Amount: 10, Symbol: EOSSymbol}, Contract: "eoscanadacom"}, "encode": ExtendedAsset{Asset: Asset{Amount: 10, Symbol: EOSSymbol}, Contract: "eoscanadacom"}, "expectedError": nil, "isOptional": false, "isArray": false},
		{"caseName": "case.33", "typeName": "bad.type.1", "value": nil, "encode": nil, "expectedError": fmt.Errorf("decoding field [testedField] of type [bad.type.1]: read value: read field of type [bad.type.1]: unknown type"), "isOptional": false, "isArray": false},
		{"caseName": "case.34", "typeName": "string", "value": "value.1", "encode": optional, "expectedError": nil, "isOptional": true, "isArray": false},
		{"caseName": "case.35", "typeName": "string", "value": nil, "encode": optionalNotPresent, "expectedError": nil, "isOptional": true, "isArray": false},
		{"caseName": "case.36", "typeName": "string", "value": nil, "encode": optionalNotPresent, "expectedError": nil, "isOptional": true, "isArray": false},
		{"caseName": "case.37", "typeName": "string", "value": nil, "encode": optionalMissingFlag, "expectedError": fmt.Errorf("decoding field [testedField] optional flag: byte required [1] byte, remaining [0]"), "isOptional": true, "isArray": false},
		{"caseName": "case.38", "typeName": "string", "value": []interface{}{"value.1", "value.2"}, "encode": []string{"value.1", "value.2"}, "expectedError": nil, "isOptional": false, "isArray": true},
		{"caseName": "case.39", "typeName": "string", "value": nil, "encode": nil, "expectedError": fmt.Errorf("reading field [testedField] array length: varint: invalid buffer size"), "isOptional": false, "isArray": true},
		{"caseName": "case.40", "typeName": "invalid.field.type", "value": nil, "encode": []string{"value.1", "value.2"}, "expectedError": fmt.Errorf("reading field [testedField] index [0]: read field of type [invalid.field.type]: unknown type"), "isOptional": false, "isArray": true},
	}

	for _, c := range testCases {

		t.Run(c["caseName"].(string), func(t *testing.T) {
			var b bytes.Buffer
			encoder := NewEncoder(&b)
			err := encoder.Encode(c["encode"])

			assert.NoError(t, err, fmt.Sprintf("encoding value %s, of type %s", c["value"], c["typeName"]), c["caseName"])

			decoder := NewABIDecoder(b.Bytes(), nil)
			result := make(ABIMap)
			err = decoder.decodeField("testedField", c["typeName"].(string), c["isOptional"].(bool), c["isArray"].(bool), result)

			assert.Equal(t, c["expectedError"], err, c["caseName"])

			if c["expectedError"] == nil {
				assert.Equal(t, c["value"], result["testedField"], c["caseName"])
			}
		})
	}
}

func TestABIDecoder_analyseFieldType(t *testing.T) {

	testCases := []map[string]interface{}{
		{"fieldName": "field.name.1", "expectedName": "field.name.1", "expectedOptional": false, "expectedArray": false},
		{"fieldName": "field.name.1?", "expectedName": "field.name.1", "expectedOptional": true, "expectedArray": false},
		{"fieldName": "field.name.1[]", "expectedName": "field.name.1", "expectedOptional": false, "expectedArray": true},
	}

	for _, c := range testCases {
		name, isOption, isArray := analyseFieldName(c["fieldName"].(string))
		assert.Equal(t, c["expectedName"], name)
		assert.Equal(t, c["expectedOptional"], isOption)
		assert.Equal(t, c["expectedArray"], isArray)
	}
}

//{
//"expiration": "2018-08-29T20:52:30",
//"ref_block_num": 11531,
//"ref_block_prefix": 740532780,
//"max_net_usage_words": 0,
//"max_cpu_usage_ms": 0,
//"delay_sec": 0,
//"context_free_actions": [],
//"actions": [
//{
//"account": "eosio",
//"name": "updateauth",
//"authorization": [
//{
//"actor": "eoscanadacom",
//"permission": "owner"
//}
//],
//"data": "202932c94c8330550000000080ab26a70000000000000000050000000006608c31c94c83305500000000a8ed32320200708c31c94c83305500000000a8ed32320200808c31c94c83305500000000a8ed32320200908c31c94c83305500000000a8ed32320200a08c31c94c83305500000000a8ed32320200b08c31c94c83305500000000a8ed3232010002805101000100803a09000200"
//}
//],
//"transaction_extensions": [],
//"signatures": [],
//"context_free_data": []
//}
