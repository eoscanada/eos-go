package eos

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

//
//import (
//	"fmt"
//	"testing"
//
//	"bytes"
//
//	"strings"
//
//	"github.com/stretchr/testify/assert"
//)
//
var abiString = `
{
	"version": "eosio::abi/1.0",
	"types": [{
		"new_type_name": "new_type_name_1",
		"type": "name"
	}],
	"structs": [
	{
		"name": "struct_name_1",
		"base": "struct_name_2",
		"fields": [
			{"name":"struct_1_field_1", "type":"new_type_name_1"},
			{"name":"struct_1_field_2", "type":"struct_name_3"},
			{"name":"struct_1_field_3?", "type":"string"},
			{"name":"struct_1_field_4?", "type":"string"},
			{"name":"struct_1_field_5[]", "type":"struct_name_4"}
		]
   },{
		"name": "struct_name_2",
		"base": "",
		"fields": [
			{"name":"struct_2_field_1", "type":"string"}
		]
   },{
		"name": "struct_name_3",
		"base": "",
		"fields": [
			{"name":"struct_3_field_1", "type":"string"}
		]
   },{
		"name": "struct_name_4",
		"base": "",
		"fields": [
			{"name":"struct_4_field_1", "type":"string"}
		]
   }
	],
  "actions": [{
		"name": "action_name_1",
		"type": "struct_name_1",
		"ricardian_contract": ""
  }]
}
`

//var abiData = M{
//	"struct_2_field.1": "struct.2.field.1.value",
//	"struct.1.field.1": Name("eoscanadacom"),
//	"struct.1.field.2": M{
//		"struct.3.field.1": "struct.3.field.1.value",
//	},
//	"struct.1.field.3": "struct.1.field.3.value",
//	//"struct.1.field.4": "struct.1.field.4.value",
//}

//
//func TestABIEncoder_Encode(t *testing.T) {
//
//	testCases := []map[string]interface{}{
//		{"caseName": "sunny path", "actionName": "action.name.1", "expectedError": nil, "writer": new(bytes.Buffer), "abi": abiString},
//		{"caseName": "missing action", "actionName": "action.name.missing", "expectedError": fmt.Errorf("action action.name.missing not found in abi"), "writer": new(bytes.Buffer), "abi": abiString},
//		{"caseName": "abi reader error", "actionName": "action.name.missing", "expectedError": fmt.Errorf("encode: read abi: unexpected EOF"), "writer": new(bytes.Buffer), "abi": "{"},
//	}
//
//	for _, c := range testCases {
//		caseName := c["caseName"].(string)
//		t.Run(caseName, func(t *testing.T) {
//
//			buf := c["writer"].(mockWriterable)
//			encoder := NewABIEncoder(strings.NewReader(c["abi"].(string)), buf)
//			err := encoder.Encode(ActionName(c["actionName"].(string)), abiData)
//			assert.Equal(t, c["expectedError"], err)
//
//			if c["expectedError"] != nil {
//				return
//			}
//
//			decoder := NewABIDecoder(buf.Bytes(), strings.NewReader(abiString))
//			result := make(M)
//			err = decoder.Decode(result, ActionName(c["actionName"].(string)))
//			assert.NoError(t, err)
//
//			assert.Equal(t, abiData, result)
//			//fmt.Println(result)
//		})
//	}
//}
//func TestABIEncoder_encodeMissingActionStruct(t *testing.T) {
//
//	abiString := `
//{
//	"version": "eosio::abi/1.0",
//	"types": [{
//		"new_type_name": "new.type.name.1",
//		"type": "name"
//	}],
//	"structs": [
//	],
//   "actions": [{
//		"name": "action.name.1",
//		"type": "struct.name.1",
//		"ricardian_contract": ""
//   }]
//}
//`
//
//	buf := new(bytes.Buffer)
//	encoder := NewABIEncoder(strings.NewReader(abiString), buf)
//
//	err := encoder.Encode("action.name.1", map[string]interface{}{})
//	assert.Equal(t, fmt.Errorf("encode: structure [struct.name.1] not found in abi"), err)
//}
//
//func TestABIEncoder_encodeErrorInBase(t *testing.T) {
//
//	abiString := `
//{
//	"version": "eosio::abi/1.0",
//	"types": [{
//		"new_type_name": "new.type.name.1",
//		"type": "name"
//	}],
//	"structs": [
//	{
//		"name": "struct.name.1",
//		"base": "struct.name.2",
//		"fields": [
//			{"name":"struct.1.field.1", "type":"new.type.name.1"}
//		]
//    }
//	],
//   "actions": [{
//		"name": "action.name.1",
//		"type": "struct.name.1",
//		"ricardian_contract": ""
//   }]
//}
//`
//
//	buf := new(bytes.Buffer)
//	encoder := NewABIEncoder(strings.NewReader(abiString), buf)
//
//	err := encoder.Encode("action.name.1", map[string]interface{}{})
//	assert.Equal(t, fmt.Errorf("encode base [struct.name.1]: encode: structure [struct.name.2] not found in abi"), err)
//}
//
//func TestABI_EncodeEncodeName(t *testing.T) {
//
//	abiString := `
//{
//	"version": "eosio::abi/1.0",
//	"types": [{
//		"new_type_name": "new.type.name.1",
//		"type": "name"
//	}],
//	"structs": [
//	{
//		"name": "struct.name.1",
//		"base": "",
//		"fields": [
//			{"name":"struct.1.field.1", "type":"new.type.name.1"}
//		]
//    }
//	],
//   "actions": [{
//		"name": "action.name.1",
//		"type": "struct.name.1",
//		"ricardian_contract": ""
//   }]
//}
//`
//	abiData := map[string]interface{}{
//		"struct.1.field.1": Name("struct.1.field.1.value"),
//	}
//	buf := new(bytes.Buffer)
//	encoder := NewABIEncoder(strings.NewReader(abiString), buf)
//	err := encoder.Encode(ActionName("action.name.1"), abiData)
//	assert.Equal(t, nil, err)
//
//	decoder := NewABIDecoder(buf.Bytes(), strings.NewReader(abiString))
//	result := make(ABIMap)
//	err = decoder.Decode(result, ActionName("action.name.1"))
//	assert.NoError(t, err)
//	fmt.Println(result)
//
//}
//
//func TestABIEncoder_encodeFields(t *testing.T) {
//
//	abiString = `
//{
//	"version": "eosio::abi/1.0",
//	"types": [
//		{
//			"new_type_name": "new.type.1",
//			"type": "string"
//		}
//	],
//	"structs": [
//		{
//			"name": "struct.name.1",
//			"base": "",
//			"fields": [
//				{"name":"struct.1.field.1", "type":"new.type.name.1"}
//			]
//    	}
//	],
//   	"actions": []
//}
//
//`
//	testCases := []map[string]interface{}{
//		{"caseName": "sunny path", "fields": []FieldDef{{Name: "field.name.1", Type: "new.type.1"}}, "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "expectedError": nil, "writer": new(bytes.Buffer)},
//		{"caseName": "encodeField error", "fields": []FieldDef{{Name: "field.name.1", Type: "new.type.1"}}, "actionMap": map[string]interface{}{}, "expectedError": fmt.Errorf("encoding fields: encode field: none optional field [field.name.1] as a nil value"), "writer": new(bytes.Buffer)},
//		{"caseName": "embedded struct wrong type", "fields": []FieldDef{{Name: "field.name.1", Type: "struct.name.1"}}, "actionMap": map[string]interface{}{"field.name.1": map[string]interface{}{}}, "expectedError": fmt.Errorf("encode fields: structure field [field.name.1] expected to be of type ABIMap"), "writer": new(bytes.Buffer)},
//		{"caseName": "encodeField error embedded struct", "fields": []FieldDef{{Name: "field.name.1", Type: "struct.name.1"}}, "actionMap": map[string]interface{}{"field.name.1": make(ABIMap)}, "expectedError": fmt.Errorf("encoding fields: encode field: none optional field [struct.1.field.1] as a nil value"), "writer": new(bytes.Buffer)},
//	}
//
//	for _, c := range testCases {
//		caseName := c["caseName"].(string)
//		t.Run(caseName, func(t *testing.T) {
//
//			buf := c["writer"].(mockWriterable)
//			encoder := NewABIEncoder(nil, buf)
//
//			abi, err := NewABI(strings.NewReader(abiString))
//			assert.NoError(t, err)
//			encoder.abi = abi
//
//			err = encoder.encodeFields(c["fields"].([]FieldDef), c["actionMap"].(map[string]interface{}))
//			assert.Equal(t, c["expectedError"], err)
//		})
//	}
//}
//
//func TestABIEncoder_encodeField(t *testing.T) {
//
//	testCases := []map[string]interface{}{
//		{"caseName": "sunny path", "fieldName": "field.name.1", "fieldType": "string", "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "isOptional": false, "isArray": false, "expectedError": nil, "writer": new(bytes.Buffer)},
//		{"caseName": "optional with value", "fieldName": "field.name.1", "fieldType": "string", "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "isOptional": true, "isArray": false, "expectedError": nil, "writer": new(bytes.Buffer)},
//		{"caseName": "optional with no value", "fieldName": "field.name.1", "fieldType": "string", "actionMap": map[string]interface{}{}, "isOptional": true, "isArray": false, "expectedError": nil, "writer": new(bytes.Buffer)},
//		{"caseName": "none optional with nil value", "fieldName": "field.name.2", "fieldType": "string", "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "isOptional": false, "isArray": false, "expectedError": fmt.Errorf("encode field: none optional field [field.name.2] as a nil value"), "writer": new(bytes.Buffer)},
//		{"caseName": "optional write flag err", "fieldName": "field.name.1", "fieldType": "string", "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "isOptional": true, "isArray": false, "expectedError": fmt.Errorf("mocked error"), "writer": mockWriter{0, fmt.Errorf("mocked error")}},
//	}
//
//	for _, c := range testCases {
//		caseName := c["caseName"].(string)
//		t.Run(caseName, func(t *testing.T) {
//			buf := c["writer"].(mockWriterable)
//			encoder := NewABIEncoder(nil, buf)
//
//			action := c["actionMap"].(map[string]interface{})
//			fieldName := c["fieldName"].(string)
//			fieldType := c["fieldType"].(string)
//			isOptional := c["isOptional"].(bool)
//			isArray := c["isArray"].(bool)
//			expectedError := c["expectedError"]
//
//			err := encoder.encodeField(fieldName, isOptional, isArray, action)
//			assert.Equal(t, expectedError, err, caseName)
//			if expectedError != nil {
//				return
//			}
//
//			decoder := NewABIDecoder(buf.Bytes(), nil)
//			result := make(ABIMap)
//			err = decoder.decodeField(fieldName, fieldType, isOptional, isArray, result)
//			assert.NoError(t, err, caseName)
//
//			assert.Equal(t, action[fieldName], result[fieldName], caseName)
//
//		})
//
//	}
//}
//

func TestABI_Write(t *testing.T) {

	Logger.ABIEncoder.SetOutput(os.Stdout)
	testCases := []map[string]interface{}{
		{"caseName": "string", "typeName": "string", "expectedValue": "0e746869732e69732e612e74657374", "json": "{\"testField\":\"this.is.a.test\""},
		{"caseName": "min int8", "typeName": "int8", "expectedValue": "80", "json": "{\"testField\":-128}"},
		{"caseName": "max int8", "typeName": "int8", "expectedValue": "7f", "json": "{\"testField\":127}", "expectedError": nil},
		{"caseName": "out of range int8", "typeName": "int8", "expectedValue": "", "json": "{\"testField\":128}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int8 : strconv.ParseInt: parsing \"128\": value out of range")},
		{"caseName": "out of range int8", "typeName": "int8", "expectedValue": "", "json": "{\"testField\":-129}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int8 : strconv.ParseInt: parsing \"-129\": value out of range")},
		{"caseName": "min uint8", "typeName": "uint8", "expectedValue": "00", "json": "{\"testField\":0}", "expectedError": nil},
		{"caseName": "max uint8", "typeName": "uint8", "expectedValue": "ff", "json": "{\"testField\":255}", "expectedError": nil},
		{"caseName": "out of range uint8", "typeName": "uint8", "expectedValue": "", "json": "{\"testField\":-1}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint8 : strconv.ParseUint: parsing \"-1\": invalid syntax")},
		{"caseName": "out of range uint8", "typeName": "uint8", "expectedValue": "", "json": "{\"testField\":256}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint8 : strconv.ParseUint: parsing \"256\": value out of range")},
		{"caseName": "min int16", "typeName": "int16", "expectedValue": "0080", "json": "{\"testField\":-32768}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int16", "typeName": "int16", "expectedValue": "ff7f", "json": "{\"testField\":32767}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range int16", "typeName": "int16", "expectedValue": "", "json": "{\"testField\":-32769}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int16 : strconv.ParseInt: parsing \"-32769\": value out of range")},
		{"caseName": "out of range int16", "typeName": "int16", "expectedValue": "", "json": "{\"testField\":32768}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int16 : strconv.ParseInt: parsing \"32768\": value out of range")},
		{"caseName": "min uint16", "typeName": "uint16", "expectedValue": "0000", "json": "{\"testField\":0}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint16", "typeName": "uint16", "expectedValue": "ffff", "json": "{\"testField\":65535}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range uint16", "typeName": "uint16", "expectedValue": "", "json": "{\"testField\":-1}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint16 : strconv.ParseUint: parsing \"-1\": invalid syntax")},
		{"caseName": "out of range uint16", "typeName": "uint16", "expectedValue": "", "json": "{\"testField\":65536}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint16 : strconv.ParseUint: parsing \"65536\": value out of range")},
		{"caseName": "min int32", "typeName": "int32", "expectedValue": "00000080", "json": "{\"testField\":-2147483648}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int32", "typeName": "int32", "expectedValue": "ffffff7f", "json": "{\"testField\":2147483647}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range int32", "typeName": "int32", "expectedValue": "", "json": "{\"testField\":-2147483649}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int32 : strconv.ParseInt: parsing \"-2147483649\": value out of range")},
		{"caseName": "out of range int32", "typeName": "int32", "expectedValue": "", "json": "{\"testField\":2147483648}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int32 : strconv.ParseInt: parsing \"2147483648\": value out of range")},
		{"caseName": "min uint32", "typeName": "uint32", "expectedValue": "00000000", "json": "{\"testField\":0}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint32", "typeName": "uint32", "expectedValue": "ffffffff", "json": "{\"testField\":4294967295}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range uint32", "typeName": "uint32", "expectedValue": "", "json": "{\"testField\":-1}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint32 : strconv.ParseUint: parsing \"-1\": invalid syntax")},
		{"caseName": "out of range uint32", "typeName": "uint32", "expectedValue": "", "json": "{\"testField\":4294967296}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint32 : strconv.ParseUint: parsing \"4294967296\": value out of range")},
		{"caseName": "min int64", "typeName": "int64", "expectedValue": "0000000000000080", "json": "{\"testField\":-9223372036854775808}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int64", "typeName": "int64", "expectedValue": "ffffffffffffff7f", "json": "{\"testField\":9223372036854775807}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range int64", "typeName": "int64", "expectedValue": "", "json": "{\"testField\":-9223372036854775809}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int64 : strconv.ParseInt: parsing \"-9223372036854775809\": value out of range")},
		{"caseName": "out of range int64", "typeName": "int64", "expectedValue": "", "json": "{\"testField\":9223372036854775808}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int64 : strconv.ParseInt: parsing \"9223372036854775808\": value out of range")},
		{"caseName": "min uint64", "typeName": "uint64", "expectedValue": "0000000000000000", "json": "{\"testField\":0}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint64", "typeName": "uint64", "expectedValue": "ffffffffffffffff", "json": "{\"testField\":18446744073709551615}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range uint64", "typeName": "uint64", "expectedValue": "", "json": "{\"testField\":-1}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint64 : strconv.ParseUint: parsing \"-1\": invalid syntax")},
		{"caseName": "out of range uint64", "typeName": "uint64", "expectedValue": "", "json": "{\"testField\":18446744073709551616}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint64 : strconv.ParseUint: parsing \"18446744073709551616\": value out of range")},
		{"caseName": "uint128 unsupported", "typeName": "uint128", "expectedValue": "", "json": "{\"testField\":18446744073709551616}", "expectedError": fmt.Errorf("writing field: uint128 support not implemented")},
		{"caseName": "varint32", "typeName": "varint32", "expectedValue": "00000080", "json": "{\"testField\":-2147483648}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "varuint32", "typeName": "varuint32", "expectedValue": "ffffffff", "json": "{\"testField\":4294967295}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"}, //{"caseName": "min varuint32", "typeName": "varuint32", "expectedValue": "0", "json": Varuint32(0), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min float32", "typeName": "float32", "expectedValue": "01000000", "json": "{\"testField\":0.000000000000000000000000000000000000000000001401298464324817}", "expectedError": nil},
		{"caseName": "max float32", "typeName": "float32", "expectedValue": "ffff7f7f", "json": "{\"testField\":340282346638528860000000000000000000000}", "expectedError": nil},
		{"caseName": "min float64", "typeName": "float64", "expectedValue": "0100000000000000", "json": "{\"testField\":0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005}", "expectedError": nil},
		{"caseName": "max float64", "typeName": "float64", "expectedValue": "ffffffffffffef7f", "json": "{\"testField\":179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000}", "expectedError": nil},
		{"caseName": "float128 unsupported", "typeName": "float128", "expectedValue": "", "json": "{\"testField\":0}", "expectedError": fmt.Errorf("writing field: float128 support not implemented")},
		{"caseName": "bool true", "typeName": "bool", "expectedValue": "01", "json": "{\"testField\":true}", "expectedError": nil},
		{"caseName": "bool false", "typeName": "bool", "expectedValue": "00", "json": "{\"testField\":false}", "expectedError": nil},
		{"caseName": "time_point", "typeName": "time_point", "expectedValue": "ffffffffffffffff", "json": "{\"testField\":18446744073709551615}", "expectedError": nil},
		{"caseName": "time_point_sec", "typeName": "time_point_sec", "expectedValue": "ffffffff", "json": "{\"testField\":4294967295", "expectedError": nil},
		{"caseName": "block_timestamp_type", "typeName": "block_timestamp_type", "expectedValue": "76c52223", "json": "{\"testField\":\"2018-09-05T12:48:54-04:00\"}", "expectedError": nil},
		{"caseName": "Name", "typeName": "name", "expectedValue": "0000000000ea3055", "json": "{\"testField\":\"eosio\"}", "expectedError": nil},
		{"caseName": "Name", "typeName": "name", "expectedValue": "", "json": "{\"testField\":\"waytolongnametomakethetestcrash\"}", "expectedError": fmt.Errorf("writing field: name: waytolongnametomakethetestcrash is to long. expexted length of max 12 characters")},
		{"caseName": "bytes", "typeName": "bytes", "expectedValue": "0e746869732e69732e612e74657374", "json": "{\"testField\":\"746869732e69732e612e74657374\"}", "expectedError": nil},
		{"caseName": "checksum160", "typeName": "checksum160", "expectedValue": "0000000000000000000000000000000000000000", "json": "{\"testField\":\"0000000000000000000000000000000000000000\"}", "expectedError": nil},
		{"caseName": "checksum256", "typeName": "checksum256", "expectedValue": "0000000000000000000000000000000000000000000000000000000000000000", "json": "{\"testField\":\"0000000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": nil},
		{"caseName": "checksum512", "typeName": "checksum512", "expectedValue": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "json": "{\"testField\":\"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": nil},
		{"caseName": "checksum160 to long", "typeName": "checksum160", "expectedValue": "", "json": "{\"testField\":\"10000000000000000000000000000000000000000\"}", "expectedError": fmt.Errorf("writing field: checksum160: expected length of 40 got 41 for value 10000000000000000000000000000000000000000")},
		{"caseName": "checksum256 to long", "typeName": "checksum256", "expectedValue": "", "json": "{\"testField\":\"10000000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": fmt.Errorf("writing field: checksum256: expected length of 64 got 65 for value 10000000000000000000000000000000000000000000000000000000000000000")},
		{"caseName": "checksum512 to long", "typeName": "checksum512", "expectedValue": "", "json": "{\"testField\":\"100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": fmt.Errorf("writing field: checksum512: expected length of 128 got 129 for value 100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")},
		{"caseName": "public_key", "typeName": "public_key", "expectedValue": "00000000000000000000000000000000000000000000000000000000000000000000", "json": "{\"testField\":\"EOS1111111111111111111111111111111114T1Anm\"}", "expectedError": nil},
		{"caseName": "public_key err", "typeName": "public_key", "expectedValue": "", "json": "{\"testField\":\"EOS1111111111111111111111114T1Anm\"}", "expectedError": fmt.Errorf("writing field: public_key: checkDecode: invalid checksum")},
		{"caseName": "signature", "typeName": "signature", "expectedValue": "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "json": "{\"testField\":\"SIG_K1_111111111111111111111111111111111111111111111111111111111111111116uk5ne\"}", "expectedError": nil},
		{"caseName": "symbol", "typeName": "symbol", "expectedValue": "0403454f53", "json": "{\"testField\":{\"Precision\":4,\"Symbol\":\"EOS\"}}", "expectedError": nil},
		{"caseName": "symbol_code", "typeName": "symbol_code", "expectedValue": "ffffffffffffffff", "json": "{\"testField\":18446744073709551615}", "expectedError": nil},
		{"caseName": "asset", "typeName": "asset", "expectedValue": "a08601000000000004454f5300000000", "json": "{\"testField\":\"10.0000 EOS\"}", "expectedError": nil},
		{"caseName": "extended_asset", "typeName": "extended_asset", "expectedValue": "0a0000000000000004454f5300000000202932c94c833055", "json": "{\"testField\":{\"asset\":\"0.0010 EOS\",\"Contract\":\"eoscanadacom\"}}", "expectedError": nil},
		{"caseName": "bad type", "typeName": "bad.type.1", "expectedValue": nil, "json": "{\"testField\":0}", "expectedError": fmt.Errorf("writing field of type [bad.type.1]: unknown type")},
		{"caseName": "optional present", "typeName": "string", "expectedValue": "0776616c75652e31", "json": "{\"testField\":\"value.1\"}", "expectedError": nil},
	}

	for _, c := range testCases {

		t.Run(c["caseName"].(string), func(t *testing.T) {
			var buffer bytes.Buffer
			encoder := NewEncoder(&buffer)

			abi := ABI{}
			fieldName := "test_field_name"
			result := gjson.Get(c["json"].(string), "testField")
			err := abi.writeField(encoder, fieldName, c["typeName"].(string), result)
			assert.Equal(t, c["expectedError"], err, c["caseName"])

			if c["expectedError"] == nil {
				assert.Equal(t, c["expectedValue"], hex.EncodeToString(buffer.Bytes()), c["caseName"])
			}
		})
	}
}

type mockWriterable interface {
	Write(p []byte) (n int, err error)
	Bytes() []byte
}
type mockWriter struct {
	length int
	err    error
}

func (w mockWriter) Write(p []byte) (n int, err error) {
	return w.length, w.err
}

func (w mockWriter) Bytes() []byte {
	return []byte{}
}
