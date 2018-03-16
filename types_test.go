package eosapi

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
		Acct: AccountName("bob"),
		A:    []*S{&S{"hello"}, &S{"world"}},
	})

	require.NoError(t, err)
	assert.Equal(t, "0000000000000e3d020568656c6c6f05776f726c64", hex.EncodeToString(cnt))
}

func TestPackTransaction(t *testing.T) {
	stamp := time.Date(1970, time.September, 1, 1, 1, 1, 1, time.UTC)
	blockID, _ := hex.DecodeString("00106438d58d4fcab54cf89ca8308e5971cff735979d6050c6c1b45d8aadcad6")
	tx := &Transaction{
		RefBlockNum:    uint16(binary.LittleEndian.Uint64(blockID[:8])),
		RefBlockPrefix: uint32(binary.LittleEndian.Uint64(blockID[16:24])),
		Expiration:     JSONTime{stamp},
		Actions: []*Action{
			{
				Account: AccountName("eosio"),
				Name:    ActionName("transfer"),
				Authorization: []PermissionLevel{
					{AccountName("eosio"), PermissionName("active")},
				},
			},
		},
	}

	buf, err := MarshalBinary(tx)
	assert.NoError(t, err)
	// "cd6a4001  0000  0010  b54cf89c  0000  0000  00  01  0000000000ea3055  000000572d3ccdcd
	// 01
	// permission level: 0000000000ea3055  00000000a8ed3232
	// data: 00"
	assert.Equal(t, `0000000000000000000000000`, hex.EncodeToString(buf))

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
			"00000003884ed1c900000000884ed1c9090000000000000000",
			Transfer{AccountName("tbcox2.3"), AccountName("tbcox2"), 9, ""},
		},
		{
			"00000003884ed1c900000000884ed1c9090000000000000004616c6c6f",
			Transfer{AccountName("tbcox2.3"), AccountName("tbcox2"), 9, "allo"},
		},
	}

	for idx, test := range tests {
		buf, err := hex.DecodeString(test.in)
		assert.NoError(t, err)

		var res Transfer
		assert.NoError(t, UnmarshalBinary(buf, &res), fmt.Sprintf("Index %d", idx))
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
