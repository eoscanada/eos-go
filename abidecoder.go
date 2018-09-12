package eos

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/sjson"
)

func (a *ABI) DecodeAction(data []byte, actionName ActionName) ([]byte, error) {

	binaryDecoder := NewDecoder(data)
	action := a.ActionForName(actionName)
	if action == nil {
		return []byte{}, fmt.Errorf("action %s not found in abi", actionName)
	}

	return a.decode(binaryDecoder, action.Type)

}

func (a *ABI) decode(binaryDecoder *Decoder, structName string) ([]byte, error) {

	Logger.ABIDecoder.Println("Decoding struct:", structName)

	structure := a.StructForName(structName)
	if structure == nil {
		return []byte{}, fmt.Errorf("structure [%s] not found in abi", structName)
	}

	resultingJson := make([]byte, 0)
	if structure.Base != "" {
		Logger.ABIDecoder.Printf("Structure %s has base structure of type: %s\n", structName, structure.Base)
		var err error
		resultingJson, err = a.decode(binaryDecoder, structure.Base)
		if err != nil {
			return resultingJson, fmt.Errorf("decode base [%s]: %s", structName, err)
		}
	}

	return a.decodeFields(binaryDecoder, structure.Fields, resultingJson)
}

func (a *ABI) decodeFields(binaryDecoder *Decoder, fields []FieldDef, json []byte) ([]byte, error) {
	defer Logger.ABIDecoder.SetPrefix(Logger.ABIDecoder.Prefix())
	Logger.ABIDecoder.SetPrefix(Logger.ABIDecoder.Prefix() + "\t")
	defer Logger.Decoder.SetPrefix(Logger.Decoder.Prefix())
	Logger.Decoder.SetPrefix(Logger.Decoder.Prefix() + "\t")

	resultingJson := json
	for _, field := range fields {

		fieldName, isOptional, isArray := analyseFieldName(field.Name)
		typeName := a.TypeNameForNewTypeName(field.Type)
		if typeName != field.Type {
			Logger.ABIDecoder.Printf("[%s] is an alias of [%s]\n", field.Type, typeName)
		}

		var err error
		resultingJson, err = a.decodeField(binaryDecoder, fieldName, typeName, isOptional, isArray, resultingJson)
		if err != nil {
			return []byte{}, fmt.Errorf("decoding fields: %s", err)
		}
	}

	return resultingJson, nil
}

func (a *ABI) decodeField(binaryDecoder *Decoder, fieldName string, fieldType string, isOptional bool, isArray bool, json []byte) ([]byte, error) {

	Logger.ABIDecoder.Printf("Decoding field [%s] of type [%s]\n", fieldName, fieldType)
	resultingJson := json
	if isOptional {
		Logger.ABIDecoder.Printf("Field [%s] is optional\n", fieldName)
		b, err := binaryDecoder.ReadByte()
		if err != nil {
			return resultingJson, fmt.Errorf("decoding field [%s] optional flag: %s", fieldName, err)
		}

		if b == 0 {
			Logger.ABIDecoder.Printf("Field [%s] is not present\n", fieldName)
			return resultingJson, nil
		}
	}

	if isArray {
		length, err := binaryDecoder.ReadUvarint64()
		if err != nil {
			return resultingJson, fmt.Errorf("reading field [%s] array length: %s", fieldName, err)
		}

		if length == 0 {
			resultingJson, err = sjson.SetBytes(resultingJson, fieldName, []interface{}{})
			if err != nil {
				return resultingJson, fmt.Errorf("reading field [%s] setting empty array: %s", fieldName, err)
			}
		}

		for i := uint64(0); i < length; i++ {
			Logger.ABIDecoder.Printf("\tAdding value for field: [%s] at index [%d]\n", fieldName, i)
			indexedFieldName := fmt.Sprintf("%s.%d", fieldName, i)
			resultingJson, err = a.read(binaryDecoder, indexedFieldName, fieldType, resultingJson)
			if err != nil {
				return resultingJson, fmt.Errorf("reading field [%s] index [%d]: %s", fieldName, i, err)
			}
		}

		return resultingJson, nil

	}

	resultingJson, err := a.read(binaryDecoder, fieldName, fieldType, resultingJson)
	if err != nil {
		return resultingJson, fmt.Errorf("decoding field [%s] of type [%s]: %s", fieldName, fieldType, err)
	}
	Logger.ABIDecoder.Printf("Set value: [%s] for field: [%s]\n", resultingJson, fieldName)
	return resultingJson, nil
}

func (a *ABI) read(binaryDecoder *Decoder, fieldName string, fieldType string, json []byte) ([]byte, error) {
	structure := a.StructForName(fieldType)

	if structure != nil {
		Logger.ABIDecoder.Printf("Field [%s] is a structure\n", fieldName)
		structureJson, err := a.decodeFields(binaryDecoder, structure.Fields, []byte{})
		if err != nil {
			return []byte{}, err
		}
		Logger.ABIDecoder.Printf("Setting [%s] for field [%s]", structureJson, fieldName)
		return sjson.SetRawBytes(json, fieldName, structureJson)
	}

	var value interface{}
	var err error
	switch fieldType {
	case "int8":
		value, err = binaryDecoder.ReadInt8()
	case "uint8":
		value, err = binaryDecoder.ReadUInt8()
	case "int16":
		value, err = binaryDecoder.ReadInt16()
	case "uint16":
		value, err = binaryDecoder.ReadUint16()
	case "int32":
		value, err = binaryDecoder.ReadInt32()
	case "uint32":
		value, err = binaryDecoder.ReadUint32()
	case "int64":
		value, err = binaryDecoder.ReadInt64()
	case "uint64":
		value, err = binaryDecoder.ReadUint64()
	case "int128":
		err = fmt.Errorf("int128 support not implemented")
	case "uint128":
		err = fmt.Errorf("uint128 support not implemented")
	case "varint32":
		value, err = binaryDecoder.ReadVarint32()
	case "varuint32":
		value, err = binaryDecoder.ReadUvarint32()
	case "float32":
		value, err = binaryDecoder.ReadFloat32()
	case "float64":
		value, err = binaryDecoder.ReadFloat64()
	case "float128":
		err = fmt.Errorf("float128 support not implemented")
	case "bool":
		value, err = binaryDecoder.ReadBool()
	case "time_point":
		timePoint, e := binaryDecoder.ReadTimePoint()
		if e == nil {
			t := time.Unix(0, int64(timePoint))
			value = t.UTC().Format("2006-01-02T15:04:05.999")
			err = e
		}
	case "time_point_sec":
		timePointSec, e := binaryDecoder.ReadTimePoint()
		if e == nil {
			t := time.Unix(0, int64(timePointSec))
			value = t.UTC().Format("2006-01-02T15:04:05")
			err = e
		}
	case "block_timestamp_type":
		value, err = binaryDecoder.ReadBlockTimestamp()
		if err == nil {
			value = value.(BlockTimestamp).Time.UTC().Format("2006-01-02T15:04:05")
		}
	case "name":
		value, err = binaryDecoder.ReadName()
	case "bytes":
		value, err = binaryDecoder.ReadByteArray()
		if err == nil {
			value = hex.EncodeToString(value.([]byte))
		}
	case "string":
		value, err = binaryDecoder.ReadString()
	case "checksum160":
		value, err = binaryDecoder.ReadChecksum160()
	case "checksum256":
		value, err = binaryDecoder.ReadChecksum256()
	case "checksum512":
		value, err = binaryDecoder.ReadChecksum512()
	case "public_key":
		value, err = binaryDecoder.ReadPublicKey()
	case "signature":
		value, err = binaryDecoder.ReadSignature()
	case "symbol":
		value, err = binaryDecoder.ReadSymbol()
	case "symbol_code":
		value, err = binaryDecoder.ReadSymbolCode()
	case "asset":
		value, err = binaryDecoder.ReadAsset()
	case "extended_asset":
		value, err = binaryDecoder.ReadExtendedAsset()
	default:
		return nil, fmt.Errorf("read field of type [%s]: unknown type", fieldType)
	}

	if err != nil {
		return []byte{}, fmt.Errorf("read: %s", err)
	}

	return sjson.SetBytes(json, fieldName, value)

}

func analyseFieldName(fieldName string) (name string, isOptional bool, isArray bool) {

	if strings.HasSuffix(fieldName, "?") {
		return fieldName[0 : len(fieldName)-1], true, false
	}

	if strings.HasSuffix(fieldName, "[]") {
		return fieldName[0 : len(fieldName)-2], false, true
	}

	return fieldName, false, false
}
