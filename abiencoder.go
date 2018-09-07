package eos

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/eoscanada/eos-go/ecc"

	"github.com/tidwall/gjson"
)

type ABIEncoder struct {
	abiReader  io.Reader
	eosEncoder *Encoder
	abi        *ABI
	pos        int
}

func (a *ABI) EncodeAction(actionName ActionName, json []byte) ([]byte, error) {

	action := a.ActionForName(actionName)
	if action == nil {
		return nil, fmt.Errorf("encode action: action %s not found in abi", actionName)
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)

	err := a.encode(encoder, action.Type, json)
	if ErrNotFound != nil {
		return nil, fmt.Errorf("encode action: %s", err)
	}
	return buffer.Bytes(), nil
}

func (a *ABI) encode(binaryEncoder *Encoder, structureName string, json []byte) error {
	Logger.ABIEncoder.Printf("Encoding structure [%s]\n", structureName)

	structure := a.StructForName(structureName)
	if structure == nil {
		return fmt.Errorf("encode: structure [%s] not found in abi", structureName)
	}

	if structure.Base != "" {
		Logger.ABIEncoder.Printf("Structure [%s] has base structure of type [%s]\n", structureName, structure.Base)
		err := a.encode(binaryEncoder, structure.Base, json)
		if err != nil {
			return fmt.Errorf("encode base [%s]: %s", structureName, err)
		}
	}
	err := a.encodeFields(binaryEncoder, structure.Fields, json)
	return err
}
func (a *ABI) encodeFields(binaryEncoder *Encoder, fields []FieldDef, json []byte) error {

	defer Logger.ABIEncoder.SetPrefix(Logger.ABIEncoder.Prefix())
	defer Logger.Encoder.SetPrefix(Logger.Encoder.Prefix())
	Logger.ABIEncoder.SetPrefix(Logger.ABIEncoder.Prefix() + "\t")
	Logger.Encoder.SetPrefix(Logger.Encoder.Prefix() + "\t")

	for _, field := range fields {

		Logger.ABIEncoder.Printf("Encoding field [%s] of type [%s]\n", field.Name, field.Type)

		fieldName, isOptional, isArray := analyseFieldName(field.Name)
		typeName := a.TypeNameForNewTypeName(field.Type)
		if typeName != field.Type {
			Logger.ABIEncoder.Printf("[%s] is an alias of [%s]\n", field.Type, typeName)
		}

		err := a.encodeField(binaryEncoder, fieldName, typeName, isOptional, isArray, json)
		if err != nil {
			return fmt.Errorf("encoding fields: %s", err)
		}
	}
	return nil
}

func (a *ABI) encodeField(binaryEncoder *Encoder, fieldName string, fieldType string, isOptional bool, isArray bool, json []byte) (err error) {

	value := gjson.GetBytes(json, fieldName)
	if isOptional {
		if value.Exists() {
			Logger.ABIEncoder.Printf("Field [%s] of type [%s] is option and present\n", fieldName, fieldType)
			if e := binaryEncoder.writeByte(1); e != nil {
				return e
			}
		} else {
			Logger.ABIEncoder.Printf("Field [%s] of type [%s] is option and not present\n", fieldName, fieldType)
			return binaryEncoder.writeByte(0)
		}

	} else if !value.Exists() {
		return fmt.Errorf("encode field: none optional field [%s] as a nil value", fieldName)
	}

	if isArray {

		Logger.ABIEncoder.Printf("Field [%s] of type [%s] is an array\n", fieldName, fieldType)
		if !value.IsArray() {
			return fmt.Errorf("encode field: expected array for field [%s] got [%s]", fieldName, value.Type.String())
		}

		results := value.Array()
		binaryEncoder.writeUVarInt(len(results))

		for _, r := range results {
			a.encodeField(binaryEncoder, fieldName, fieldType, false, false, []byte(r.Raw))
		}

		return nil
	}

	structure := a.StructForName(fieldType)
	if structure != nil {
		Logger.ABIEncoder.Printf("Field [%s] is a structure\n", fieldName)

		structureJSON := gjson.GetBytes(json, fieldName)
		err := a.encodeFields(binaryEncoder, structure.Fields, []byte(structureJSON.Raw))
		if err != nil {
			return err
		}
	}

	return a.writeField(binaryEncoder, fieldName, fieldType, value)
}

func (a *ABI) writeField(binaryEncoder *Encoder, fieldName string, fieldType string, value gjson.Result) error {

	Logger.ABIEncoder.Printf("Writing value [%s]\n", value.Raw)

	var object interface{}
	switch fieldType {
	case "int8":
		i, err := valueToInt(fieldName, value, 8)
		if err != nil {
			return err
		}
		object = int8(i)
	case "uint8":
		i, err := valueToUint(fieldName, value, 8)
		if err != nil {
			return err
		}
		object = uint8(i)
	case "int16":
		i, err := valueToInt(fieldName, value, 16)
		if err != nil {
			return err
		}
		object = int16(i)
	case "uint16":
		i, err := valueToUint(fieldName, value, 16)
		if err != nil {
			return err
		}
		object = uint16(i)
	case "int32", "varint32":
		i, err := valueToInt(fieldName, value, 32)
		if err != nil {
			return err
		}
		object = int32(i)
	case "uint32", "varuint32", "time_point_sec":
		i, err := valueToUint(fieldName, value, 32)
		if err != nil {
			return err
		}
		object = uint32(i)
	case "int64":
		i, err := valueToInt(fieldName, value, 64)
		if err != nil {
			return err
		}
		object = i
	case "uint64", "time_point":
		i, err := valueToUint(fieldName, value, 64)
		if err != nil {
			return err
		}
		object = i
	case "int128":
		return fmt.Errorf("writing field: int128 support not implemented")
	case "uint128":
		return fmt.Errorf("writing field: uint128 support not implemented")
	case "float32":
		f, err := valueToFloat(fieldName, value, 32)
		if err != nil {
			return err
		}
		object = float32(f)
	case "float64":
		f, err := valueToFloat(fieldName, value, 64)
		if err != nil {
			return err
		}
		object = f
	case "float128":
		return fmt.Errorf("writing field: float128 support not implemented")
	case "bool":
		object = value.Bool()
	case "block_timestamp_type":
		time, err := time.Parse("2006-01-02T15:04:05.999999-07:00", value.Str)
		if err != nil {
			return fmt.Errorf("writing field: block_timestamp_type: %s", err)
		}
		object = BlockTimestamp{
			Time: time,
		}
	case "name":
		if len(value.Str) > 12 {
			return fmt.Errorf("writing field: name: %s is to long. expexted length of max 12 characters", value.Str)
		}
		object = Name(value.Str)
	case "bytes":
		data, err := hex.DecodeString(value.String())
		if err != nil {
			return fmt.Errorf("writing field: bytes: %s", err)
		}
		object = data

	case "string":
		object = value.String()
	case "checksum160":
		if len(value.Str) != 40 {
			return fmt.Errorf("writing field: checksum160: expected length of 40 got %d for value %s", len(value.Str), value.String())
		}
		data, err := hex.DecodeString(value.Str)
		if err != nil {
			return fmt.Errorf("writing field: checksum160: %s", err)
		}
		object = Checksum160(data)
	case "checksum256":
		if len(value.Str) != 64 {
			return fmt.Errorf("writing field: checksum256: expected length of 64 got %d for value %s", len(value.Str), value.String())
		}
		data, err := hex.DecodeString(value.Str)
		if err != nil {
			return fmt.Errorf("writing field: checksum256: %s", err)
		}
		object = Checksum256(data)
	case "checksum512":
		if len(value.Str) != 128 {
			return fmt.Errorf("writing field: checksum512: expected length of 128 got %d for value %s", len(value.Str), value.String())
		}
		data, err := hex.DecodeString(value.String())
		if err != nil {
			return fmt.Errorf("writing field: checksum512: %s", err)
		}
		object = Checksum512(data)
	case "public_key":
		pk, err := ecc.NewPublicKey(value.String())
		if err != nil {
			return fmt.Errorf("writing field: public_key: %s", err)
		}
		object = pk
	case "signature":
		signature, err := ecc.NewSignature(value.String())
		if err != nil {
			return fmt.Errorf("writing field: public_key: %s", err)
		}
		object = signature
	case "symbol":
		var symbol Symbol
		err := json.Unmarshal([]byte(value.Raw), &symbol)
		if err != nil {
			return fmt.Errorf("writing field: symbol: %s", err)
		}
		object = symbol
	case "symbol_code":
		object = SymbolCode(value.Uint())
	case "asset":
		asset, err := NewAsset(value.String())
		if err != nil {
			return fmt.Errorf("writing field: asset: %s", err)
		}
		object = asset
	case "extended_asset":
		var extendedAsset ExtendedAsset
		err := json.Unmarshal([]byte(value.Raw), &extendedAsset)
		if err != nil {
			return fmt.Errorf("writing field: extended_asset: %s", err)
		}
		object = extendedAsset
	default:
		return fmt.Errorf("writing field of type [%s]: unknown type", fieldType)
	}

	Logger.ABIEncoder.Printf("Writing object [%s]\n", object)
	return binaryEncoder.Encode(object)

}

func valueToInt(fieldName string, value gjson.Result, bitSize int) (int64, error) {
	i, err := strconv.ParseInt(value.Raw, 10, bitSize)
	if err != nil {
		return i, fmt.Errorf("writing field: [%s] type int%d : %s", fieldName, bitSize, err)
	}
	return i, nil
}

func valueToUint(fieldName string, value gjson.Result, bitSize int) (uint64, error) {
	i, err := strconv.ParseUint(value.Raw, 10, bitSize)
	if err != nil {
		return i, fmt.Errorf("writing field: [%s] type uint%d : %s", fieldName, bitSize, err)
	}
	return i, nil
}

func valueToFloat(fieldName string, value gjson.Result, bitSize int) (float64, error) {
	f, err := strconv.ParseFloat(value.Raw, bitSize)
	if err != nil {
		return f, fmt.Errorf("writing field: [%s] type float%d : %s", fieldName, bitSize, err)
	}
	return f, nil
}
