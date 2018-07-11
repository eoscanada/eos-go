package eos

import (
	"testing"

	"encoding/json"

	eos "github.com/eoscanada/eos-go"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalAccountResp(t *testing.T) {
	resp := &eos.AccountResp{}

	err := json.Unmarshal([]byte(jsonData), resp)
	assert.NoError(t, err)

	assert.Equal(t, eos.AccountName("eosriobrazil"), resp.AccountName)
}

var jsonData = `{
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
