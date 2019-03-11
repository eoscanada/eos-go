package eos

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/eoscanada/eos-go/ecc"
	"go.uber.org/zap"

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
	if err != nil {
		return nil, fmt.Errorf("encode action: %s", err)
	}
	return buffer.Bytes(), nil
}

func (a *ABI) encode(binaryEncoder *Encoder, structureName string, json []byte) error {
	abiEncoderLog.Debug("abi encode struct", zap.String("name", structureName))

	structure := a.StructForName(structureName)
	if structure == nil {
		return fmt.Errorf("encode struct [%s] not found in abi", structureName)
	}

	if structure.Base != "" {
		abiEncoderLog.Debug("struct has base struct", zap.String("struct", structureName), zap.String("base", structure.Base))
		err := a.encode(binaryEncoder, structure.Base, json)
		if err != nil {
			return fmt.Errorf("encode base [%s]: %s", structureName, err)
		}
	}
	err := a.encodeFields(binaryEncoder, structure.Fields, json)
	return err
}
func (a *ABI) encodeFields(binaryEncoder *Encoder, fields []FieldDef, json []byte) error {

	defer func(prev *zap.Logger) { abiEncoderLog = prev }(abiEncoderLog)
	abiEncoderLog = encoderLog.Named("fields")
	defer func(prev *zap.Logger) { encoderLog = prev }(encoderLog)
	encoderLog = encoderLog.Named("fields")

	for _, field := range fields {

		abiEncoderLog.Debug("encode field", zap.String("name", field.Name), zap.String("type", field.Type))

		fieldType, isOptional, isArray := analyzeFieldType(field.Type)
		typeName, isAlias := a.TypeNameForNewTypeName(fieldType)
		fieldName := field.Name
		if isAlias {
			abiEncoderLog.Debug("type is an alias", zap.String("from", field.Type), zap.String("to", typeName))
		}

		err := a.encodeField(binaryEncoder, fieldName, typeName, isOptional, isArray, json)
		if err != nil {
			return fmt.Errorf("encoding fields: %s", err)
		}
	}
	return nil
}

func (a *ABI) encodeField(binaryEncoder *Encoder, fieldName string, fieldType string, isOptional bool, isArray bool, json []byte) (err error) {

	abiEncoderLog.Debug("encode field json", zap.ByteString("json", json))

	value := gjson.GetBytes(json, fieldName)
	if isOptional {
		if value.Exists() {
			abiEncoderLog.Debug("field is optional and present", zap.String("name", fieldName), zap.String("type", fieldType))
			if e := binaryEncoder.writeByte(1); e != nil {
				return e
			}
		} else {
			abiEncoderLog.Debug("field is optional and *not* present", zap.String("name", fieldName), zap.String("type", fieldType))
			return binaryEncoder.writeByte(0)
		}

	} else if !value.Exists() {
		return fmt.Errorf("encode field: none optional field [%s] as a nil value", fieldName)
	}

	if isArray {

		abiEncoderLog.Debug("field is an array", zap.String("name", fieldName), zap.String("type", fieldType))
		if !value.IsArray() {
			return fmt.Errorf("encode field: expected array for field [%s] got [%s]", fieldName, value.Type.String())
		}

		results := value.Array()
		binaryEncoder.writeUVarInt(len(results))

		for _, r := range results {
			a.writeField(binaryEncoder, fieldName, fieldType, r)
		}

		return nil
	}

	return a.writeField(binaryEncoder, fieldName, fieldType, value)
}

func (a *ABI) writeField(binaryEncoder *Encoder, fieldName string, fieldType string, value gjson.Result) error {

	abiEncoderLog.Debug("write field", zap.String("name", fieldName), zap.String("type", fieldType), zap.String("json", value.Raw))

	structure := a.StructForName(fieldType)
	if structure != nil {
		abiEncoderLog.Debug("field is a struct", zap.String("name", fieldName))

		err := a.encodeFields(binaryEncoder, structure.Fields, []byte(value.Raw))
		if err != nil {
			return err
		}
		return nil
	}

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
	case "uint32", "varuint32":
		i, err := valueToUint(fieldName, value, 32)
		if err != nil {
			return err
		}
		object = uint32(i)
	case "int64":
		var in Int64
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return fmt.Errorf("encoding int64: %s", err)
		}
		object = in
	case "uint64":
		var in Uint64
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return fmt.Errorf("encoding uint64: %s", err)
		}
		object = in
	case "int128":
		var in Int128
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return err
		}
		object = in
	case "uint128":
		var in Uint128
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return err
		}
		object = in
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
		var in Float128
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return err
		}
		object = in
	case "bool":
		object = value.Bool()
	case "time_point_sec":
		t, err := time.Parse("2006-01-02T15:04:05", value.Str)
		if err != nil {
			return fmt.Errorf("writing field: time_point_sec: %s", err)
		}
		object = TimePointSec(t.UTC().Second())
	case "time_point":
		t, err := time.Parse("2006-01-02T15:04:05.999", value.Str)
		if err != nil {
			return fmt.Errorf("writing field: time_point: %s", err)
		}
		object = TimePoint(t.UTC().Nanosecond() / int(time.Millisecond))
	case "block_timestamp_type":
		t, err := time.Parse("2006-01-02T15:04:05.999999-07:00", value.Str)
		if err != nil {
			return fmt.Errorf("writing field: block_timestamp_type: %s", err)
		}
		object = BlockTimestamp{
			Time: t,
		}
	case "name":
		if len(value.Str) > 12 {
			return fmt.Errorf("writing field: name: %s is to long. expected length of max 12 characters", value.Str)
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
		parts := strings.Split(value.Str, ",")
		if len(parts) != 2 {
			return fmt.Errorf("writing field: symbol: symbol should be of format '4,EOS'")
		}

		i, err := strconv.ParseUint(parts[0], 10, 8)
		if err != nil {
			return fmt.Errorf("writing field: symbol: %s", err)
		}
		object = Symbol{
			Precision: uint8(i),
			Symbol:    parts[1],
		}

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

	abiEncoderLog.Debug("write object", zap.Reflect("value", object))

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
