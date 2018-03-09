package eosapi

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/lunixbochs/struc"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"in", "out"},
	}

	for _, test := range tests {
		res := (test.in)
		assert.Equal(t, test.out, res)
	}
}

func TestUnpackBinaryTableRows(t *testing.T) {
	resp := &GetTableRowsResp{
		Rows: json.RawMessage(`["044355520000000004435552000000000000000000000000"]`),
	}
	assert.NoError(t, resp.BinaryToStructs(true))
}

func TestUnpackAccountName(t *testing.T) {
	tests := []struct {
		in  string
		out []byte
	}{
		{"eosio", []byte{0x55, 0x30, 0xea, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"eosio.system", []byte{0x55, 0x30, 0xea, 0x03, 0x1e, 0xc6, 0x55, 0x00}},
	}

	for _, test := range tests {
		acct := AccountName(test.in)
		var buf bytes.Buffer
		assert.NoError(t, struc.Pack(&buf, &acct))
		assert.Equal(t, test.out, buf.Bytes())
	}
}

type TestStruct struct {
	Account *AccountName `struc:""`
}
