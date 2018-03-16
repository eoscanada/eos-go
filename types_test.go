package eosapi

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimplePacking(t *testing.T) {

	type S struct {
		P string
	}
	type M struct {
		Acct AccountName
		A    []*S
	}
	cnt, err := MarshalBinary(&M{
		Acct: AccountName("."),
		A:    []*S{},
	})

	// type M struct {
	// 	NumA Varint `struc:"sizeof=A"`
	// 	A    []string
	// }

	// var buf bytes.Buffer
	// err := struc.Pack(&buf, &M{
	// 	A: []string{"hello", "world"},
	// })
	require.NoError(t, err)
	assert.Equal(t, `000000`, hex.EncodeToString(cnt))
}

func TestUnpackBinaryTableRows(t *testing.T) {

	resp := &GetTableRowsResp{
		Rows: json.RawMessage(`["044355520000000004435552000000000000000000000000"]`),
	}
	assert.NoError(t, resp.BinaryToStructs(true))
}

func TestStringToName(t *testing.T) {
	i, err := StringToName("tbcox2.3")
	require.NoError(t, err)
	assert.Equal(t, uint64(0xc9d14e8803000000), i)
}

func TestNameToString(t *testing.T) {
	assert.Equal(t, "tbcox2.3", NameToString(uint64(0xc9d14e8803000000)))
}

func TestPackAccountName(t *testing.T) {
	// SHOULD IMPLEMENT: string_to_name from contracts/eosiolib/types.hpp
	tests := []struct {
		in  string
		out []byte
	}{
		{"eosio", []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55}},
		{"eosio.system", []byte{0x20, 0x55, 0xc6, 0x1e, 0x3, 0xea, 0x30, 0x55}},
		{"tbcox2.3", []byte{0x0, 0x0, 0x0, 0x3, 0x88, 0x4e, 0xd1, 0xc9}},
		{"tbcox2.", []byte{0x0, 0x0, 0x0, 0x0, 0x88, 0x4e, 0xd1, 0xc9}},
	}

	for idx, test := range tests {
		acct := AccountName(test.in)
		buf, err := MarshalBinary(acct)
		assert.NoError(t, err)
		assert.Equal(t, test.out, buf, fmt.Sprintf("index %d", idx))
	}
}

func TestUnpackActionTransfer(t *testing.T) {
	tests := []struct {
		in  string
		out Transfer
	}{
		{
			"00000003884ed1c900000000884ed1c9090000000000000000000000000000000000000000000000",
			Transfer{AccountName("tbcox2.3"), AccountName("tbcox2"), 9, ""},
		},
		{
			"00000003884ed1c900000000884ed1c9090000000000000000",
			Transfer{AccountName("tbcox2.3"), AccountName("tbcox2"), 9, ""},
		},
	}

	for _, test := range tests {
		buf, err := hex.DecodeString(test.in)
		assert.NoError(t, err)

		var res Transfer
		assert.NoError(t, UnmarshalBinary(buf, &res))
		assert.Equal(t, test.out, res)
	}

}

func TestActionMetaTypes(t *testing.T) {
	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("transfer"),
		Fields: &Transfer{
			From: AccountName("abourget"),
			To:   AccountName("mama"),
		},
	}

	cnt, err := json.Marshal(a)
	require.NoError(t, err)
	assert.Equal(t,
		`{"account":"eosio","from":"abourget","memo":"","name":"transfer","quantity":0,"to":"mama"}`,
		string(cnt),
	)

	var newAction *Action
	require.NoError(t, json.Unmarshal(cnt, &newAction))
}

func TestActionNoFields(t *testing.T) {
	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("transfer"),
	}

	cnt, err := json.Marshal(a)
	require.NoError(t, err)
	assert.Equal(t,
		`{"account":"eosio","name":"transfer"}`,
		string(cnt),
	)

	var newAction *Action
	require.NoError(t, json.Unmarshal(cnt, &newAction))
}

func TestHexBytes(t *testing.T) {
	a := HexBytes("hello world")
	cnt, err := json.Marshal(a)
	require.NoError(t, err)
	assert.Equal(t,
		`"68656c6c6f20776f726c64"`,
		string(cnt),
	)

	var b HexBytes
	require.NoError(t, json.Unmarshal(cnt, &b))
	assert.Equal(t, a, b)
}
