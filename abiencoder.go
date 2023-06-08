package eos

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
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
		return nil, fmt.Errorf("encode action: %w", err)
	}
	return buffer.Bytes(), nil
}

func (a *ABI) EncodeTable(tableName TableName, json []byte) ([]byte, error) {
	table := a.TableForName(tableName)
	if table == nil {
		return nil, fmt.Errorf("encode table: table %s not found in abi", tableName)
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)

	err := a.encode(encoder, table.Type, json)
	if err != nil {
		return nil, fmt.Errorf("encode table: %w", err)
	}
	return buffer.Bytes(), nil
}

func (a *ABI) EncodeStruct(structName string, json []byte) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)

	err := a.encode(encoder, structName, json)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	return buffer.Bytes(), nil
}

func (a *ABI) encode(binaryEncoder *Encoder, structName string, json []byte) error {
	if tracer.Enabled() {
		zlog.Debug("abi encode struct", zap.String("name", structName))
	}

	structure := a.StructForName(structName)
	if structure == nil {
		return fmt.Errorf("encode struct [%s] not found in abi", structName)
	}

	if structure.Base != "" {
		if tracer.Enabled() {
			zlog.Debug("struct has base struct", zap.String("struct", structName), zap.String("base", structure.Base))
		}
		err := a.encode(binaryEncoder, structure.Base, json)
		if err != nil {
			return fmt.Errorf("encode base [%s]: %w", structName, err)
		}
	}
	err := a.encodeFields(binaryEncoder, structure.Fields, json)
	return err
}
func (a *ABI) encodeFields(binaryEncoder *Encoder, fields []FieldDef, json []byte) error {

	if tracer.Enabled() {
		defer func(prev *zap.Logger) { zlog = prev }(zlog)
		zlog = zlog.Named("fields")
		defer func(prev *zap.Logger) { zlog = prev }(zlog)
		zlog = zlog.Named("fields")
	}

	for _, field := range fields {

		if tracer.Enabled() {
			zlog.Debug("encode field", zap.String("name", field.Name), zap.String("type", field.Type))
		}

		fieldType, isOptional, isArray, _ := analyzeFieldType(field.Type)
		typeName, isAlias := a.TypeNameForNewTypeName(fieldType)
		fieldName := field.Name
		if isAlias && tracer.Enabled() {
			zlog.Debug("type is an alias", zap.String("from", field.Type), zap.String("to", typeName))
		}

		err := a.encodeField(binaryEncoder, fieldName, typeName, isOptional, isArray, json)
		if err != nil {
			return fmt.Errorf("encoding fields: %w", err)
		}
	}
	return nil
}

func (a *ABI) encodeField(binaryEncoder *Encoder, fieldName string, fieldType string, isOptional bool, isArray bool, json []byte) (err error) {

	if tracer.Enabled() {
		zlog.Debug("encode field json", zap.ByteString("json", json))
	}

	value := gjson.GetBytes(json, fieldName)
	if isOptional {
		if value.Exists() {
			if tracer.Enabled() {
				zlog.Debug("field is optional and present", zap.String("name", fieldName), zap.String("type", fieldType))
			}
			if e := binaryEncoder.writeByte(1); e != nil {
				return e
			}
		} else {
			if tracer.Enabled() {
				zlog.Debug("field is optional and *not* present", zap.String("name", fieldName), zap.String("type", fieldType))
			}
			return binaryEncoder.writeByte(0)
		}

	} else if !value.Exists() {
		return fmt.Errorf("encode field: none optional field [%s] as a nil value", fieldName)
	}

	if isArray {

		if tracer.Enabled() {
			zlog.Debug("field is an array", zap.String("name", fieldName), zap.String("type", fieldType))
		}
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

	if tracer.Enabled() {
		zlog.Debug("write field", zap.String("name", fieldName), zap.String("type", fieldType), zap.String("json", value.Raw))
	}

	structure := a.StructForName(fieldType)
	if structure != nil {
		if tracer.Enabled() {
			zlog.Debug("field is a struct", zap.String("name", fieldName))
		}

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
	case "int32":
		i, err := valueToInt(fieldName, value, 32)
		if err != nil {
			return err
		}
		object = int32(i)
	case "varint32":
		v, err := strconv.ParseInt(value.Raw, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid int32 value %q", value.Raw)
		}

		object = Varint32(v)

	case "uint32":
		i, err := valueToUint(fieldName, value, 32)
		if err != nil {
			return err
		}
		object = uint32(i)
	case "varuint32":
		v, err := strconv.ParseUint(value.Raw, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid uint32 value %q", value.Raw)
		}

		object = Varuint32(v)

	case "int64":
		var in Int64
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return fmt.Errorf("encoding int64: %w", err)
		}
		object = in
	case "uint64":
		var in Uint64
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return fmt.Errorf("encoding uint64: %w", err)
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
			return fmt.Errorf("writing field: time_point_sec: %w", err)
		}
		object = TimePointSec(t.UTC().Unix())
	case "time_point":
		t, err := time.Parse("2006-01-02T15:04:05.999", value.Str)
		if err != nil {
			return fmt.Errorf("writing field: time_point: %w", err)
		}
		object = TimePoint(t.UTC().Nanosecond() / int(time.Millisecond))
	case "block_timestamp_type":
		t, err := time.Parse("2006-01-02T15:04:05.999999-07:00", value.Str)
		if err != nil {
			return fmt.Errorf("writing field: block_timestamp_type: %w", err)
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
			return fmt.Errorf("writing field: bytes: %w", err)
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
			return fmt.Errorf("writing field: checksum160: %w", err)
		}
		object = Checksum160(data)
	case "checksum256":
		if len(value.Str) != 64 {
			return fmt.Errorf("writing field: checksum256: expected length of 64 got %d for value %s", len(value.Str), value.String())
		}
		data, err := hex.DecodeString(value.Str)
		if err != nil {
			return fmt.Errorf("writing field: checksum256: %w", err)
		}
		object = Checksum256(data)
	case "checksum512":
		if len(value.Str) != 128 {
			return fmt.Errorf("writing field: checksum512: expected length of 128 got %d for value %s", len(value.Str), value.String())
		}
		data, err := hex.DecodeString(value.String())
		if err != nil {
			return fmt.Errorf("writing field: checksum512: %w", err)
		}
		object = Checksum512(data)
	case "public_key":
		pk, err := ecc.NewPublicKey(value.String())
		if err != nil {
			return fmt.Errorf("writing field: public_key: %w", err)
		}
		object = pk
	case "signature":
		signature, err := ecc.NewSignature(value.String())
		if err != nil {
			return fmt.Errorf("writing field: public_key: %w", err)
		}
		object = signature
	case "symbol":
		parts := strings.Split(value.Str, ",")
		if len(parts) != 2 {
			return fmt.Errorf("writing field: symbol: symbol should be of format '4,EOS'")
		}

		i, err := strconv.ParseUint(parts[0], 10, 8)
		if err != nil {
			return fmt.Errorf("writing field: symbol: %w", err)
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
			return fmt.Errorf("writing field: asset: %w", err)
		}
		object = asset
	case "extended_asset":
		var extendedAsset ExtendedAsset
		err := json.Unmarshal([]byte(value.Raw), &extendedAsset)
		if err != nil {
			return fmt.Errorf("writing field: extended_asset: %w", err)
		}
		object = extendedAsset
	default:
		return fmt.Errorf("writing field of type [%s]: unknown type", fieldType)
	}

	if tracer.Enabled() {
		zlog.Debug("write object", zap.Reflect("value", object))
	}

	return binaryEncoder.Encode(object)
}

func valueToInt(fieldName string, value gjson.Result, bitSize int) (int64, error) {
	// Compatible with conversion of quoted strings to int
	i, err := strconv.ParseInt(strings.Trim(value.Raw, `"`), 10, bitSize)
	if err != nil {
		return i, fmt.Errorf("writing field: [%s] type int%d : %w", fieldName, bitSize, err)
	}
	return i, nil
}

func valueToUint(fieldName string, value gjson.Result, bitSize int) (uint64, error) {
	// Compatible with conversion of quoted strings to uint
	i, err := strconv.ParseUint(strings.Trim(value.Raw, `"`), 10, bitSize)
	if err != nil {
		return i, fmt.Errorf("writing field: [%s] type uint%d : %w", fieldName, bitSize, err)
	}
	return i, nil
}

func valueToFloat(fieldName string, value gjson.Result, bitSize int) (float64, error) {
	switch value.Raw {
	case "inf":
		return math.Inf(1), nil
	case "-inf":
		return math.Inf(-1), nil
	case "nan":
		return math.NaN(), nil
	default:
	}
	// Compatible with conversion of quoted strings to float
	f, err := strconv.ParseFloat(strings.Trim(value.Raw, `"`), bitSize)
	if err != nil {
		return f, fmt.Errorf("writing field: [%s] type float%d : %w", fieldName, bitSize, err)
	}
	return f, nil
}
