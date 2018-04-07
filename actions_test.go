package eos

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: Move this test to the `system` contract.. and take out
// `NewAccount` from this package.
func TestActionNewAccount(t *testing.T) {
	pubKey, err := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	require.NoError(t, err)
	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("newaccount"),
		Authorization: []PermissionLevel{
			{AccountName("eosio"), PermissionName("active")},
		},
		Data: NewActionData(NewAccount{
			Creator: AccountName("eosio"),
			Name:    AccountName("abourget"),
			Owner: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					KeyWeight{
						PublicKey: pubKey,
						Weight:    1,
					},
				},
			},
			Active: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					KeyWeight{
						PublicKey: pubKey,
						Weight:    1,
					},
				},
			},
			Recovery: Authority{
				Threshold: 1,
				Accounts: []PermissionLevelWeight{
					PermissionLevelWeight{
						Permission: PermissionLevel{AccountName("eosio"), PermissionName("active")},
						Weight:     1,
					},
				},
			},
		}),
	}
	tx := &Transaction{
		Actions: []*Action{a},
	}

	buf, err := MarshalBinary(tx)
	assert.NoError(t, err)

	assert.Equal(t, `00096e88000000000000000000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100`, hex.EncodeToString(buf))

	buf, err = json.Marshal(a)
	assert.NoError(t, err)
	assert.Equal(t, `{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100"}`, string(buf))

	buf, err = json.Marshal(a.Data)
	assert.NoError(t, err)
	assert.Equal(t, "{\"creator\":\"eosio\",\"name\":\"abourget\",\"owner\":{\"threshold\":1,\"keys\":[{\"public_key\":\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\",\"weight\":1}],\"accounts\":null},\"active\":{\"threshold\":1,\"keys\":[{\"public_key\":\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\",\"weight\":1}],\"accounts\":null},\"recovery\":{\"threshold\":1,\"keys\":null,\"accounts\":[{\"permission\":{\"actor\":\"eosio\",\"permission\":\"active\"},\"weight\":1}]}}", string(buf))
	// 00096e88 0000 0000 00000000 00 00 00 00 01 0000000000ea3055

	// WUTz that ?
	// var newAct *Action
	// newAct.DecodeAs(&NewAccount{})
	// require.NoError(t, UnmarshalBinary(buf, &newAct))
	// assert.Equal(t, a, newAct)
}

func TestMarshalTransactionAndSigned(t *testing.T) {
	a := &Action{
		Account: AccountName("eosio"),
		Name:    ActionName("newaccount"),
		Authorization: []PermissionLevel{
			{AccountName("eosio"), PermissionName("active")},
		},
		Data: NewActionData(NewAccount{
			Creator: AccountName("eosio"),
			Name:    AccountName("abourget"),
		}),
	}
	tx := &SignedTransaction{Transaction: &Transaction{
		Actions: []*Action{a},
	}}

	buf, err := MarshalBinary(tx)
	assert.NoError(t, err)
	// 00096e88 0000 0000 00000000 0000 0000 00
	// actions: 01
	// 0000000000ea3055 00409e9a2264b89a 01 0000000000ea3055 00000000a8ed3232
	// len: 22
	// 0000000000ea3055 00000059b1abe931 000000000000000000000000000000000000

	assert.Equal(t, `00096e88000000000000000000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232220000000000ea305500000059b1abe9310000000000000000000000000000000000000000`, hex.EncodeToString(buf))

	buf, err = json.Marshal(a)
	assert.NoError(t, err)
	assert.Equal(t, `{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500000059b1abe931000000000000000000000000000000000000"}`, string(buf))
}

func TestActionUnmarshalBinary(t *testing.T) {
	tests := []struct {
		in     string
		jsonTx string
	}{
		{
			"967abe5a000003002a48328c0000000000010000000000ea3055000000572d3ccdcd010000000000ea305500000000a8ed32322f0000000000ea3055000000023baca66200a0724e1809000004454f53000000000e57656c636f6d6520306461303930",
			"{\"expiration\":\"2018-03-30T17:57:42\",\"region\":0,\"ref_block_num\":3,\"ref_block_prefix\":2352105514,\"max_net_usage_words\":0,\"max_kcpu_usage\":0,\"delay_sec\":0}",
		},
	}

	for _, test := range tests {
		var tx Transaction
		b, err := hex.DecodeString(test.in)
		require.NoError(t, err)
		require.NoError(t, UnmarshalBinary(b, &tx))

		js, err := json.Marshal(tx)
		require.NoError(t, err)
		assert.Equal(t, test.jsonTx, string(js))
	}

}

// FETCHED FROM A SIMILAR TRANSACTION VIA `eosioc`, includes the Transaction headers though:
// This was BEFORE the `keys` and `accounts` were swapped on `Authority`.
// transaction header:
//    expiration epoch: 1e76ac5a
//    region: 0000
//    blocknum: 62cf
//    blockprefix: 50090bd8
//    packedbandwidthwords: 0000
//    contexfreecpubandwidth: 0000
//    []ContextFreeActions: 00
//    []Actions: 01
// Action idx 0:
//  account: 0000000000ea3055 (eosio)
//  name: 00409e9a2264b89a (newaccount)
//  []authorizations: 01
//   - actor: 0000000000ea3055 (eosio)
//     permission: 00000000a8ed3232 (active)
//  data len: 7c (124, indeed the length of the following...
//  creator: 0000000000ea3055
//  name: 0000001e4d75af46
//  owner authority:
//   threshold: 01000000
//   []accounts: 00
//   []keys: 01
//     - publickey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf // fixed width.
//       weight: 0100
//  active authority:
//   threshold: 01000000
//   []accounts: 00
//   []keys: 01
//     - publickey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
//       weight: 0100
//  recovery authority:
//    threshold: 01000000
//    []accounts: 01
//    - name: 0000000000ea3055 (eosio)
//    []keys: 00
// now the `newaccount` struct is done.. what,s that ?
// a list of a new object: 01
// an account name:
// a permission name: 00000000a8ed3232 (active)
// some list with one thing: 01
//   - an empty list: 00
//   - another empty list: 00

// 0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32326a0000000000ea305500000059b1abe9310100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000000000000
// 0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100
// Generated by the first run:
// account: 0000000000ea3055 (eosio)
// name: 00409e9a2264b89a (newaccount)
// []authorizations: 01
//  - actor: 0000000000ea3055 (eosio)
//    permission: 00000000a8ed3232 (active)
// data length: 6a (106) which MATCHES the lengths to follow.
// NewAccount:
//  creator: 0000000000ea3055 (eosio)
//  name: 00000059b1abe931 (abourget)
// owner-authority:
//  threshold: 01000000
//  []keys: 01
//  - publickey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
//    weight: 0100
//  []accounts: 00
// active-authority:
//  threshold: 01000000
//  []keys: 01
//  - pubkey: 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
//    weight: 0100
//  []accounts: 00
// recovery-authority:  // the last bit is the Recovery authority.. it works :)
//  threshold: 00000000
//  []keys: 00
//  []accounts: 00

// 0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232a4010000000000ea305500000059b1abe9310100000001   BINARY SERIALIZER FAILED: 35454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed32320100

/**
574fbd5a
0000
0300
0b859e4e
0000
0000
00 = ctx-free-actions
01
0000000000ea3055
000000572d3ccdcd
01
0000000000ea3055
00000000a8ed3232
2f
0000000000ea3055 eosio
000000003baca662 genesis
00a0724e18090000 a billion
04454f5300000000 precision 4 + EOS string

0e57656c636f6d6520303030304231
*/

/**
1055bd5a
0000
0200
b34c9a0e
000000000001
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
7c
0000000000ea3055
010000003baca662
0100000001000260520ba1782b60f9a658aff7b6d8536cf9088d509608bca5aae66dc171cba9030100000100000001000260520ba1782b60f9a658aff7b6d8536cf9088d509608bca5aae66dc171cba9030100000100000000010000000000ea305500000000a8ed32320100


c157bd5a00000400e1c97e580000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055
62a6ac3b0100000001
00000001000260520ba1782b60f9a658aff7b6d8536cf9088d509608bca5aae66dc171cba9030100000100000001000260520ba1782b60f9a658aff7b6d8536cf9088d509608bca5aae66dc171cba9030100000100000000010000000000ea305500000000a8ed32320100

0c58bd5a00000100adf72e080000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055
010000003baca66201
00000001000260520ba1782b60f9a658aff7b6d8536cf9088d509608bca5aae66dc171cba9030100000100000001000260520ba1782b60f9a658aff7b6d8536cf9088d509608bca5aae66dc171cba9030100000100000000010000000000ea305500000000a8ed32320100

967abe5a
0000
0300
2a48328c
0000
0000
00
01
0000000000ea3055000000572d3ccdcd010000000000ea305500000000a8ed32322f0000000000ea3055000000023baca66200a0724e1809000004454f53000000000e57656c636f6d6520306461303930
*/
