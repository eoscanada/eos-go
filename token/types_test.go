package token

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	eos "github.com/eosioca/eosapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackAction(t *testing.T) {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("transfer"),
		Authorization: []eos.PermissionLevel{
			{AN("eosio"), PN("active")},
		},
		Data: Transfer{
			From:     AN("abourget"),
			To:       AN("eosio"),
			Quantity: eos.Asset{Amount: 123123, Symbol: eos.EOSSymbol},
		},
	}

	buf, err := eos.MarshalBinary(a)
	assert.NoError(t, err)
	assert.Equal(t, `0000000000ea3055000000572d3ccdcd010000000000ea305500000000a8ed32322100000059b1abe9310000000000ea3055f3e001000000000004454f530000000000`, hex.EncodeToString(buf))

	buf, err = json.Marshal(a)
	assert.NoError(t, err)
	assert.Equal(t, `{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"00000059b1abe9310000000000ea3055f3e001000000000004454f530000000000","name":"transfer"}`, string(buf))

	/* 0000000000ea3055 000000572d3ccdcd 01 0000000000ea3055 00000000a8ed3232
	   21
	   00000059b1abe931 0000000000ea3055 f3e0010000000000 04 454f5300000000 00 */
}

func TestUnpackActionTransfer(t *testing.T) {
	tests := []struct {
		in  string
		out Transfer
	}{
		{
			"00000003884ed1c900000000884ed1c90900000000000000000000000000000000",
			Transfer{AN("tbcox2.3"), AN("tbcox2"), eos.Asset{Amount: 9}, ""},
		},
		{
			"00000003884ed1c900000000884ed1c90900000000000000000000000000000004616c6c6f",
			Transfer{AN("tbcox2.3"), AN("tbcox2"), eos.Asset{Amount: 9}, "allo"},
		},
	}

	for idx, test := range tests {
		buf, err := hex.DecodeString(test.in)
		assert.NoError(t, err)

		var res Transfer
		assert.NoError(t, eos.UnmarshalBinary(buf, &res), fmt.Sprintf("Index %d", idx))
		assert.Equal(t, test.out, res)
	}

}

func TestActionMetaTypes(t *testing.T) {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("transfer"),
		Data: &Transfer{
			From: AN("abourget"),
			To:   AN("mama"),
		},
	}

	cnt, err := json.Marshal(a)
	require.NoError(t, err)
	assert.Equal(t,
		`{"account":"eosio","data":"00000059b1abe931000000000060a4910000000000000000000000000000000000","name":"transfer"}`,
		string(cnt),
	)

	var newAction eos.Action
	require.NoError(t, json.Unmarshal(cnt, &newAction))

	tx := &eos.Transaction{Actions: []*eos.Action{a}}
	stx := eos.NewSignedTransaction(tx)
	packed, err := stx.Pack(eos.TxOptions{})
	require.NoError(t, err)
	fmt.Println("MAMA1", stx.MaxNetUsageWords)
	fmt.Println("MAMA2", stx.Transaction.MaxNetUsageWords)
	assert.Equal(t, 123, stx.MaxNetUsageWords)
	packedData, err := json.Marshal(packed)
	fmt.Println("MAMA3", stx.Transaction.MaxNetUsageWords)
	assert.Equal(t, `000000000000000000000`, string(packedData))
	assert.Equal(t, `000000000000000000000`, hex.EncodeToString(packed.PackedTransaction))
}
