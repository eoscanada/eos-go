package eos

import (
	"encoding/json"
	"fmt"
	"io"
)

var DEBUG = false

type Result map[string]interface{}

type ABIDecoder struct {
	eosDecoder *Decoder
	abi        *ABI
	abiReader  io.Reader
	pos        int
}

func NewABIDecoder(data []byte, abiReader io.Reader) *ABIDecoder {

	return &ABIDecoder{
		eosDecoder: NewDecoder(data),
		abiReader:  abiReader,
	}
}

func (d *ABIDecoder) Decode(result Result, actionName ActionName) error {

	abi := &ABI{}
	abiDecoder := json.NewDecoder(d.abiReader)
	err := abiDecoder.Decode(abi)
	if err != nil {
		return fmt.Errorf("read abi: %s", err)
	}
	d.abi = abi

	action := abi.ActionForName(actionName)
	if action == nil {
		return fmt.Errorf("action %s not found in abi", actionName)
	}

	return d.decode(action.Type, result)

}

func (d *ABIDecoder) decode(structName string, result Result) error {

	fmt.Println("Decoding struct:", structName)

	structure := d.abi.StructForName(structName)
	if structure == nil {
		return fmt.Errorf("structure [%s] not found in abi", structName)
	}

	if structure.Base != "" {
		fmt.Printf("Structure: %s has base structure of type: %s\n", structName, structure.Base)
		err := d.decode(structure.Base, result)
		if err != nil {
			return fmt.Errorf("decode base [%s]: %s", structName, err)
		}
	}

	return d.decodeFields(structure.Fields, result)
}

func (d *ABIDecoder) decodeFields(fields []FieldDef, result Result) error {

	for _, field := range fields {

		fmt.Printf("Decoding field [%s] of type [%s]\n", field.Name, field.Type)

		typeName := d.abi.TypeNameForNewTypeName(field.Type)
		if typeName != field.Type {
			fmt.Printf("-- type [%s] is an alias of [%s]\n", field.Type, typeName)
		}

		structure := d.abi.StructForName(typeName)
		if structure != nil {
			fmt.Printf("Field [%s] is a structure\n", field.Name)
			err := d.decodeFields(structure.Fields, result)
			if err != nil {
				return err
			}
		} else {
			err := d.read(field.Name, typeName, result)
			if err != nil {
				return fmt.Errorf("decode field [%s] of type[%s]: %s", field.Name, typeName, err)
			}
		}

	}

	return nil
}

func (d *ABIDecoder) read(fieldName string, fieldType string, result Result) (err error) {

	//todo: check for array flag "[]"
	//todo: check for optional "?"

	fmt.Printf("\tReading field [%s] of type [%s]\n", fieldName, fieldType)

	//optional := false

	var value interface{}

	switch fieldType {
	case "int8":
		value, err = d.eosDecoder.ReadInt8()
	case "uint8":
		value, err = d.eosDecoder.ReadUInt8()
	case "int16":
		value, err = d.eosDecoder.ReadInt16()
	case "uint16":
		value, err = d.eosDecoder.ReadUint16()
	case "int32":
		value, err = d.eosDecoder.ReadInt32()
	case "uint32":
		value, err = d.eosDecoder.ReadUint32()
	case "int64":
		value, err = d.eosDecoder.ReadInt64()
	case "uint64":
		value, err = d.eosDecoder.ReadUint64()
	case "int128":
		err = fmt.Errorf("read field: int128 support not implemented")
	case "uint128":
		err = fmt.Errorf("read field: uint128 support not implemented")
	case "varint32":
		value, err = d.eosDecoder.ReadVarint()
	case "varuint32":
		value, err = d.eosDecoder.ReadUvarint()
	case "float32":
		value, err = d.eosDecoder.ReadFloat32()
	case "float64":
		value, err = d.eosDecoder.ReadFloat64()
	case "float128":
		err = fmt.Errorf("read field: float128 support not implemented")
	case "bool":
		value, err = d.eosDecoder.ReadBool()
	case "time_point":
		value, err = d.eosDecoder.ReadTimePoint()
	case "time_point_sec":
		value, err = d.eosDecoder.ReadTimePointSec()
	case "block_timestamp_type":
		value, err = d.eosDecoder.ReadBlockTimestamp()
	case "name":
		value, err = d.eosDecoder.ReadName()
	case "bytes":
		value, err = d.eosDecoder.ReadByteArray()
	case "string":
		value, err = d.eosDecoder.ReadString()
	case "checksum160":
		value, err = d.eosDecoder.ReadChecksum160()
	case "checksum256":
		value, err = d.eosDecoder.ReadChecksum256()
	case "checksum512":
		value, err = d.eosDecoder.ReadChecksum512()
	case "public_key":
		value, err = d.eosDecoder.ReadPublicKey()
	case "signature":
		value, err = d.eosDecoder.ReadSignature()
	case "symbol":
		value, err = d.eosDecoder.ReadSymbol()
	case "symbol_code":
		value, err = d.eosDecoder.ReadSymbolCode()
	case "asset":
		value, err = d.eosDecoder.ReadAsset()
	case "extended_asset":
		value, err = d.eosDecoder.ReadExtendedAsset()
	default:
		return fmt.Errorf("read field [%s] of type [%s]: unknown type", fieldName, fieldType)
	}

	if err != nil {
		return err
	}

	//if value == nil && !optional {
	//	return fmt.Errorf("read field [%s] of type [%s]: is not optional", fieldName, fieldType)
	//}

	fmt.Printf("\tSet value: [%s] for field: [%s]\n", value, fieldName)
	result[fieldName] = value

	return
}
