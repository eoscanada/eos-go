package eos

import (
	"fmt"
	"io"
)

type ABIEncoder struct {
	abiReader  io.Reader
	eosEncoder *Encoder
	abi        *ABI
	pos        int
}

func NewABIEncoder(abiReader io.Reader, writer io.Writer) *ABIEncoder {

	return &ABIEncoder{
		eosEncoder: NewEncoder(writer),
		abiReader:  abiReader,
	}
}

func (e *ABIEncoder) Encode(actionName ActionName, v map[string]interface{}) error {

	abi, err := NewABI(e.abiReader)
	if err != nil {
		return fmt.Errorf("encode: %s", err)
	}
	e.abi = abi

	action := abi.ActionForName(actionName)
	if action == nil {
		return fmt.Errorf("action %s not found in abi", actionName)
	}

	return e.encode(action.Type, v)
}

func (e *ABIEncoder) encode(structName string, action map[string]interface{}) error {

	fmt.Println("Decoding struct:", structName)

	structure := e.abi.StructForName(structName)
	if structure == nil {
		return fmt.Errorf("encode: structure [%s] not found in abi", structName)
	}

	if structure.Base != "" {
		fmt.Printf("Structure: %s has base structure of type: %s\n", structName, structure.Base)
		err := e.encode(structure.Base, action)
		if err != nil {
			return fmt.Errorf("decode base [%s]: %s", structName, err)
		}
	}

	return e.encodeFields(structure.Fields, action)
}
func (e *ABIEncoder) encodeFields(fields []FieldDef, action map[string]interface{}) error {

	for _, field := range fields {

		fmt.Printf("Encoding field [%s] of type [%s]\n", field.Name, field.Type)

		fieldName, isOptinal, isArray := analyseFieldName(field.Name)
		typeName := e.abi.TypeNameForNewTypeName(field.Type)
		if typeName != field.Type {
			fmt.Printf("-- type [%s] is an alias of [%s]\n", field.Type, typeName)
		}

		e.encodeField(fieldName, isOptinal, isArray, action)
	}

	return nil

}

func (e *ABIEncoder) encodeField(fieldName string, isOptional bool, isArray bool, action map[string]interface{}) (err error) {

	value := action[fieldName]
	if isOptional {
		if value == nil {
			return e.eosEncoder.writeByte(0)
		} else {
			if e := e.eosEncoder.writeByte(1); e != nil {
				return e
			}

			fmt.Println("Writing optional value:", value)
			e.eosEncoder.Encode(value)
		}

	} else if value == nil {
		return fmt.Errorf("encode field: none optional field [%s] as a nil value", fieldName)
	}

	return e.eosEncoder.Encode(value)
}
