package eos

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/eoscanada/eos-go/ecc"

	"github.com/tidwall/gjson"

	"bytes"

	"strings"

	"github.com/stretchr/testify/assert"
)

func TestABI_DecodeAction(t *testing.T) {

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1    string
		F1     Name
		F2     string
		F3FLAG byte //this a hack until we have the abi encoder
		F3     string
		F4FLAG byte //this a hack until we have the abi encoder
		F5     []string
	}{
		BF1:    "value_struct_2_field_1",
		F1:     Name("eoscanadacom"),
		F2:     "value_struct_3_field_1",
		F3FLAG: 1,
		F3:     "value_struct_1_field_3",
		F4FLAG: 0,
		F5:     []string{"value_struct_4_field_1_1", "value_struct_4_field_1_2", "value_struct_4_field_1_3"},
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	assert.NoError(t, err)

	abi, err := NewABI(abiReader)
	assert.NoError(t, err)

	json, err := abi.DecodeAction(buffer.Bytes(), "action_name_1")
	assert.NoError(t, err)

	assert.Equal(t, "eoscanadacom", gjson.GetBytes(json, "struct_1_field_1").String())
	assert.Equal(t, "value_struct_2_field_1", gjson.GetBytes(json, "struct_2_field_1").String())
	assert.Equal(t, "value_struct_3_field_1", gjson.GetBytes(json, "struct_1_field_2.struct_3_field_1").String())
	assert.Equal(t, "value_struct_1_field_3", gjson.GetBytes(json, "struct_1_field_3").String())
	assert.Equal(t, "", gjson.GetBytes(json, "struct_1_field_4").String())
	assert.Equal(t, "value_struct_4_field_1_1", gjson.GetBytes(json, "struct_1_field_5.0.struct_4_field_1").String())
	assert.Equal(t, "value_struct_4_field_1_2", gjson.GetBytes(json, "struct_1_field_5.1.struct_4_field_1").String())
	assert.Equal(t, "value_struct_4_field_1_3", gjson.GetBytes(json, "struct_1_field_5.2.struct_4_field_1").String())

}

func TestABI_DecodeMissingData(t *testing.T) {
	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1 string
		F1  Name
	}{
		BF1: "value_struct_2_field_1",
		F1:  Name("eoscanadacom"),
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	assert.NoError(t, err)

	abi, err := NewABI(abiReader)
	assert.NoError(t, err)

	_, err = abi.DecodeAction(buffer.Bytes(), "action_name_1")
	assert.Equal(t, fmt.Errorf("decoding fields: decoding field [struct_1_field_2] of type [struct_name_3]: decoding fields: decoding field [struct_3_field_1] of type [string]: read: varint: invalid buffer size"), err)

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

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	assert.NoError(t, err)

	abi, err := NewABI(abiReader)
	assert.NoError(t, err)

	_, err = abi.DecodeAction(buffer.Bytes(), "bad.action.name")
	assert.Equal(t, fmt.Errorf("action bad.action.name not found in abi"), err)
}

func TestABI_DecodeTable(t *testing.T) {

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1    string
		F1     Name
		F2     string
		F3FLAG byte //this a hack until we have the abi encoder
		F3     string
		F4FLAG byte //this a hack until we have the abi encoder
		F5     []string
	}{
		BF1:    "value_struct_2_field_1",
		F1:     Name("eoscanadacom"),
		F2:     "value_struct_3_field_1",
		F3FLAG: 1,
		F3:     "value_struct_1_field_3",
		F4FLAG: 0,
		F5:     []string{"value_struct_4_field_1_1", "value_struct_4_field_1_2", "value_struct_4_field_1_3"},
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	assert.NoError(t, err)

	abi, err := NewABI(abiReader)
	assert.NoError(t, err)

	json, err := abi.DecodeTableRow("table_name_1", buffer.Bytes())
	assert.NoError(t, err)

	assert.Equal(t, "eoscanadacom", gjson.GetBytes(json, "struct_1_field_1").String())
	assert.Equal(t, "value_struct_2_field_1", gjson.GetBytes(json, "struct_2_field_1").String())
	assert.Equal(t, "value_struct_3_field_1", gjson.GetBytes(json, "struct_1_field_2.struct_3_field_1").String())
	assert.Equal(t, "value_struct_1_field_3", gjson.GetBytes(json, "struct_1_field_3").String())
	assert.Equal(t, "", gjson.GetBytes(json, "struct_1_field_4").String())
	assert.Equal(t, "value_struct_4_field_1_1", gjson.GetBytes(json, "struct_1_field_5.0.struct_4_field_1").String())
	assert.Equal(t, "value_struct_4_field_1_2", gjson.GetBytes(json, "struct_1_field_5.1.struct_4_field_1").String())
	assert.Equal(t, "value_struct_4_field_1_3", gjson.GetBytes(json, "struct_1_field_5.2.struct_4_field_1").String())

}

func TestABI_DecodeTableRowMissingTable(t *testing.T) {

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1 string
		F1  Name
	}{
		BF1: "value.base.field.1",
		F1:  Name("eoscanadacom"),
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	assert.NoError(t, err)

	abi, err := NewABI(abiReader)
	assert.NoError(t, err)

	_, err = abi.DecodeTableRow("bad.action.name", buffer.Bytes())
	assert.Equal(t, fmt.Errorf("table name bad.action.name not found in abi"), err)
}

func TestABI_DecodeBadABI(t *testing.T) {

	abiReader := strings.NewReader("{")
	_, err := NewABI(abiReader)
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

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(s)
	assert.NoError(t, err)

	json, err := abi.decode(NewDecoder(buffer.Bytes()), "struct.1")
	assert.NoError(t, err)

	assert.Equal(t, "value.field.1", gjson.GetBytes(json, "field.1").String())
	assert.Equal(t, "value.base.field.1", gjson.GetBytes(json, "base.field.1").String())

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

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(s)
	assert.NoError(t, err)

	_, err = abi.decode(NewDecoder(buffer.Bytes()), "struct.1")
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

	_, err = abi.decode(NewDecoder(b.Bytes()), "struct.1")
	assert.Equal(t, fmt.Errorf("structure [struct.1] not found in abi"), err)
}

func TestABI_decodeFields(t *testing.T) {

	types := []ABIType{
		{NewTypeName: "action.type.1", Type: "name"},
	}
	fields := []FieldDef{
		{Name: "F1", Type: "uint64"},
		{Name: "F2", Type: "action.type.1"},
	}
	abi := &ABI{
		Types: types,
		Structs: []StructDef{
			{Fields: fields},
		},
	}

	s := struct {
		F1 uint64
		F2 Name
	}{
		F1: uint64(18446744073709551615),
		F2: Name("eoscanadacom"),
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(s)
	assert.NoError(t, err)

	json, err := abi.decodeFields(NewDecoder(buffer.Bytes()), fields, []byte{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(18446744073709551615), gjson.GetBytes(json, "F1").Uint())
	assert.Equal(t, "eoscanadacom", gjson.GetBytes(json, "F2").String())
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

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(s)
	assert.NoError(t, err)

	_, err = abi.decodeFields(NewDecoder(buffer.Bytes()), fields, []byte{})
	assert.Equal(t, fmt.Errorf("decoding fields: decoding field [field.with.bad.type.1] of type [bad.type.1]: read field of type [bad.type.1]: unknown type"), err)

}

func TestABI_Read(t *testing.T) {
	someTime, err := time.Parse("2006-01-02T15:04:05", "2018-09-05T12:48:54")
	assert.NoError(t, err)
	bt := BlockTimestamp{
		Time: someTime,
	}
	assert.NoError(t, err)

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
		{"caseName": "string", "typeName": "string", "value": "\"this.is.a.test\"", "encode": "this.is.a.test", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min int8", "typeName": "int8", "value": "-128", "encode": int8(-128), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int8", "typeName": "int8", "value": "127", "encode": int8(127), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min uint8", "typeName": "uint8", "value": "0", "encode": uint8(0), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint8", "typeName": "uint8", "value": "255", "encode": uint8(255), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min int16", "typeName": "int16", "value": "-32768", "encode": int16(-32768), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int16", "typeName": "int16", "value": "32767", "encode": int16(32767), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min uint16", "typeName": "uint16", "value": "0", "encode": uint16(0), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint16", "typeName": "uint16", "value": "65535", "encode": uint16(65535), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min int32", "typeName": "int32", "value": "-2147483648", "encode": int32(-2147483648), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int32", "typeName": "int32", "value": "2147483647", "encode": int32(2147483647), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min uint32", "typeName": "uint32", "value": "0", "encode": uint32(0), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint32", "typeName": "uint32", "value": "4294967295", "encode": uint32(4294967295), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min int64", "typeName": "int64", "value": `"-9223372036854775808"`, "encode": int64(-9223372036854775808), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int64", "typeName": "int64", "value": `"9223372036854775807"`, "encode": int64(9223372036854775807), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "mid int64", "typeName": "int64", "value": `4096`, "encode": int64(4096), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "stringified lower int64", "typeName": "int64", "value": `"-5000000000"`, "encode": int64(-5000000000), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min uint64", "typeName": "uint64", "value": "0", "encode": uint64(0), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint64", "typeName": "uint64", "value": `"18446744073709551615"`, "encode": uint64(18446744073709551615), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "int128", "typeName": "int128", "value": `"0x01000000000000000200000000000000"`, "encode": Int128{Lo: 1, Hi: 2}, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "uint128", "typeName": "uint128", "value": `"0x01000000000000000200000000000000"`, "encode": Uint128{Lo: 1, Hi: 2}, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min varint32", "typeName": "varint32", "value": "-2147483648", "encode": Varint32(-2147483648), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max varint32", "typeName": "varint32", "value": "2147483647", "encode": Varint32(2147483647), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min varuint32", "typeName": "varuint32", "value": "0", "encode": Varuint32(0), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max varuint32", "typeName": "varuint32", "value": "4294967295", "encode": Varuint32(4294967295), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min float 32", "typeName": "float32", "value": "0.000000000000000000000000000000000000000000001401298464324817", "encode": float32(math.SmallestNonzeroFloat32), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max float 32", "typeName": "float32", "value": "340282346638528860000000000000000000000", "encode": float32(math.MaxFloat32), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min float64", "typeName": "float64", "value": "0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005", "encode": math.SmallestNonzeroFloat64, "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max float64", "typeName": "float64", "value": "179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "encode": math.MaxFloat64, "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "float128", "typeName": "float128", "value": `"0x01000000000000000200000000000000"`, "encode": Float128{Lo: 1, Hi: 2}, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "bool true", "typeName": "bool", "value": "true", "encode": true, "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "bool false", "typeName": "bool", "value": "false", "encode": false, "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "time_point", "typeName": "time_point", "value": "\"2018-11-01T15:13:07.001\"", "encode": TimePoint(1541085187001001), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "time_point_sec", "typeName": "time_point_sec", "value": "\"2023-04-14T10:55:53\"", "encode": TimePointSec(1681469753), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "block_timestamp_type", "typeName": "block_timestamp_type", "value": "\"2018-09-05T12:48:54\"", "encode": bt, "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "Name", "typeName": "name", "value": "\"eoscanadacom\"", "encode": Name("eoscanadacom"), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "bytes", "typeName": "bytes", "value": "\"746869732e69732e612e74657374\"", "encode": []byte("this.is.a.test"), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "checksum160", "typeName": "checksum160", "value": "\"0000000000000000000000000000000000000000\"", "encode": Checksum160(make([]byte, TypeSize.Checksum160)), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "checksum256", "typeName": "checksum256", "value": "\"0000000000000000000000000000000000000000000000000000000000000000\"", "encode": Checksum256(make([]byte, TypeSize.Checksum256)), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "checksum512", "typeName": "checksum512", "value": "\"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\"", "encode": Checksum512(make([]byte, TypeSize.Checksum512)), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "public_key", "typeName": "public_key", "value": "\"EOS1111111111111111111111111111111114T1Anm\"", "encode": ecc.MustNewPublicKey("EOS1111111111111111111111111111111114T1Anm"), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "signature", "typeName": "signature", "value": "\"SIG_K1_111111111111111111111111111111111111111111111111111111111111111116uk5ne\"", "encode": ecc.Signature{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 65)}, "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "symbol", "typeName": "symbol", "value": "\"4,EOS\"", "encode": EOSSymbol, "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "symbol_code", "typeName": "symbol_code", "value": "18446744073709551615", "encode": SymbolCode(18446744073709551615), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "asset", "typeName": "asset", "value": "\"10.0000 EOS\"", "encode": Asset{Amount: 100000, Symbol: EOSSymbol}, "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "extended_asset", "typeName": "extended_asset", "value": "{\"asset\":\"0.0010 EOS\",\"Contract\":\"eoscanadacom\"}", "encode": ExtendedAsset{Asset: Asset{Amount: 10, Symbol: EOSSymbol}, Contract: "eoscanadacom"}, "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "bad type", "typeName": "bad.type.1", "value": nil, "encode": nil, "expectedError": fmt.Errorf("decoding field [testedField] of type [bad.type.1]: read field of type [bad.type.1]: unknown type"), "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "optional present", "typeName": "string", "value": "\"value.1\"", "encode": optional, "expectedError": nil, "isOptional": true, "isArray": false, "fieldName": "testedField"},
		{"caseName": "optional not present", "typeName": "string", "value": "", "encode": optionalNotPresent, "expectedError": nil, "isOptional": true, "isArray": false, "fieldName": "testedField"},
		{"caseName": "optional missing flag", "typeName": "string", "value": nil, "encode": optionalMissingFlag, "expectedError": fmt.Errorf("decoding field [testedField] optional flag: byte required [1] byte, remaining [0]"), "isOptional": true, "isArray": false, "fieldName": "testedField"},
		{"caseName": "array", "typeName": "string", "value": "[\"value.1\",\"value.2\"]", "encode": []string{"value.1", "value.2"}, "expectedError": nil, "isOptional": false, "isArray": true, "fieldName": "testedField"},
		{"caseName": "array empty", "typeName": "string", "value": "[]", "encode": []string{}, "expectedError": nil, "isOptional": false, "isArray": true, "fieldName": "testedField"},
		{"caseName": "missing array", "typeName": "string", "value": nil, "encode": nil, "expectedError": fmt.Errorf("reading field [testedField] array length: varint: invalid buffer size"), "isOptional": false, "isArray": true, "fieldName": "testedField"},
		{"caseName": "array item unknown type", "typeName": "invalid.field.type", "value": nil, "encode": []string{"value.1", "value.2"}, "expectedError": fmt.Errorf("reading field [testedField] index [0]: read field of type [invalid.field.type]: unknown type"), "isOptional": false, "isArray": true, "fieldName": "testedField"},
	}

	for _, c := range testCases {

		t.Run(c["caseName"].(string), func(t *testing.T) {
			var buffer bytes.Buffer
			encoder := NewEncoder(&buffer)
			err := encoder.Encode(c["encode"])

			assert.NoError(t, err, fmt.Sprintf("encoding value %s, of type %s", c["value"], c["typeName"]), c["caseName"])

			abi := ABI{}
			json, err := abi.decodeField(NewDecoder(buffer.Bytes()), c["fieldName"].(string), c["typeName"].(string), c["isOptional"].(bool), c["isArray"].(bool), []byte{})

			//fmt.Println("JSON:", string(json))
			assert.Equal(t, c["expectedError"], err, c["caseName"])

			if c["expectedError"] == nil {
				assert.Equal(t, c["value"], gjson.GetBytes(json, c["fieldName"].(string)).Raw, c["caseName"])
			}
		})
	}
}

func TestABI_Read_TimePointSec(t *testing.T) {
	// EnableDecoderLogging()
	// EnableABIDecoderLogging()

	abi := ABI{}
	data, err := hex.DecodeString("919dd85b")
	require.NoError(t, err)
	out, err := abi.decodeField(NewDecoder(data), "name", "time_point_sec", false, false, []byte("{}"))
	//out, err := abi.decodeField(NewDecoder([]byte("c15dd35b")), "name", "time_point_sec", false, false, []byte("{}"))
	//out, err := abi.decodeField(NewDecoder([]byte("919dd85b")), "name", "time_point_sec", false, false, []byte("{}"))
	require.NoError(t, err)
	assert.Equal(t, `{"name":"2018-10-30T18:06:09"}`, string(out))
}

func TestABIDecoder_analyseFieldType(t *testing.T) {

	testCases := []map[string]interface{}{
		{"fieldType": "field.type.1", "expectedName": "field.type.1", "expectedOptional": false, "expectedArray": false},
		{"fieldType": "field.type.1?", "expectedName": "field.type.1", "expectedOptional": true, "expectedArray": false},
		{"fieldType": "field.type.1[]", "expectedName": "field.type.1", "expectedOptional": false, "expectedArray": true},
	}

	for _, c := range testCases {
		name, isOption, isArray := analyzeFieldType(c["fieldType"].(string))
		assert.Equal(t, c["expectedName"], name)
		assert.Equal(t, c["expectedOptional"], isOption)
		assert.Equal(t, c["expectedArray"], isArray)
	}
}
