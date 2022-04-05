package eos

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestABI_DecodeAction(t *testing.T) {
	abiReader := strings.NewReader(abiString)
	mockData := struct {
		BF1    string
		F1     Name
		F2     string
		F3FLAG byte //this a hack until we have the abi encoder
		F3     string
		F4FLAG byte //this a hack until we have the abi encoder
		F5     []string
	}{
		BF1:    "value_struct_2_field_1",
		F1:     Name("eoscanadacom"),
		F2:     "value_struct_3_field_1",
		F3FLAG: 1,
		F3:     "value_struct_1_field_3",
		F4FLAG: 0,
		F5:     []string{"value_struct_4_field_1_1", "value_struct_4_field_1_2", "value_struct_4_field_1_3"},
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	require.NoError(t, err)

	abi, err := NewABI(abiReader)
	require.NoError(t, err)

	json, err := abi.DecodeAction(buffer.Bytes(), "action_name_1")
	require.NoError(t, err)

	assert.Equal(t, "eoscanadacom", gjson.GetBytes(json, "struct_1_field_1").String())
	assert.Equal(t, "value_struct_2_field_1", gjson.GetBytes(json, "struct_2_field_1").String())
	assert.Equal(t, "value_struct_3_field_1", gjson.GetBytes(json, "struct_1_field_2.struct_3_field_1").String())
	assert.Equal(t, "value_struct_1_field_3", gjson.GetBytes(json, "struct_1_field_3").String())
	assert.Equal(t, "", gjson.GetBytes(json, "struct_1_field_4").String())
	assert.Equal(t, "value_struct_4_field_1_1", gjson.GetBytes(json, "struct_1_field_5.0.struct_4_field_1").String())
	assert.Equal(t, "value_struct_4_field_1_2", gjson.GetBytes(json, "struct_1_field_5.1.struct_4_field_1").String())
	assert.Equal(t, "value_struct_4_field_1_3", gjson.GetBytes(json, "struct_1_field_5.2.struct_4_field_1").String())

}

func TestABI_DecodeMissingData(t *testing.T) {
	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1 string
		F1  Name
	}{
		BF1: "value_struct_2_field_1",
		F1:  Name("eoscanadacom"),
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	require.NoError(t, err)

	abi, err := NewABI(abiReader)
	require.NoError(t, err)

	_, err = abi.DecodeAction(buffer.Bytes(), "action_name_1")
	assert.Equal(t, fmt.Errorf("decoding field struct_1_field_2: decoding [struct_name_3]: decoding field struct_3_field_1: read: varint: invalid buffer size").Error(), err.Error())

}

func TestABI_DecodeMissingAction(t *testing.T) {

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1 string
		F1  Name
	}{
		BF1: "value.base.field.1",
		F1:  Name("eoscanadacom"),
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	require.NoError(t, err)

	abi, err := NewABI(abiReader)
	require.NoError(t, err)

	_, err = abi.DecodeAction(buffer.Bytes(), "bad.action.name")
	assert.Equal(t, "action bad.action.name not found in abi", err.Error())
}

func TestABI_DecodeTable(t *testing.T) {

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1    string
		F1     Name
		F2     string
		F3FLAG byte //this a hack until we have the abi encoder
		F3     string
		F4FLAG byte //this a hack until we have the abi encoder
		F5     []string
	}{
		BF1:    "value_struct_2_field_1",
		F1:     Name("eoscanadacom"),
		F2:     "value_struct_3_field_1",
		F3FLAG: 1,
		F3:     "value_struct_1_field_3",
		F4FLAG: 0,
		F5:     []string{"value_struct_4_field_1_1", "value_struct_4_field_1_2", "value_struct_4_field_1_3"},
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	require.NoError(t, err)

	abi, err := NewABI(abiReader)
	require.NoError(t, err)

	json, err := abi.DecodeTableRow("table_name_1", buffer.Bytes())
	require.NoError(t, err)

	assert.Equal(t, "eoscanadacom", gjson.GetBytes(json, "struct_1_field_1").String())
	assert.Equal(t, "value_struct_2_field_1", gjson.GetBytes(json, "struct_2_field_1").String())
	assert.Equal(t, "value_struct_3_field_1", gjson.GetBytes(json, "struct_1_field_2.struct_3_field_1").String())
	assert.Equal(t, "value_struct_1_field_3", gjson.GetBytes(json, "struct_1_field_3").String())
	assert.Equal(t, "", gjson.GetBytes(json, "struct_1_field_4").String())
	assert.Equal(t, "value_struct_4_field_1_1", gjson.GetBytes(json, "struct_1_field_5.0.struct_4_field_1").String())
	assert.Equal(t, "value_struct_4_field_1_2", gjson.GetBytes(json, "struct_1_field_5.1.struct_4_field_1").String())
	assert.Equal(t, "value_struct_4_field_1_3", gjson.GetBytes(json, "struct_1_field_5.2.struct_4_field_1").String())

}

func TestABI_DecodeTableRowMissingTable(t *testing.T) {

	abiReader := strings.NewReader(abiString)

	mockData := struct {
		BF1 string
		F1  Name
	}{
		BF1: "value.base.field.1",
		F1:  Name("eoscanadacom"),
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(mockData)
	require.NoError(t, err)

	abi, err := NewABI(abiReader)
	require.NoError(t, err)

	_, err = abi.DecodeTableRow("bad.action.name", buffer.Bytes())
	assert.Equal(t, fmt.Errorf("table name bad.action.name not found in abi"), err)
}

func TestABI_DecodeBadABI(t *testing.T) {

	abiReader := strings.NewReader("{")
	_, err := NewABI(abiReader)
	assert.Equal(t, fmt.Errorf("read abi: unexpected EOF").Error(), err.Error())
}

func TestABI_decode(t *testing.T) {
	abi := &ABI{
		Structs: []StructDef{
			{
				Name: "struct.base.1",
				Fields: []FieldDef{
					{Name: "basefield1", Type: "string"},
				},
			},
			{
				Name: "struct.1",
				Base: "struct.base.1",
				Fields: []FieldDef{
					{Name: "field1", Type: "string"},
				},
			},
		},
	}

	s := struct {
		BF1 string
		F1  string
	}{
		BF1: "value.base.field.1",
		F1:  "value.field.1",
	}
	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(s)
	require.NoError(t, err)

	encodedTwoStrings := buffer.Bytes()
	//fmt.Println("EY", hex.EncodeToString(encodedTwoStrings))

	out, err := abi.Decode(NewDecoder(encodedTwoStrings), "struct.1")
	require.NoError(t, err)

	//fmt.Println("HOO", out, string(out))

	// FIXME: woah, the previous test was relying on the `FieldDef.Name`'s
	// value of `field.1` to expect `gjson` to fetch a nested field (!)

	// What was that?
	//var into interface{}
	//require.NoError(t, json.Unmarshal(out, &into))
	//spew.Dump(into)

	assert.Equal(t, "value.field.1", gjson.GetBytes(out, "field1").String())
	assert.Equal(t, "value.base.field.1", gjson.GetBytes(out, "basefield1").String())
}

func TestABI_decode_Float32FitNodeos(t *testing.T) {

	abi := &ABI{
		fitNodeos: true,
		Types:     []ABIType{},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "field", Type: "float32"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("cdcc8c3f")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"field":"1.10000002384185791"}`, string(json))
}

func TestABI_decode_StructFieldTypeUint128(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "extern_amount", Type: "uint128"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("0000f4b7062e7059ee11000000000000")
	require.NoError(t, err)

	abi.fitNodeos = false
	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"extern_amount":"0x0000f4b7062e7059ee11000000000000"}`, string(json))

	abi.fitNodeos = true
	json, err = abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"extern_amount":"84677000000000000000000"}`, string(json))
}

func TestABI_decode_StructFieldTypeTimePoint(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "timestamp", Type: "time_point"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("0020deedd5920500")
	require.NoError(t, err)

	abi.fitNodeos = true
	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"timestamp":"2019-09-18T16:00:00.000"}`, string(json))

	abi.fitNodeos = false
	json, err = abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"timestamp":"2019-09-18T16:00:00"}`, string(json))

}

func TestABI_decode_StructHasAliasedBase(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "base", NewTypeName: "aliasbase"},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Base: "aliasbase",
			},
			{
				Name: "base",
				Fields: []FieldDef{
					{Name: "data", Type: "uint8[]"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("02010a")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"data":[1,10]}`, string(json))
}

func TestABI_decode_StructFieldTypeHasAliasedBase(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "base", NewTypeName: "aliasbase"},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "item", Type: "item"},
				},
			},
			{
				Name: "base",
				Fields: []FieldDef{
					{Name: "data", Type: "uint8[]"},
				},
			},
			{
				Name: "item",
				Base: "aliasbase",
				Fields: []FieldDef{
					{Name: "name", Type: "name"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("02010a000000000010aa91")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"item":{"name":"map1","data":[1,10]}}`, string(json))
}

func TestABI_decode_StructFieldTypeHasArrayOfAliasArray(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "string[]", NewTypeName: "vec_string"},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "player_hands_string", Type: "vec_string[]"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("0104a7024469616d6f6e6420337c643634653132386664313837636330346139343935653065336535633230633938326330353731666239363932386137383033633931366139303232313133662c303666636565323034363465376138616333323130643034316533353638636434653430393437643766373431633763323432636361313432356635363265312c454f5337476151426267597769314d69676237786b417638385370767a76735363783242766f6d384a3236593273556d467070335a2c5349475f4b315f4b3666376b78327563635a75425950734d4c564a5a71786d4665634c4276796956677838336d6448567139467973726642634c6638674a627037525478334d45577246724a547061636762376b374b684e52644d73386f734e3572516f73a7024469616d6f6e6420417c326464323464613764383163393664353333383236323338343934336433663530316163323838316637633633376562616564663332653239653339346535312c303666636565323230633537643535623536386465313337623838376538653964303432316533313430303031623130643466396532663831393064653030322c454f5337476151426267597769314d69676237786b417638385370767a76735363783242766f6d384a3236593273556d467070335a2c5349475f4b315f4b32626d3636676d636f3553516b6e476364624c4d4d7838464d3334695832774e7543674e4443377254323870317a336765486d7057503771397547336a72706f386a6871474d33596970576162477237584131716a394a545a4a376f47a502486561727420337c373530636631666139366337646564306463316336373536626162323661356564316630353262373062373036633237396666303638656639306136333136332c303666636565336662363866343830386662313431343439653536383866633830613837386566326331663036373031643733383362363535666534633666362c454f5337476151426267597769314d69676237786b417638385370767a76735363783242766f6d384a3236593273556d467070335a2c5349475f4b315f4a76484a5877663361664652364a6879375a415a57744b576e736d5a4e6a5753556e635452386a696e38557a624c486a4759764c714d38373166637056326a73386953735658324d6a6e4d37426e594c37536739777a526f48746d716f35a502537061646520347c613766613637353565366632373735343635333432663533633733393337313638383432386333316130663066666634616566313764343163626232613766342c303666636565366135633366343835363930303864383338376666663538363332626164396334663935653832376531333238356563633862333337366633342c454f5337476151426267597769314d69676237786b417638385370767a76735363783242766f6d384a3236593273556d467070335a2c5349475f4b315f4b5a644e337055583235574254387879574e3478553936356e594c72734c4857613351383931486169484578744b7253594b31623636377679376d52633166565148474471763268316248313733646b76446a7552364e36503733554644")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"player_hands_string":[["Diamond 3|d64e128fd187cc04a9495e0e3e5c20c982c0571fb96928a7803c916a9022113f,06fcee20464e7a8ac3210d041e3568cd4e40947d7f741c7c242cca1425f562e1,EOS7GaQBbgYwi1Migb7xkAv88SpvzvsScx2Bvom8J26Y2sUmFpp3Z,SIG_K1_K6f7kx2uccZuBYPsMLVJZqxmFecLBvyiVgx83mdHVq9FysrfBcLf8gJbp7RTx3MEWrFrJTpacgb7k7KhNRdMs8osN5rQos","Diamond A|2dd24da7d81c96d5338262384943d3f501ac2881f7c637ebaedf32e29e394e51,06fcee220c57d55b568de137b887e8e9d0421e3140001b10d4f9e2f8190de002,EOS7GaQBbgYwi1Migb7xkAv88SpvzvsScx2Bvom8J26Y2sUmFpp3Z,SIG_K1_K2bm66gmco5SQknGcdbLMMx8FM34iX2wNuCgNDC7rT28p1z3geHmpWP7q9uG3jrpo8jhqGM3YipWabGr7XA1qj9JTZJ7oG","Heart 3|750cf1fa96c7ded0dc1c6756bab26a5ed1f052b70b706c279ff068ef90a63163,06fcee3fb68f4808fb141449e5688fc80a878ef2c1f06701d7383b655fe4c6f6,EOS7GaQBbgYwi1Migb7xkAv88SpvzvsScx2Bvom8J26Y2sUmFpp3Z,SIG_K1_JvHJXwf3afFR6Jhy7ZAZWtKWnsmZNjWSUncTR8jin8UzbLHjGYvLqM871fcpV2js8iSsVX2MjnM7BnYL7Sg9wzRoHtmqo5","Spade 4|a7fa6755e6f2775465342f53c739371688428c31a0f0fff4aef17d41cbb2a7f4,06fcee6a5c3f48569008d8387fff58632bad9c4f95e827e13285ecc8b3376f34,EOS7GaQBbgYwi1Migb7xkAv88SpvzvsScx2Bvom8J26Y2sUmFpp3Z,SIG_K1_KZdN3pUX25WBT8xyWN4xU965nYLrsLHWa3Q891HaiHExtKrSYK1b667vy7mRc1fVQHGDqv2h1bH173dkvDjuR6N6P73UFD"]]}`, string(json))
}

func TestABI_decode_StructFieldWithUint128(t *testing.T) {
	abi := &ABI{
		fitNodeos: true,
		Types:     []ABIType{},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "extern_amount", Type: "uint128"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("9ea6ce00000000000000000000000000")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"extern_amount": "13543070"}`, string(json))
}

func TestABI_decode_StructFieldTypeHasAlias(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "uint8", NewTypeName: "alias"},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "item", Type: "alias"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("0a")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"item":10}`, string(json))
}

func TestABI_decode_StructFieldArrayTypeHasAlias(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "uint8", NewTypeName: "alias"},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "item", Type: "alias[]"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("02010a")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"item":[1,10]}`, string(json))
}

func TestABI_decode_StructFieldTypeHasAliasArray(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "uint8[]", NewTypeName: "alias"},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "item", Type: "alias"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("02010a")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"item":[1,10]}`, string(json))
}

func TestABI_decode_StructFieldHasAliasWithStructType(t *testing.T) {

	abi := &ABI{
		fitNodeos: true,
		Types: []ABIType{
			{Type: "collab_data[]", NewTypeName: "approvals_t"},
		},
		Structs: []StructDef{
			{
				Name: "struct_with_alias_with_struct_type",
				Fields: []FieldDef{
					{Name: "requested", Type: "approvals_t"},
				},
			},
			{
				Name: "collab_data",
				Fields: []FieldDef{
					{Name: "asset_owner", Type: "name"},
					{Name: "asset_id", Type: "uint64"},
					{Name: "asset_uid", Type: "uint64"},
					{Name: "percentage", Type: "uint64"},
					{Name: "accepted", Type: "bool"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("0190316d4c65338d54510b000000000000f1477db8479a050040420f000000000001")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "struct_with_alias_with_struct_type")

	require.NoError(t, err)
	assert.JSONEq(t, `{"requested": [{"asset_owner": "emanateghost", "accepted": 1, "asset_id": 2897, "asset_uid": "1577007712126961", "percentage": 1000000}]}`, string(json))
}

func TestABI_decode_StructFieldHasAlias(t *testing.T) {

	abi := &ABI{
		Types: []ABIType{
			{Type: "name", NewTypeName: "alias"},
		},
		Structs: []StructDef{
			{
				Name: "struct_with_alias",
				Fields: []FieldDef{
					{Name: "item", Type: "alias"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("000000000010aa91")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "struct_with_alias")
	require.NoError(t, err)

	assert.JSONEq(t, `{"item":"map1"}`, string(json))
}

func TestABI_decode_StructFieldTypeArray_ResetsOuterFields(t *testing.T) {
	abi := &ABI{
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "name", Type: "name"},
					{Name: "items", Type: "item[]"},
				},
			},
			{
				Name: "item",
				Fields: []FieldDef{
					{Name: "data", Type: "uint8[]"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("000000000010aa910102010a")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"items":[{"data":[1,10]}],"name":"map1"}`, string(json))
}

func TestABI_decode_StructNotFound(t *testing.T) {
	abi := &ABI{
		Structs: []StructDef{},
	}

	_, err := abi.Decode(NewDecoder(nil), "struct.1")
	assert.Equal(t, fmt.Errorf("structure [struct.1] not found in abi"), err)
}

func TestABI_decode_StructBaseNotFound(t *testing.T) {
	abi := &ABI{
		Structs: []StructDef{
			{
				Name: "struct.1",
				Base: "struct.base.1",
				Fields: []FieldDef{
					{Name: "field.1", Type: "name"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("202932c94c833055")
	require.NoError(t, err)

	_, err = abi.Decode(NewDecoder(buffer), "struct.1")
	assert.Equal(t, "decode base [struct.1]: structure [struct.base.1] not found in abi", err.Error())
}

func TestABI_decode_StructFieldArrayType_HasSjsonPathLikeName(t *testing.T) {
	t.Skipf("This test exhibits the problem where the field name has sjson array like path, the expectOutput is completely wrong.")

	abi := &ABI{
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "item.1", Type: "uint8[]"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("02010a")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"item.1":[1,10]}`, string(json))
}

func TestABI_decode_StructVariantField(t *testing.T) {
	abi := &ABI{
		Variants: []VariantDef{
			{
				Name:  "variant_",
				Types: []string{"name", "uint32"},
			},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "field", Type: "variant_"},
				},
			},
		},
	}
	buffer, err := hex.DecodeString("00000050df45e3aec2")
	require.NoError(t, err)
	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)
	assert.JSONEq(t, `{"field": ["name", "serialize"]}`, string(json))
	buffer, err = hex.DecodeString("0164000000")
	require.NoError(t, err)
	json, err = abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)
	assert.JSONEq(t, `{"field":["uint32", 100]}`, string(json))
}

func TestABI_decode_StructArrayOfVariantField_OneOfVariantIsAlias(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "name", NewTypeName: "my_name"},
		},
		Variants: []VariantDef{
			{
				Name:  "variant_",
				Types: []string{"my_name", "uint32"},
			},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "field", Type: "variant_[]"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("0200000050df45e3aec20164000000")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"field":[["name","serialize"],["uint32",100]]}`, string(json))
}

func TestABI_decode_Struct2DArrayOfVariantField_OneOfVariantIsAlias(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "name", NewTypeName: "my_name"},
		},
		Variants: []VariantDef{
			{
				Name:  "variant_",
				Types: []string{"my_name", "uint32"},
			},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "field", Type: "variant_[][]"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("010200000050df45e3aec20164000000")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"field":[[["name","serialize"],["uint32",100]]]}`, string(json))
}

func TestABI_decode_StructVariantField_OneOfVariantIsAlias(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "name", NewTypeName: "my_name"},
		},
		Variants: []VariantDef{
			{
				Name:  "variant_",
				Types: []string{"my_name", "uint32"},
			},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "field", Type: "variant_"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("00000050df45e3aec2")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"field":["name","serialize"]}`, string(json))

	buffer, err = hex.DecodeString("0164000000")
	require.NoError(t, err)

	json, err = abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"field":["uint32",100]}`, string(json))
}

func TestABI_decode_StructAliasToAVariantField(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{
			{Type: "variant_", NewTypeName: "my_variant"},
		},
		Variants: []VariantDef{
			{
				Name:  "variant_",
				Types: []string{"name", "uint32"},
			},
		},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "field", Type: "my_variant"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("00000050df45e3aec2")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"field":["name","serialize"]}`, string(json))

	buffer, err = hex.DecodeString("0164000000")
	require.NoError(t, err)

	json, err = abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"field":["uint32",100]}`, string(json))
}

func TestABI_decode_Uint8ArrayVec(t *testing.T) {

	abi := &ABI{
		Types: []ABIType{
			{Type: "name", NewTypeName: "alias"},
		},
		Structs: []StructDef{
			{
				Name: "endgame",
				Fields: []FieldDef{
					{Name: "player_hands", Type: "uint8[][]"},
				},
			},
		},
	}

	json, err := abi.Decode(NewDecoder(HexString("01020d13")), "endgame")
	require.NoError(t, err)
	assert.JSONEq(t, `{"player_hands": [[13,19]]}`, string(json))
}

func TestABI_decode_BinaryExtension(t *testing.T) {
	abi := &ABI{
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "id", Type: "uint8"},
					{Name: "name", Type: "name$"},
				},
			},
		},
	}

	json, err := abi.Decode(NewDecoder(HexString("00")), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"id":0}`, string(json))

	json, err = abi.Decode(NewDecoder(HexString("000000000000a0a499")), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"id":0,"name":"name"}`, string(json))
}

func TestABI_decode_BinaryExtensionArray(t *testing.T) {
	abi := &ABI{
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "id", Type: "uint8"},
					{Name: "name", Type: "name[]$"},
				},
			},
		},
	}

	json, err := abi.Decode(NewDecoder(HexString("00")), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"id":0}`, string(json))

	json, err = abi.Decode(NewDecoder(HexString("00010000000000a0a499")), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"id":0,"name":["name"]}`, string(json))
}

func TestABI_decodeFields(t *testing.T) {
	types := []ABIType{
		{NewTypeName: "action.type.1", Type: "name"},
	}
	fields := []FieldDef{
		{Name: "F1", Type: "uint64"},
		{Name: "F2", Type: "action.type.1"},
	}
	abi := &ABI{
		Types: types,
		Structs: []StructDef{
			{Fields: fields},
		},
	}

	buffer, err := hex.DecodeString("ffffffffffffffff202932c94c833055")
	require.NoError(t, err)

	json, err := abi.decodeFields(NewDecoder(buffer), fields, map[string]interface{}{})
	require.NoError(t, err)

	assert.JSONEq(t, `{"F1":"18446744073709551615", "F2":"eoscanadacom"}`, toJSON(t, json))
}

func toJSON(t *testing.T, in interface{}) string {
	cnt, err := json.Marshal(in)
	require.NoError(t, err)
	return string(cnt)
}

func TestABI_decodeFieldsErr(t *testing.T) {

	types := []ABIType{}
	fields := []FieldDef{
		{
			Name: "field.with.bad.type.1",
			Type: "bad.type.1",
		},
	}

	s := struct{}{}

	abi := &ABI{
		Types: types,
		Structs: []StructDef{
			{Fields: fields},
		},
	}

	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	err := encoder.Encode(s)
	require.NoError(t, err)

	_, err = abi.decodeFields(NewDecoder(buffer.Bytes()), fields, map[string]interface{}{})
	assert.Equal(t, fmt.Errorf("decoding field field.with.bad.type.1: read field of type [bad.type.1]: unknown type").Error(), err.Error())

}

func TestABI_decodeOptionalField(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "field", Type: "string?"},
				},
			},
		},
	}

	optional := &testOptionalField{
		B: 1,
		S: "value.1",
	}

	optionalNotPresent := &testOptionalField{
		B: 0,
	}
	optionalMissingFlag := struct{}{}

	tests := []struct {
		name         string
		data         interface{}
		fitNodeos    bool
		expectError  bool
		expectedJSON string
	}{
		{
			name:         "optional present",
			data:         optional,
			expectedJSON: `{"field":"value.1"}`,
		},
		{
			name:         "optional not present",
			data:         optionalNotPresent,
			expectedJSON: `{}`,
		},
		{
			name:         "optional not present fit nodeos",
			fitNodeos:    true,
			data:         optionalNotPresent,
			expectedJSON: `{"field": null}`,
		},
		{
			name:        "optional missing flag",
			data:        optionalMissingFlag,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var buffer bytes.Buffer
			encoder := NewEncoder(&buffer)
			err := encoder.Encode(test.data)
			require.NoError(t, err)

			abi.fitNodeos = test.fitNodeos
			json, err := abi.Decode(NewDecoder(buffer.Bytes()), "root")
			if test.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if test.expectedJSON != "" {
					assert.JSONEq(t, test.expectedJSON, string(json))
				} else {
					assert.Equal(t, test.expectedJSON, string(json))
				}

			}
		})
	}

	//// Option is present
	//// Option is not present
	//
	//buffer, err := hex.DecodeString("000000000010aa91")
	//require.NoError(t, err)
	//
	//json, err := abi.Decode(NewDecoder(buffer), "struct_with_alias")
	//require.NoError(t, err)
	//
	//assert.JSONEq(t, `{"item":"map1"}`, string(json))
}

func TestABI_decode_StructFieldArray(t *testing.T) {
	abi := &ABI{
		Types: []ABIType{},
		Structs: []StructDef{
			{
				Name: "root",
				Fields: []FieldDef{
					{Name: "item", Type: "string[]"},
				},
			},
		},
	}

	buffer, err := hex.DecodeString("00")
	require.NoError(t, err)

	json, err := abi.Decode(NewDecoder(buffer), "root")
	require.NoError(t, err)

	assert.JSONEq(t, `{"item":[]}`, string(json))
}

func TestABI_Read(t *testing.T) {
	someTime, err := time.Parse("2006-01-02T15:04:05", "2018-09-05T12:48:54")
	require.NoError(t, err)
	bt := BlockTimestamp{
		Time: someTime,
	}
	require.NoError(t, err)

	signatureBuffer, err := hex.DecodeString("001f69c3e7b2789bccd5c4be1030129f35e93de2e8e18468dca94c65600cac25b4636e5d75342499e5519a0df74c714fd5ad682662204068eff4ca9fac86254ae416")
	require.NoError(t, err)

	tests := []struct {
		name          string
		typeName      string
		fitNodeos     bool
		data          interface{}
		expectError   bool
		expectedValue *string
	}{
		{name: "string", typeName: "string", data: "this.is.a.test", expectedValue: s(`"this.is.a.test"`)},
		{name: "min int8", typeName: "int8", data: int8(-128), expectedValue: s("-128")},
		{name: "max int8", typeName: "int8", data: int8(127), expectedValue: s("127")},
		{name: "min uint8", typeName: "uint8", data: uint8(0), expectedValue: s("0")},
		{name: "max uint8", typeName: "uint8", data: uint8(255), expectedValue: s("255")},
		{name: "min int16", typeName: "int16", data: int16(-32768), expectedValue: s("-32768")},
		{name: "max int16", typeName: "int16", data: int16(32767), expectedValue: s("32767")},
		{name: "min uint16", typeName: "uint16", data: uint16(0), expectedValue: s("0")},
		{name: "max uint16", typeName: "uint16", data: uint16(65535), expectedValue: s("65535")},
		{name: "min int32", typeName: "int32", data: int32(-2147483648), expectedValue: s("-2147483648")},
		{name: "max int32", typeName: "int32", data: int32(2147483647), expectedValue: s("2147483647")},
		{name: "min uint32", typeName: "uint32", data: uint32(0), expectedValue: s("0")},
		{name: "max uint32", typeName: "uint32", data: uint32(4294967295), expectedValue: s("4294967295")},
		{name: "min int64", typeName: "int64", data: int64(-9223372036854775808), expectedValue: s(`"-9223372036854775808"`)},
		{name: "max int64", typeName: "int64", data: int64(9223372036854775807), expectedValue: s(`"9223372036854775807"`)},
		{name: "mid int64", typeName: "int64", data: int64(4096), expectedValue: s(`4096`)},
		{name: "stringified lower int64", typeName: "int64", data: int64(-5000000000), expectedValue: s(`"-5000000000"`)},
		{name: "min uint64", typeName: "uint64", data: uint64(0), expectedValue: s("0")},
		{name: "max uint64", typeName: "uint64", data: uint64(18446744073709551615), expectedValue: s(`"18446744073709551615"`)},
		{name: "int128 1", typeName: "int128", data: Int128{Lo: 1, Hi: 0}, expectedValue: s(`"0x01000000000000000000000000000000"`)},
		{name: "int128 -1", typeName: "int128", data: Int128{Lo: math.MaxUint64, Hi: math.MaxUint64}, expectedValue: s(`"0xffffffffffffffffffffffffffffffff"`)},
		{name: "int128", typeName: "int128", data: Int128{Lo: 925, Hi: 125}, expectedValue: s(`"0x9d030000000000007d00000000000000"`)},
		{name: "int128 negative ", typeName: "int128", data: Int128{Lo: 1, Hi: math.MaxUint64}, expectedValue: s(`"0x0100000000000000ffffffffffffffff"`)},
		{name: "int128 fit nodeos", fitNodeos: true, typeName: "int128", data: Int128{Lo: 925, Hi: 125}, expectedValue: s(`"2305843009213693952925"`)},
		{name: "int128 negative fit nodeos ", fitNodeos: true, typeName: "int128", data: Int128{Lo: 1, Hi: math.MaxUint64}, expectedValue: s(`"-18446744073709551615"`)},
		{name: "uint128 1", typeName: "uint128", data: Int128{Lo: 1, Hi: 0}, expectedValue: s(`"0x01000000000000000000000000000000"`)},
		{name: "uint128", typeName: "uint128", data: Uint128{Lo: 925, Hi: 125}, expectedValue: s(`"0x9d030000000000007d00000000000000"`)},
		{name: "uint128 fit nodeos", fitNodeos: true, typeName: "uint128", data: Uint128{Lo: 925, Hi: 125}, expectedValue: s(`"2305843009213693952925"`)},
		{name: "uint128 max", typeName: "uint128", data: Uint128{Lo: math.MaxUint64, Hi: math.MaxUint64}, expectedValue: s(`"0xffffffffffffffffffffffffffffffff"`)},
		{name: "uint128 fit nodeos", fitNodeos: true, typeName: "uint128", data: Uint128{Lo: math.MaxUint64, Hi: math.MaxUint64}, expectedValue: s(`"340282366920938463463374607431768211455"`)},
		{name: "min varint32", typeName: "varint32", data: Varint32(-2147483648), expectedValue: s("-2147483648")},
		{name: "max varint32", typeName: "varint32", data: Varint32(2147483647), expectedValue: s("2147483647")},
		{name: "min varuint32", typeName: "varuint32", data: Varuint32(0), expectedValue: s("0")},
		{name: "max varuint32", typeName: "varuint32", data: Varuint32(4294967295), expectedValue: s("4294967295")},
		{name: "min float 32", typeName: "float32", data: float32(math.SmallestNonzeroFloat32), expectedValue: s("0.000000000000000000000000000000000000000000001401298464324817")},
		{name: "max float 32", typeName: "float32", data: float32(math.MaxFloat32), expectedValue: s("340282346638528860000000000000000000000")},
		{name: "min float64", typeName: "float64", data: math.SmallestNonzeroFloat64, expectedValue: s("0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005")},
		{name: "max float64", typeName: "float64", data: math.MaxFloat64, expectedValue: s("179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")},
		{name: "float128", typeName: "float128", data: Float128{Lo: 1, Hi: 2}, expectedValue: s(`"0x01000000000000000200000000000000"`)},
		{name: "bool true", typeName: "bool", data: true, expectedValue: s("true")},
		{name: "bool false", typeName: "bool", data: false, expectedValue: s("false")},
		{name: "time_point", typeName: "time_point", data: TimePoint(1541085187001001), expectedValue: s(`"2018-11-01T15:13:07.001"`)},
		{name: "time_point_sec", typeName: "time_point_sec", data: TimePointSec(1681469753), expectedValue: s(`"2023-04-14T10:55:53"`)},
		{name: "block_timestamp_type", typeName: "block_timestamp_type", data: bt, expectedValue: s(`"2018-09-05T12:48:54"`)},
		{name: "Name", typeName: "name", data: Name("eoscanadacom"), expectedValue: s(`"eoscanadacom"`)},
		{name: "bytes", typeName: "bytes", data: []byte("this.is.a.test"), expectedValue: s(`"746869732e69732e612e74657374"`)},
		{name: "checksum160", typeName: "checksum160", data: Checksum160(make([]byte, TypeSize.Checksum160)), expectedValue: s(`"0000000000000000000000000000000000000000"`)},
		{name: "checksum256", typeName: "checksum256", data: Checksum256(make([]byte, TypeSize.Checksum256)), expectedValue: s(`"0000000000000000000000000000000000000000000000000000000000000000"`)},
		{name: "checksum512", typeName: "checksum512", data: Checksum512(make([]byte, TypeSize.Checksum512)), expectedValue: s(`"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`)},
		{name: "public_key", typeName: "public_key", data: ecc.MustNewPublicKey(ecc.PublicKeyPrefixCompat + "1111111111111111111111111111111114T1Anm"), expectedValue: s(`"` + ecc.PublicKeyPrefixCompat + `1111111111111111111111111111111114T1Anm"`)},
		{name: "public_key_k1", typeName: "public_key", data: ecc.MustNewPublicKey("PUB_K1_1111111111111111111111111111111114T1Anm"), expectedValue: s(`"` + ecc.PublicKeyPrefixCompat + `1111111111111111111111111111111114T1Anm"`)},
		{name: "public_key_wa", typeName: "public_key", data: ecc.MustNewPublicKey("PUB_WA_5hyixc7vkMbKiThWi1TnFtXw7HTDcHfjREj2SzxCtgw3jQGepa5T9VHEy1Tunjzzj"), expectedValue: s(`"PUB_WA_5hyixc7vkMbKiThWi1TnFtXw7HTDcHfjREj2SzxCtgw3jQGepa5T9VHEy1Tunjzzj"`)},
		{name: "signature", typeName: "signature", data: ecc.MustNewSignatureFromData(signatureBuffer), expectedValue: s(`"SIG_K1_K96L1au4xFJg5edn6qBK6UDbSsC2RKsMs4cXCA2LoCPZxBDMXehdZFWPh1GeRhzGoQjBwNK2eBmUXf4L8SBApL69pGdUJm"`)},
		{name: "signature_wa", typeName: "signature", data: ecc.MustNewSignature("SIG_WA_28AzYsRYSSA85Q4Jjp4zkiyBA8G85AcPsHU3HUuqLkY3LooYcFiSMGGxhEQcCzAhaZJqdaUXG16p8t63sDhqh9L4xc24CDxbf81D6FW4SXGjxQSM2D7FAJSSQCogjbqJanTP5CbSF8FWyaD4pVVAs4Z9ubqNhHCkiLDesEukwGYu6ujgwQkFqczow5cSwTqTirdgqCBjkGQLMT3KV2JwjN7b2qPAyDa2vvjsGWFP8HVTw2tctD6FBPHU9nFgtfcztkc3eqxVU9UbvUbKayU62dLZBwNCwHxmyPymH5YfoJLhBkS8s"), expectedValue: s(`"SIG_WA_28AzYsRYSSA85Q4Jjp4zkiyBA8G85AcPsHU3HUuqLkY3LooYcFiSMGGxhEQcCzAhaZJqdaUXG16p8t63sDhqh9L4xc24CDxbf81D6FW4SXGjxQSM2D7FAJSSQCogjbqJanTP5CbSF8FWyaD4pVVAs4Z9ubqNhHCkiLDesEukwGYu6ujgwQkFqczow5cSwTqTirdgqCBjkGQLMT3KV2JwjN7b2qPAyDa2vvjsGWFP8HVTw2tctD6FBPHU9nFgtfcztkc3eqxVU9UbvUbKayU62dLZBwNCwHxmyPymH5YfoJLhBkS8s"`)},
		{name: "symbol", typeName: "symbol", data: EOSSymbol, expectedValue: s(`"4,EOS"`)},
		{name: "symbol_code", typeName: "symbol_code", data: SymbolCode(22606239386324546), expectedValue: s(`"BNTDAPP"`)},
		{name: "asset", typeName: "asset", data: Asset{Amount: 100000, Symbol: EOSSymbol}, expectedValue: s(`"10.0000 EOS"`)},
		{name: "extended_asset", typeName: "extended_asset", data: ExtendedAsset{Asset: Asset{Amount: 10, Symbol: EOSSymbol}, Contract: "eoscanadacom"}, expectedValue: s("{\"quantity\":\"0.0010 EOS\",\"contract\":\"eoscanadacom\"}")},
		{name: "bad type", typeName: "bad.type.1", data: nil, expectedValue: nil, expectError: true},
		{name: "min float64 fit nodeos", fitNodeos: true, typeName: "float64", data: math.SmallestNonzeroFloat64, expectedValue: s(`"0.00000000000000000"`)},
		{name: "max float64 fit nodeos", fitNodeos: true, typeName: "float64", data: math.MaxFloat64, expectedValue: s(`"179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.00000000000000000"`)},
		{name: "bool true fit nodeos", fitNodeos: true, typeName: "bool", data: true, expectedValue: s("1")},
		{name: "bool false  fit nodeos", fitNodeos: true, typeName: "bool", data: false, expectedValue: s("0")},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			var buffer bytes.Buffer
			encoder := NewEncoder(&buffer)
			err := encoder.Encode(c.data)
			require.NoError(t, err)

			value := ""
			if c.expectedValue != nil {
				value = *c.expectedValue
			}
			require.NoError(t, err, fmt.Sprintf("encoding value %s, of type %s", value, c.typeName), c.name)

			abi := ABI{}
			abi.fitNodeos = c.fitNodeos
			json, err := abi.read(NewDecoder(buffer.Bytes()), c.typeName)
			if c.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, value, toJSON(t, json), c.name)
			}
		})
	}
}

func TestABI_Read_Symbol(t *testing.T) {
	abi := ABI{}
	data, err := hex.DecodeString("04454f5300000000")
	require.NoError(t, err)

	out, err := abi.read(NewDecoder(data), "symbol")
	require.NoError(t, err)
	assert.Equal(t, `"4,EOS"`, toJSON(t, out))
}

func TestABI_Read_TimePointSec(t *testing.T) {
	abi := ABI{}
	data, err := hex.DecodeString("919dd85b")
	require.NoError(t, err)
	out, err := abi.read(NewDecoder(data), "time_point_sec")
	//out, err := abi.read(NewDecoder([]byte("c15dd35b")), "name", "time_point_sec", []byte("{}"))
	//out, err := abi.read(NewDecoder([]byte("919dd85b")), "name", "time_point_sec", []byte("{}"))
	require.NoError(t, err)
	assert.Equal(t, `"2018-10-30T18:06:09"`, toJSON(t, out))
}

func TestABI_Read_SymbolCode(t *testing.T) {
	abi := ABI{}
	data, err := hex.DecodeString("424e544441505000")
	require.NoError(t, err)

	out, err := abi.read(NewDecoder(data), "symbol_code")
	require.NoError(t, err)
	assert.Equal(t, `"BNTDAPP"`, toJSON(t, out))
}

func TestABIDecoder_analyseFieldType(t *testing.T) {
	testCases := []struct {
		fieldType               string
		expectedName            string
		expectedOptional        bool
		expectedArray           bool
		expectedBinaryExtension bool
	}{
		{"field.type.1", "field.type.1", false, false, false},
		{"field.type.1?", "field.type.1", true, false, false},
		{"field.type.1[]", "field.type.1", false, true, false},
		{"field.type.1$", "field.type.1", false, false, true},
		//ultra-andrey-bezrukov --- BLOCK-178 Dfuse cannot produce JSON data for migration
		{"field.type.1?$", "field.type.1", true, true, true},
	}

	for i, test := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			name, isOptional, isArray, isBinaryExtension := analyzeFieldType(test.fieldType)
			assert.Equal(t, test.expectedName, name)
			assert.Equal(t, test.expectedOptional, isOptional)
			assert.Equal(t, test.expectedArray, isArray)
			assert.Equal(t, test.expectedBinaryExtension, isBinaryExtension)
		})
	}
}

func Test_formatTimePoint(t *testing.T) {
	tests := []struct {
		name            string
		time            TimePoint
		shouldFitNodeos bool
		expectedOutput  string
	}{
		{
			name:            "golden path with fit nodeos",
			time:            1588450213523000,
			shouldFitNodeos: false,
			expectedOutput:  "2020-05-02T20:10:13.523",
		},
		{
			name:            "golden path without fit nodeos",
			time:            1588450213523000,
			shouldFitNodeos: false,
			expectedOutput:  "2020-05-02T20:10:13.523",
		},
		{
			name:            "0 nano second with fit nodeos",
			time:            1568822400000000,
			shouldFitNodeos: true,
			expectedOutput:  "2019-09-18T16:00:00.000",
		},
		{
			name:            "0 nano second without fit nodeos",
			time:            1568822400000000,
			shouldFitNodeos: false,
			expectedOutput:  "2019-09-18T16:00:00",
		},
		{
			name:            "500 nano second with fit nodeos",
			time:            1588450213500000,
			shouldFitNodeos: true,
			expectedOutput:  "2020-05-02T20:10:13.500",
		},
		{
			name:            "500 nano second without fit nodeos",
			time:            1588450213500000,
			shouldFitNodeos: false,
			expectedOutput:  "2020-05-02T20:10:13.5",
		},
		{
			name:            "520 nano second with fit nodeos",
			time:            1588450213520000,
			shouldFitNodeos: true,
			expectedOutput:  "2020-05-02T20:10:13.520",
		},
		{
			name:            "520 nano second without fit nodeos",
			time:            1588450213520000,
			shouldFitNodeos: false,
			expectedOutput:  "2020-05-02T20:10:13.52",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			time := formatTimePoint(test.time, test.shouldFitNodeos)
			//fmt.Println(time)
			assert.Equal(t, test.expectedOutput, time)
		})
	}

}

func HexString(input string) []byte {
	buffer, err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}

	return buffer
}

func s(str string) *string {
	return &str
}

type testOptionalField struct {
	B byte
	S string
}
