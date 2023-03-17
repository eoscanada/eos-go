package eos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func (a *ABI) DecodeAction(data []byte, actionName ActionName) ([]byte, error) {

	binaryDecoder := NewDecoder(data)
	action := a.ActionForName(actionName)
	if action == nil {
		return nil, fmt.Errorf("action %s not found in abi", actionName)
	}

	builtStruct, err := a.decode(binaryDecoder, action.Type)
	if err != nil {
		return nil, err
	}
	return json.Marshal(builtStruct)
}

func (a *ABI) DecodeActionResult(data []byte, actionName ActionName) ([]byte, error) {

	binaryDecoder := NewDecoder(data)
	actionResult := a.ActionResultForName(actionName)
	if actionResult == nil {
		return nil, fmt.Errorf("action_result %s not found in abi", actionName)
	}

	res, err := a.resolveField(binaryDecoder, actionResult.ResultType)
	if err != nil {
		return nil, err
	}

	return json.Marshal(res)
}

func (a *ABI) DecodeTableRow(tableName TableName, data []byte) ([]byte, error) {
	binaryDecoder := NewDecoder(data)
	tbl := a.TableForName(tableName)
	if tbl == nil {
		return nil, fmt.Errorf("table name %s not found in abi", tableName)
	}

	builtStruct, err := a.decode(binaryDecoder, tbl.Type)
	if err != nil {
		return nil, err
	}
	return json.Marshal(builtStruct)

}

func (a *ABI) DecodeTableRowTyped(tableType string, data []byte) ([]byte, error) {
	binaryDecoder := NewDecoder(data)
	builtStruct, err := a.decode(binaryDecoder, tableType)
	if err != nil {
		return nil, err
	}
	return json.Marshal(builtStruct)
}

func (a *ABI) Decode(binaryDecoder *Decoder, structName string) ([]byte, error) {
	builtStruct, err := a.decode(binaryDecoder, structName)
	if err != nil {
		return nil, err
	}
	return json.Marshal(builtStruct)
}

func (a *ABI) decode(binaryDecoder *Decoder, structName string) (map[string]interface{}, error) {
	if tracer.Enabled() {
		zlog.Debug("decode struct", zap.String("name", structName))
	}

	if variant := a.VariantForName(structName); variant != nil {
		out, err := binaryDecoder.ReadUvarint32()
		if err != nil {
			zlog.Error("error reading variant", zap.Error(err))
		}
		structName = variant.Types[out]
	}

	structure := a.StructForName(structName)
	if structure == nil {
		return nil, fmt.Errorf("structure [%s] not found in abi", structName)
	}

	builtStruct := map[string]interface{}{}
	if structure.Base != "" {
		if tracer.Enabled() {
			zlog.Debug("struct has base struct", zap.String("name", structName), zap.String("base", structure.Base))
		}

		baseName, isAlias := a.TypeNameForNewTypeName(structure.Base)
		if isAlias && tracer.Enabled() {
			zlog.Debug("base is an alias", zap.String("from", structure.Base), zap.String("to", baseName))
		}

		var err error
		builtStruct, err = a.decode(binaryDecoder, baseName)
		if err != nil {
			return nil, fmt.Errorf("decode base [%s]: %w", structName, err)
		}
	}

	return a.decodeFields(binaryDecoder, structure.Fields, builtStruct)
}

func (a *ABI) decodeFields(binaryDecoder *Decoder, fields []FieldDef, builtStruct map[string]interface{}) (out map[string]interface{}, err error) {
	out = builtStruct

	for _, field := range fields {
		resultingValue, err := a.resolveField(binaryDecoder, field.Type)
		if err != nil {
			return nil, fmt.Errorf("decoding field %s: %w", field.Name, err)
		}

		if resultingValue != skipField {
			out[field.Name] = resultingValue
		}
	}

	return
}

type skipFieldType int

var skipField = skipFieldType(0)

type field struct {
	name  string
	value interface{}
}

func (a *ABI) resolveField(binaryDecoder *Decoder, initialFieldType string) (out interface{}, err error) {
	// retrieve the fields characteristics, note we can be a few depth down here....
	fieldType, isOptional, isArray, isBinaryExtension := analyzeFieldType(initialFieldType)
	//fmt.Println("resolveField", isOptional, isArray, initialFieldType, fieldType)

	if tracer.Enabled() {
		zlog.Debug("analyzed field",
			zap.String("field_type", fieldType),
			zap.Bool("is_optional", isOptional),
			zap.Bool("is_array", isArray),
			zap.Bool("is_binaryExtension", isBinaryExtension),
		)
	}

	// check if this field is an alias
	aliasFieldType, isAlias := a.TypeNameForNewTypeName(fieldType)
	if isAlias {
		if tracer.Enabled() {
			zlog.Debug("type is an alias",
				zap.String("from", fieldType),
				zap.String("to", aliasFieldType),
			)
		}
		fieldType = aliasFieldType
	}

	// check if the field is a binary extension
	if isBinaryExtension && !binaryDecoder.hasRemaining() {
		if tracer.Enabled() {
			zlog.Debug("type is a binary extension and no more data, skipping field", zap.String("type", fieldType))
		}
		return skipField, nil
	}

	// check if the field is optional
	if isOptional {
		b, err := binaryDecoder.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("reading optional flag: %w", err)
		}

		if b == 0 {
			if tracer.Enabled() {
				zlog.Debug("field is not present")
			}
			if !a.fitNodeos {
				return skipField, nil
			}

			// TODO: Not sure about this right now
			fieldType = "null"
		}
	}

	// if we have an array, we will loop in it and handle the subField (note that if we have an array of ALIAS we would loop here)
	if isArray {
		retVal, err := a.readArray(binaryDecoder, fieldType)
		if err != nil {
			return nil, err
		}
		return retVal, nil
	}

	// if the fiels is not an array, but is an alias we need to re-resolve the field again
	if isAlias {
		return a.resolveField(binaryDecoder, fieldType)
	}

	return a.read(binaryDecoder, fieldType)
}

func (a *ABI) readArray(binaryDecoder *Decoder, fieldType string) ([]interface{}, error) {
	//fmt.Println("Read array", fieldType)
	length, err := binaryDecoder.ReadUvarint64()
	if err != nil {
		return nil, fmt.Errorf("reading array length: %w", err)
	}

	if length == 0 {
		return []interface{}{}, nil // we just want "[]" in the final output
	}

	var elements []interface{}
	for i := uint64(0); i < length; i++ {

		//fmt.Println("Field type", fieldType)
		retVal, err := a.resolveField(binaryDecoder, fieldType)
		if err != nil {
			return nil, fmt.Errorf("resolve array index [%d]: %w", i, err)
		}

		if retVal != skipField {
			elements = append(elements, retVal)
		}
	}

	return elements, nil
}

// Decodes the EOS ABIs built-in types
func (a *ABI) read(binaryDecoder *Decoder, fieldType string) (interface{}, error) {
	variant := a.VariantForName(fieldType)
	if variant != nil {
		variantIndex, err := binaryDecoder.ReadUvarint32()
		if err != nil {
			return nil, fmt.Errorf("unable to read variant type index: %w", err)
		}

		if int(variantIndex) >= len(variant.Types) {
			return nil, fmt.Errorf("variant type index is unknown, got type index %d, know up to index %d", variantIndex, len(variant.Types)-1)
		}

		variantFieldType := variant.Types[variantIndex]
		if tracer.Enabled() {
			zlog.Debug("field is a variant", zap.String("type", variantFieldType))
		}

		resolvedVariantFieldType, isAlias := a.TypeNameForNewTypeName(variantFieldType)
		if isAlias && tracer.Enabled() {
			zlog.Debug("variant type is an alias", zap.String("from", fieldType), zap.String("to", resolvedVariantFieldType))
		}

		fieldType = resolvedVariantFieldType
	}

	structure := a.StructForName(fieldType)
	if structure != nil {
		builtStruct, err := a.decode(binaryDecoder, fieldType)
		if err != nil {
			return nil, fmt.Errorf("decoding [%s]: %w", fieldType, err)
		}

		return builtStruct, nil
	}

	var value interface{}
	var err error
	switch fieldType {
	case "null":
		value = nil
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
		var val int64
		val, err = binaryDecoder.ReadInt64()
		value = Int64(val)
	case "uint64":
		var val uint64
		val, err = binaryDecoder.ReadUint64()
		value = Uint64(val)
	case "int128":
		v, e := binaryDecoder.ReadInt128()
		if e == nil {
			if a.fitNodeos {
				value = v.DecimalString()
			} else {
				value = v
			}
		}
		err = e
	case "uint128":
		v, e := binaryDecoder.ReadUint128("uint128")
		if e == nil {
			if a.fitNodeos {
				value = v.DecimalString()
			} else {
				value = v
			}
		}
		err = e
	case "varint32":
		value, err = binaryDecoder.ReadVarint32()
	case "varuint32":
		value, err = binaryDecoder.ReadUvarint32()
	case "float32":
		v, e := binaryDecoder.ReadFloat32()
		if e == nil {
			if a.fitNodeos {
				value = strconv.FormatFloat(float64(v), 'f', 17, 32)
			} else {
				value = json.RawMessage(strconv.FormatFloat(float64(v), 'f', -1, 64)) // as sjson does
			}
		}
		err = e
	case "float64":
		v, e := binaryDecoder.ReadFloat64()
		if e == nil {
			value = formatFloat(v, a.fitNodeos)
		}
		err = e
	case "float128":
		value, err = binaryDecoder.ReadUint128("float128")
	case "bool":
		if a.fitNodeos {
			value, err = binaryDecoder.ReadByte()
		} else {
			value, err = binaryDecoder.ReadBool()
		}
	case "time_point":
		timePoint, e := binaryDecoder.ReadTimePoint() //todo double check
		if e == nil {
			value = formatTimePoint(timePoint, a.fitNodeos)
		}
		err = e
	case "time_point_sec":
		timePointSec, e := binaryDecoder.ReadTimePointSec()
		if e == nil {
			t := time.Unix(int64(timePointSec), 0)
			value = t.UTC().Format("2006-01-02T15:04:05")
		}
		err = e
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
		symbol, e := binaryDecoder.ReadSymbol()
		err = e
		if err == nil {
			value = fmt.Sprintf("%d,%s", symbol.Precision, symbol.Symbol)
		}
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
		return nil, fmt.Errorf("read: %w", err)
	}

	if tracer.Enabled() {
		zlog.Debug("set field value",
			zap.Reflect("value", value),
		)
	}

	if variant != nil {
		// As a variant we need to include the field type in the json
		return []interface{}{fieldType, value}, nil
	}

	// t0 := time.Now()
	// defer func() {
	// 	//fmt.Println("Doing field", fieldName, value, time.Since(t0))
	// }()
	return value, nil
}

func analyzeFieldType(fieldType string) (typeName string, isOptional bool, isArray bool, isBinaryExtension bool) {
	if strings.HasSuffix(fieldType, "[]$") {
		return fieldType[0 : len(fieldType)-3], false, true, true
	}

	if strings.HasSuffix(fieldType, "?") {
		return fieldType[0 : len(fieldType)-1], true, false, false
	}

	if strings.HasSuffix(fieldType, "$") {
		return fieldType[0 : len(fieldType)-1], false, false, true
	}

	if strings.HasSuffix(fieldType, "[]") {
		return fieldType[0 : len(fieldType)-2], false, true, false
	}

	return fieldType, false, false, false
}

const standardTimePointFormat = "2006-01-02T15:04:05.999"
const nodeosTimePointFormat = "2006-01-02T15:04:05.000"

func formatTimePoint(timePoint TimePoint, shouldFitNodeos bool) string {
	t := time.Unix(0, int64(timePoint*1000))
	if shouldFitNodeos {
		return t.UTC().Format(nodeosTimePointFormat)
	}

	return t.UTC().Format(standardTimePointFormat)
}

const standardTimePointSecFormat = "2006-01-02T15:04:05"

func formatTimePointSec(timePoint TimePointSec) string {
	t := time.Unix(int64(timePoint), 0)

	return t.UTC().Format(standardTimePointSecFormat)
}

func formatFloat(v float64, fitNodeos bool) interface{} {
	switch {
	case math.IsInf(v, 1):
		return "inf"
	case math.IsInf(v, -1):
		return "-inf"
	case math.IsNaN(v): // cannot check equality on math.NaN()
		return "nan"
	default:
	}

	if fitNodeos {
		return strconv.FormatFloat(v, 'f', 17, 64)
	} else {
		return json.RawMessage(strconv.FormatFloat(float64(v), 'f', -1, 64))
	}

}
