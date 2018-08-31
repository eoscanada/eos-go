package eos

import (
	"fmt"
	"testing"

	"bytes"

	"strings"

	"github.com/stretchr/testify/assert"
)

func TestABIEncoder_Encode(t *testing.T) {

	testCases := []map[string]interface{}{
		{"caseName": "sunny path", "actionName": "action.name.1", "expectedError": nil, "writer": new(bytes.Buffer)},
	}

	for _, c := range testCases {
		caseName := c["caseName"].(string)
		t.Run(caseName, func(t *testing.T) {

			buf := c["writer"].(mockWriterable)
			encoder := NewABIEncoder(strings.NewReader(abiString), buf)
			err := encoder.Encode(ActionName(c["actionName"].(string)), abiData)
			assert.Equal(t, c["expectedError"], err)

			decoder := NewABIDecoder(buf.Bytes(), strings.NewReader(abiString))
			result := make(Result)
			err = decoder.Decode(result, ActionName(c["actionName"].(string)))
			assert.NoError(t, err)
			fmt.Println(result)
		})
	}
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
	"structs": [],
   	"actions": []
}
	
`
	testCases := []map[string]interface{}{
		{"caseName": "sunny path", "fields": []FieldDef{{Name: "field.name.1", Type: "new.type.1"}}, "actionMap": map[string]interface{}{"field.name.1": "field.1.value.1"}, "expectedError": nil, "writer": new(bytes.Buffer)},
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
			result := make(Result)
			Debug = true
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
