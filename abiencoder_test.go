package eos

import (
	"fmt"
	"os"
	"testing"

	"bytes"

	"strings"

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
			{"name":"struct.1.field.2", "type":"struct.name.3"},
			{"name":"struct.1.field.3?", "type":"string"},
			{"name":"struct.1.field.4?", "type":"string"}
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
var abiData = ABIMap{
	"struct.2.field.1": "struct.2.field.1.value",
	"struct.1.field.1": Name("eoscanadacom"),
	"struct.1.field.2": ABIMap{
		"struct.3.field.1": "struct.3.field.1.value",
	},
	"struct.1.field.3": "struct.1.field.3.value",
	//"struct.1.field.4": "struct.1.field.4.value",
}

func TestABIEncoder_Encode(t *testing.T) {
	Logger.Decoder.SetOutput(os.Stdout)
	Logger.Encoder.SetOutput(os.Stdout)
	Logger.ABIEncoder.SetOutput(os.Stdout)
	Logger.ABIDecoder.SetOutput(os.Stdout)

	testCases := []map[string]interface{}{
		{"caseName": "sunny path", "actionName": "action.name.1", "expectedError": nil, "writer": new(bytes.Buffer), "abi": abiString},
		{"caseName": "missing action", "actionName": "action.name.missing", "expectedError": fmt.Errorf("action action.name.missing not found in abi"), "writer": new(bytes.Buffer), "abi": abiString},
		{"caseName": "abi reader error", "actionName": "action.name.missing", "expectedError": fmt.Errorf("encode: read abi: unexpected EOF"), "writer": new(bytes.Buffer), "abi": "{"},
	}

	for _, c := range testCases {
		caseName := c["caseName"].(string)
		t.Run(caseName, func(t *testing.T) {

			buf := c["writer"].(mockWriterable)
			encoder := NewABIEncoder(strings.NewReader(c["abi"].(string)), buf)
			err := encoder.Encode(ActionName(c["actionName"].(string)), abiData)
			assert.Equal(t, c["expectedError"], err)

			if c["expectedError"] != nil {
				return
			}

			decoder := NewABIDecoder(buf.Bytes(), strings.NewReader(abiString))
			result := make(ABIMap)
			err = decoder.Decode(result, ActionName(c["actionName"].(string)))
			assert.NoError(t, err)

			assert.Equal(t, abiData, result)
			//fmt.Println(result)
		})
	}
}
func TestABIEncoder_encodeMissingActionStruct(t *testing.T) {

	abiString := `
{
	"version": "eosio::abi/1.0",
	"types": [{
		"new_type_name": "new.type.name.1",
		"type": "name"
	}],
	"structs": [
	],
   "actions": [{
		"name": "action.name.1",
		"type": "struct.name.1",
		"ricardian_contract": ""
   }]
}
`

	buf := new(bytes.Buffer)
	encoder := NewABIEncoder(strings.NewReader(abiString), buf)

	err := encoder.Encode("action.name.1", map[string]interface{}{})
	assert.Equal(t, fmt.Errorf("encode: structure [struct.name.1] not found in abi"), err)
}

func TestABIEncoder_encodeErrorInBase(t *testing.T) {

	abiString := `
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
			{"name":"struct.1.field.1", "type":"new.type.name.1"}
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

	buf := new(bytes.Buffer)
	encoder := NewABIEncoder(strings.NewReader(abiString), buf)

	err := encoder.Encode("action.name.1", map[string]interface{}{})
	assert.Equal(t, fmt.Errorf("encode base [struct.name.1]: encode: structure [struct.name.2] not found in abi"), err)
}

func TestABI_EncodeEncodeName(t *testing.T) {

	abiString := `
{
	"version": "eosio::abi/1.0",
	"types": [{
		"new_type_name": "new.type.name.1",
		"type": "name"
	}],
	"structs": [
	{
		"name": "struct.name.1",
		"base": "",
		"fields": [
			{"name":"struct.1.field.1", "type":"new.type.name.1"}
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
	abiData := map[string]interface{}{
		"struct.1.field.1": Name("struct.1.field.1.value"),
	}
	buf := new(bytes.Buffer)
	encoder := NewABIEncoder(strings.NewReader(abiString), buf)
	err := encoder.Encode(ActionName("action.name.1"), abiData)
	assert.Equal(t, nil, err)

	decoder := NewABIDecoder(buf.Bytes(), strings.NewReader(abiString))
	result := make(ABIMap)
	err = decoder.Decode(result, ActionName("action.name.1"))
	assert.NoError(t, err)
	fmt.Println(result)

}

func TestABIEncoder_encodeFields(t *testing.T) {

	abiString = `
{
	"version": "eosio::abi/1.0",
	"types": [
		{
			"new_type_name": "new.type.1",
			"type": "string"
		}
	],
	"structs": [
		{
			"name": "struct.name.1",
			"base": "",
			"fields": [
				{"name":"struct.1.field.1", "type":"new.type.name.1"}
			]
    	}
	],
   	"actions": []
}
	
`
	testCases := []map[string]interface{}{
		{"caseName": "sunny path", "fields": []FieldDef{{Name: "field.name.1", Type: "new.type.1"}}, "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "expectedError": nil, "writer": new(bytes.Buffer)},
		{"caseName": "encodeField error", "fields": []FieldDef{{Name: "field.name.1", Type: "new.type.1"}}, "actionMap": map[string]interface{}{}, "expectedError": fmt.Errorf("encoding fields: encode field: none optional field [field.name.1] as a nil value"), "writer": new(bytes.Buffer)},
		{"caseName": "embedded struct wrong type", "fields": []FieldDef{{Name: "field.name.1", Type: "struct.name.1"}}, "actionMap": map[string]interface{}{"field.name.1": map[string]interface{}{}}, "expectedError": fmt.Errorf("encode fields: structure field [field.name.1] expected to be of type ABIMap"), "writer": new(bytes.Buffer)},
		{"caseName": "encodeField error embedded struct", "fields": []FieldDef{{Name: "field.name.1", Type: "struct.name.1"}}, "actionMap": map[string]interface{}{"field.name.1": make(ABIMap)}, "expectedError": fmt.Errorf("encoding fields: encode field: none optional field [struct.1.field.1] as a nil value"), "writer": new(bytes.Buffer)},
	}

	for _, c := range testCases {
		caseName := c["caseName"].(string)
		t.Run(caseName, func(t *testing.T) {

			buf := c["writer"].(mockWriterable)
			encoder := NewABIEncoder(nil, buf)

			abi, err := NewABI(strings.NewReader(abiString))
			assert.NoError(t, err)
			encoder.abi = abi

			err = encoder.encodeFields(c["fields"].([]FieldDef), c["actionMap"].(map[string]interface{}))
			assert.Equal(t, c["expectedError"], err)
		})
	}
}

func TestABIEncoder_encodeField(t *testing.T) {

	testCases := []map[string]interface{}{
		{"caseName": "sunny path", "fieldName": "field.name.1", "fieldType": "string", "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "isOptional": false, "isArray": false, "expectedError": nil, "writer": new(bytes.Buffer)},
		{"caseName": "optional with value", "fieldName": "field.name.1", "fieldType": "string", "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "isOptional": true, "isArray": false, "expectedError": nil, "writer": new(bytes.Buffer)},
		{"caseName": "optional with no value", "fieldName": "field.name.1", "fieldType": "string", "actionMap": map[string]interface{}{}, "isOptional": true, "isArray": false, "expectedError": nil, "writer": new(bytes.Buffer)},
		{"caseName": "none optional with nil value", "fieldName": "field.name.2", "fieldType": "string", "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "isOptional": false, "isArray": false, "expectedError": fmt.Errorf("encode field: none optional field [field.name.2] as a nil value"), "writer": new(bytes.Buffer)},
		{"caseName": "optional write flag err", "fieldName": "field.name.1", "fieldType": "string", "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "isOptional": true, "isArray": false, "expectedError": fmt.Errorf("mocked error"), "writer": mockWriter{0, fmt.Errorf("mocked error")}},
	}

	for _, c := range testCases {
		caseName := c["caseName"].(string)
		t.Run(caseName, func(t *testing.T) {
			buf := c["writer"].(mockWriterable)
			encoder := NewABIEncoder(nil, buf)

			action := c["actionMap"].(map[string]interface{})
			fieldName := c["fieldName"].(string)
			fieldType := c["fieldType"].(string)
			isOptional := c["isOptional"].(bool)
			isArray := c["isArray"].(bool)
			expectedError := c["expectedError"]

			err := encoder.encodeField(fieldName, isOptional, isArray, action)
			assert.Equal(t, expectedError, err, caseName)
			if expectedError != nil {
				return
			}

			decoder := NewABIDecoder(buf.Bytes(), nil)
			result := make(ABIMap)
			err = decoder.decodeField(fieldName, fieldType, isOptional, isArray, result)
			assert.NoError(t, err, caseName)

			assert.Equal(t, action[fieldName], result[fieldName], caseName)

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
