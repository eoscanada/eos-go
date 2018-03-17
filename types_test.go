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

func TestPackAction(t *testing.T) {
	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("transfer"),
		Authorization: []PermissionLevel{
			{AccountName("eosio"), PermissionName("active")},
		},
		Data: Transfer{
			From:     AccountName("abourget"),
			To:       AccountName("eosio"),
			Quantity: 123123,
		},
	}

	buf, err := MarshalBinary(a)
	assert.NoError(t, err)
	// 0000000000ea3055000000572d3ccdcd010000000000ea305500000000a8ed3232
	// data: 32 (50 chars, the length of the data)
	//   from: 00000059b1abe931
	//   to: 0000000000ea3055
	//   quantity: f3e0010000000000   (quantity)
	//   memo: 00  (string)
	assert.Equal(t, `0000000000ea3055000000572d3ccdcd010000000000ea305500000000a8ed32323200000059b1abe9310000000000ea3055f3e001000000000000`, hex.EncodeToString(buf))

	buf, err = json.Marshal(a)
	assert.NoError(t, err)
	assert.Equal(t, `{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000000572d3ccdcd010000000000ea305500000000a8ed32323200000059b1abe9310000000000ea3055f3e001000000000000","from":"abourget","memo":"","name":"transfer","quantity":123123,"to":"eosio"}`, string(buf))
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

	//h, _ := hex.DecodeString("00409e9a2264b89a")
	h, _ := hex.DecodeString("0000001e4d75af46")
	fmt.Println("NAMETOSTRING", NameToString(binary.LittleEndian.Uint64(h)))
}

func TestNameToString(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"0000001e4d75af46", "currency"},
		{"0000000000ea3055", "eosio"},
		{"00409e9a2264b89a", "newaccount"},
		{"00000003884ed1c9", "tbcox2.3"},
		{"00000000a8ed3232", "active"},
		{"000000572d3ccdcd", "transfer"},
		{"00000059b1abe931", "abourget"},
	}

	for _, test := range tests {
		h, err := hex.DecodeString(test.in)
		require.NoError(t, err)
		res := NameToString(binary.LittleEndian.Uint64(h))
		assert.Equal(t, test.out, res)
	}
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
		Data: &Transfer{
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

func TestAuthorityBinaryMarshal(t *testing.T) {
	a := Authority{
		Threshold: 2,
		Keys: []KeyWeight{
			KeyWeight{
				PublicKey: "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
				Weight:    5,
			},
		},
	}
	cnt, err := MarshalBinary(a)
	require.NoError(t, err)

	// threshold: 02000000
	// []accounts: 00
	// []keys: 01
	// - pubkey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
	//   weight: 0500
	assert.Equal(t, `0200000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0500`, hex.EncodeToString(cnt))
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
