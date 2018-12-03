package eos

import (
	"encoding/hex"
	"testing"
	"time"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalAccountResp(t *testing.T) {
	resp := &AccountResp{}

	err := json.Unmarshal([]byte(accountResponseJSONData), resp)
	assert.NoError(t, err)

	assert.Equal(t, AccountName("eosriobrazil"), resp.AccountName)
}

func TestUnmarshalBlockResp(t *testing.T) {
	RegisterAction(AccountName("eosdactokens"), ActionName("transfer"), Transfer{})

	resp := &BlockResp{}
	err := json.Unmarshal([]byte(blockResponseJSONData), resp)
	assert.NoError(t, err)

	timestamp, _ := time.Parse("2006-01-02T15:04:05.000", "2018-07-17T11:18:10.000")

	// Block Header
	assert.Equal(t, BlockTimestamp{timestamp}, resp.Timestamp)
	assert.Equal(t, AccountName("bitfinexeos1"), resp.Producer)
	assert.Equal(t, uint16(0), resp.Confirmed)
	assert.Equal(t, hexToChecksum256("0060c5c3527812daa61b8a93ac82329eb1995bbafd37fcba4e15b25cc94cea1a"), resp.Previous)
	assert.Equal(t, hexToChecksum256("f9590621c12ac92b4ab621a30468497918287fd0dd2535907bd85409947f1988"), resp.TransactionMRoot)
	assert.Equal(t, hexToChecksum256("5827604590c993e2489639b739af79ab12dd0442fee6551e9cd2932a36fa5026"), resp.ActionMRoot)
	assert.Equal(t, uint32(149), resp.ScheduleVersion)
	assert.Nil(t, resp.NewProducers)
	assert.Empty(t, resp.HeaderExtensions)

	// Signed Block Header
	assert.Equal(t, "SIG_K1_K6ocyo6k5zcYKxJ3XAUxphrwBs8Y3Y5etotX8wvr3hxDJiaQeGseV79rXwdHt6dEy2zFoyrmReGATob6EYJ85assYKudc4", resp.ProducerSignature.String())

	// Signed Block
	assert.Len(t, resp.Transactions, 1)
	assert.Empty(t, resp.BlockExtensions)

	// Packged Transaction
	receipt := resp.Transactions[0]
	assert.Equal(t, uint32(1529), receipt.CPUUsageMicroSeconds)
	assert.Equal(t, Varuint32(16), receipt.NetUsageWords)
	assert.Equal(t, TransactionStatusExecuted, receipt.Status)
	assert.Equal(t, CompressionNone, receipt.Transaction.Packed.Compression)
	assert.Equal(t, HexBytes{}, receipt.Transaction.Packed.PackedContextFreeData)

	id, err := receipt.Transaction.Packed.ID()
	assert.NoError(t, err)

	assert.Equal(t, hexToChecksum256("7074b6caaac4dfe1d19903a41b88a53b595e963bab02139a508785eba6e11ba5"), id)
	assert.Len(t, receipt.Transaction.Packed.Signatures, 1)
	assert.Equal(t, "SIG_K1_KXsd17mt6qf8JAHvRiVLRH93tMoQrkC69qhoS2suG8N3YYF54LTVkSwnh4t4wscDJXPnSAdbJZpSfHjJjSurDmwGCAxvTs", receipt.Transaction.Packed.Signatures[0].String())

	// Signed Transaction
	transaction, err := receipt.Transaction.Packed.Unpack()
	assert.NoError(t, err)

	expiration, _ := time.Parse("2006-01-02T15:04:05", "2018-07-17T11:19:05")

	assert.Equal(t, Varuint32(0), transaction.DelaySec)
	assert.Equal(t, JSONTime{expiration}, transaction.Expiration)
	assert.Equal(t, uint8(0), transaction.MaxCPUUsageMS)
	assert.Equal(t, Varuint32(0), transaction.MaxNetUsageWords)
	assert.Equal(t, uint16(50301), transaction.RefBlockNum)
	assert.Equal(t, uint32(2432041012), transaction.RefBlockPrefix)
	assert.Empty(t, transaction.Extensions)

	assert.Empty(t, transaction.ContextFreeActions)
	assert.Len(t, transaction.Actions, 1)

	// Action
	action := transaction.Actions[0]

	assert.Equal(t, AccountName("eosdactokens"), action.Account)
	assert.Len(t, action.Authorization, 1)
	assert.Equal(t, AccountName("gy4dkmjzhege"), action.Authorization[0].Actor)
	assert.Equal(t, PermissionName("active"), action.Authorization[0].Permission)
	assert.Equal(t, ActionName("transfer"), action.Name)
	assert.Equal(t, hexToHexBytes("a0986aff49988867100261f9519b8867881300000000000004454f534441430000"), action.HexData)

	// For this to work correctly, you must call `RegisterAction` with the right Data struct prior deserialization
	//  @see Top of this test method for the `RegisterAction` call
	transfer, ok := action.Data.(*Transfer)
	assert.True(t, ok)

	assert.Equal(t, AccountName("gy4dkmjzhege"), transfer.From)
	assert.Equal(t, AccountName("gy4dqojtg411"), transfer.To)
	assert.Equal(t, Asset{Amount: 5000, Symbol: Symbol{Precision: uint8(4), Symbol: "EOSDAC"}}, transfer.Quantity)
	assert.Equal(t, "", transfer.Memo)
}

var accountResponseJSONData = `{
	"account_name": "eosriobrazil",
	"head_block_num": 5264738,
	"head_block_time": "2018-07-11T04:29:16.500",
	"privileged": false,
	"last_code_update": "1970-01-01T00:00:00.000",
	"created": "2018-06-10T13:09:26.500",
	"core_liquid_balance": "695.2674 EOS",
	"ram_quota": 145360,
	"net_weight": 324628,
	"cpu_weight": 329628,
	"net_limit": {
		"used": 550,
		"available": 17636233,
		"max": 17636783
	},
	"cpu_limit": {
		"used": 14728,
		"available": 3396286,
		"max": 3411014
	},
	"ram_usage": 5935,
	"permissions": [
		{
		"perm_name": "active",
		"parent": "owner",
		"required_auth": {
			"threshold": 1,
			"keys": [
			{
				"key": "EOS6HSE9SVvNmGF4Dv8cHLUjF8BigorYykUG2z8UbHZd1BQ9qF88r",
				"weight": 1
			}
			],
			"accounts": [],
			"waits": []
		}
		},
		{
		"perm_name": "claim",
		"parent": "active",
		"required_auth": {
			"threshold": 1,
			"keys": [
			{
				"key": "EOS7FJJ7igorHoTq6y6yd7GmRei9cc6CRhR7L2TXP6H9UFEP49jNc",
				"weight": 1
			}
			],
			"accounts": [],
			"waits": []
		}
		},
		{
		"perm_name": "owner",
		"parent": "",
		"required_auth": {
			"threshold": 1,
			"keys": [
			{
				"key": "EOS5UhWBMYKPPzb4tigorbnrH9Ft7mogW1MmvViaHJkBif2kSa1f4",
				"weight": 1
			}
			],
			"accounts": [],
			"waits": []
		}
		}
	],
	"total_resources": {
		"owner": "eosriobrazil",
		"net_weight": "32.4628 EOS",
		"cpu_weight": "32.9628 EOS",
		"ram_bytes": 145360
	},
	"self_delegated_bandwidth": {
		"from": "eosriobrazil",
		"to": "eosriobrazil",
		"net_weight": "32.4628 EOS",
		"cpu_weight": "32.9628 EOS"
	},
	"refund_request": {
		"owner": "eosriobrazil",
		"request_time": "2018-07-09T20:54:31",
		"net_amount": "2.9284 EOS",
		"cpu_amount": "2.9284 EOS"
	},
	"voter_info": {
		"owner": "eosriobrazil",
		"proxy": "",
		"producers": [],
		"staked": 804256,
		"last_vote_weight": "171334771736.95532226562500000",
		"proxied_vote_weight": "0.00000000000000000",
		"is_proxy": 0
	}
}`

type Transfer struct {
	From     AccountName `json:"from"`
	To       AccountName `json:"to"`
	Quantity Asset       `json:"quantity"`
	Memo     string      `json:"memo"`
}

func hexToHexBytes(data string) HexBytes {
	bytes, _ := hex.DecodeString(data)

	return HexBytes(bytes)
}

func hexToChecksum256(data string) Checksum256 {
	return Checksum256(hexToHexBytes(data))
}

var blockResponseJSONData = `
{
	"action_mroot": "5827604590c993e2489639b739af79ab12dd0442fee6551e9cd2932a36fa5026",
	"block_extensions": [],
	"block_num": 6342084,
	"confirmed": 0,
	"header_extensions": [],
	"id": "0060c5c404a7e85d5c3e35cbaabfafad847c7c7c7035bdde307fb2c8777413f1",
	"new_producers": null,
	"previous": "0060c5c3527812daa61b8a93ac82329eb1995bbafd37fcba4e15b25cc94cea1a",
	"producer": "bitfinexeos1",
	"producer_signature": "SIG_K1_K6ocyo6k5zcYKxJ3XAUxphrwBs8Y3Y5etotX8wvr3hxDJiaQeGseV79rXwdHt6dEy2zFoyrmReGATob6EYJ85assYKudc4",
	"ref_block_prefix": 3409264220,
	"schedule_version": 149,
	"timestamp": "2018-07-17T11:18:10.000",
	"transaction_mroot": "f9590621c12ac92b4ab621a30468497918287fd0dd2535907bd85409947f1988",
	"transactions": [
	  {
		"cpu_usage_us": 1529,
		"net_usage_words": 16,
		"status": "executed",
		"trx": {
		  "compression": "none",
		  "context_free_data": [],
		  "id": "7074b6caaac4dfe1d19903a41b88a53b595e963bab02139a508785eba6e11ba5",
		  "packed_context_free_data": "",
		  "packed_trx": "a9d04d5b7dc43400f690000000000180a7823423933055000000572d3ccdcd01a0986aff4998886700000000a8ed323221a0986aff49988867100261f9519b8867881300000000000004454f53444143000000",
		  "signatures": [
			"SIG_K1_KXsd17mt6qf8JAHvRiVLRH93tMoQrkC69qhoS2suG8N3YYF54LTVkSwnh4t4wscDJXPnSAdbJZpSfHjJjSurDmwGCAxvTs"
		  ],
		  "transaction": {
			"actions": [
			  {
				"account": "eosdactokens",
				"authorization": [
				  {
					"actor": "gy4dkmjzhege",
					"permission": "active"
				  }
				],
				"data": {
				  "from": "gy4dkmjzhege",
				  "memo": "",
				  "quantity": "0.5000 EOSDAC",
				  "to": "gy4dqojtg411"
				},
				"hex_data": "a0986aff49988867100261f9519b8867881300000000000004454f534441430000",
				"name": "transfer"
			  }
			],
			"context_free_actions": [],
			"delay_sec": 0,
			"expiration": "2018-07-17T11:19:05",
			"max_cpu_usage_ms": 0,
			"max_net_usage_words": 0,
			"ref_block_num": 50301,
			"ref_block_prefix": 2432041012,
			"transaction_extensions": []
		  }
		}
	  }
	]
  }
`
