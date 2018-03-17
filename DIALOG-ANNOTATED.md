RUNNING THE COMMAND:

ec create account eosio currency EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV




POST /v1/chain/get_info HTTP/1.0
Host: testnet-dawn3.eosio.ca
content-length: 0
Accept: */*
Connection: close

HTTP/1.0 200 OK
Content-Length: 370
Content-Type: application/json
Server: WebSocket++/0.7.0
Date: Sat, 17 Mar 2018 01:57:20 GMT

{
  "server_version": "4c9eed11",
  "head_block_num": 1232738,
  "last_irreversible_block_num": 1232737,
  "head_block_id": "0012cf6247be7e2050090bd83b473369b705ba1d280cd55d3aef79998c784b9b",
  "head_block_time": "2018-03-17T01:57:20",
  "head_block_producer": "eosio",
  "recent_slots": "1111111111111111111111111111111111111111111111111111111111111111",
  "participation_rate": "1.00000000000000000"
}


----

POST /v1/chain/get_required_keys

{
  "transaction": {
    "expiration": "2018-03-17T01:57:50",
    "region": 0,
    "ref_block_num": 53090,
    "ref_block_prefix": 3624601936,
    "packed_bandwidth_words": 0,
    "context_free_cpu_bandwidth": 0,
    "context_free_actions": [],
    "actions": [
      {
        "account": "eosio",
        "name": "newaccount",
        "authorization": [
          {
            "actor": "eosio",
            "permission": "active"
          }
        ],
        "data": "0000000000ea30550000001e4d75af460100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf010001000000010000000000ea305500000000a8ed3232010000"
      }
    ]
  },
  "available_keys": [
    "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"
  ]
}

HTTP/1.0 200 OK
Content-Length: 75
Content-Type: application/json
Server: WebSocket++/0.7.0
Date: Sat, 17 Mar 2018 01:57:21 GMT

{"required_keys":["EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"]}

-------


POST /v1/wallet/sign_transaction HTTP/1.0
Host: localhost
content-length: 717
Accept: */*
Connection: close

[
  {
    "expiration": "2018-03-17T01:57:50",
    "region": 0,
    "ref_block_num": 53090,
    "ref_block_prefix": 3624601936,
    "packed_bandwidth_words": 0,
    "context_free_cpu_bandwidth": 0,
    "context_free_actions": [],
    "actions": [
      {
        "account": "eosio",
        "name": "newaccount",
        "authorization": [
          {
            "actor": "eosio",
            "permission": "active"
          }
        ],
        "data": "0000000000ea30550000001e4d75af460100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf010001000000010000000000ea305500000000a8ed3232010000"
      }
    ],
    "signatures": [],
    "context_free_data": []
  },
  [
    "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"
  ],
  "0000000000000000000000000000000000000000000000000000000000000000"
]

HTTP/1.1 201 Created
Content-Length: 689
Content-type: application/json
Server: WebSocket++/0.7.0

{
  "expiration": "2018-03-17T01:57:50",
  "region": 0,
  "ref_block_num": 53090,
  "ref_block_prefix": 3624601936,
  "packed_bandwidth_words": 0,
  "context_free_cpu_bandwidth": 0,
  "context_free_actions": [],
  "actions": [
    {
      "account": "eosio",
      "name": "newaccount",
      "authorization": [
        {
          "actor": "eosio",
          "permission": "active"
        }
      ],
      "data": "0000000000ea30550000001e4d75af460100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf010001000000010000000000ea305500000000a8ed3232010000"
    }
  ],
  "signatures": [
    "EOSKkA4MtvipSgKJWMp6J8G2f3RmVTLqQ47zG4rnjo6YGXfE6DS6s7iVDZrMNfaxAGfnyJ43vfuBbKSi9TL4ahwmRU2d7tHou"
  ],
  "context_free_data": []
}


-----------------------------------

POST /v1/chain/push_transaction HTTP/1.0
Host: testnet-dawn3.eosio.ca
content-length: 499
Accept: */*
Connection: close

{
  "signatures": [
    "EOSKkA4MtvipSgKJWMp6J8G2f3RmVTLqQ47zG4rnjo6YGXfE6DS6s7iVDZrMNfaxAGfnyJ43vfuBbKSi9TL4ahwmRU2d7tHou"
  ],
  "compression": "none",
  "data": "1e76ac5a000062cf50090bd80000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea30550000001e4d75af460100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf010001000000010000000000ea305500000000a8ed3232010000"
}

HTTP/1.0 202 Accepted
Content-Length: 744
Content-Type: application/json
Server: WebSocket++/0.7.0
Date: Sat, 17 Mar 2018 01:57:22 GMT

{
  "transaction_id": "779304bc16683538943d9aa83972e90194829a12c6c437627c4df272b76f7ec7",
  "processed": {
    "status": "executed",
    "id": "779304bc16683538943d9aa83972e90194829a12c6c437627c4df272b76f7ec7",
    "action_traces": [
      {
        "receiver": "eosio",
        "act": {
          "account": "eosio",
          "name": "newaccount",
          "authorization": [
            {
              "actor": "eosio",
              "permission": "active"
            }
          ],
          "data": "0000000000ea30550000001e4d75af460100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf010001000000010000000000ea305500000000a8ed3232010000"
        },
        "console": "",
        "region_id": 0,
        "cycle_index": 0,
        "data_access": [
          {
            "type": "write",
            "code": "eosio",
            "scope": "eosio.auth",
            "sequence": 0
          }
        ]
      }
    ],
    "deferred_transactions": []
  }
}

-----


POST /v1/account_history/get_transaction HTTP/1.0
Host: testnet-dawn3.eosio.ca
content-length: 85
Accept: */*
Connection: close

{"transaction_id":"779304bc16683538943d9aa83972e90194829a12c6c437627c4df272b76f7ec7"}


HTTP/1.0 200 OK
Content-Length: 1162
Content-Type: application/json
Server: WebSocket++/0.7.0
Date: Sat, 17 Mar 2018 01:57:58 GMT

{
  "transaction_id": "779304bc16683538943d9aa83972e90194829a12c6c437627c4df272b76f7ec7",
  "transaction": {
    "signatures": [
      "EOSKkA4MtvipSgKJWMp6J8G2f3RmVTLqQ47zG4rnjo6YGXfE6DS6s7iVDZrMNfaxAGfnyJ43vfuBbKSi9TL4ahwmRU2d7tHou"
    ],
    "compression": "none",
    "hex_data": "1e76ac5a000062cf50090bd80000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea30550000001e4d75af460100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf010001000000010000000000ea305500000000a8ed3232010000",
    "data": {
      "expiration": "2018-03-17T01:57:50",
      "region": 0,
      "ref_block_num": 53090,
      "ref_block_prefix": 3624601936,
      "packed_bandwidth_words": 0,
      "context_free_cpu_bandwidth": 0,
      "context_free_actions": [],
      "actions": [
        {
          "account": "eosio",
          "name": "newaccount",
          "authorization": [
            {
              "actor": "eosio",
              "permission": "active"
            }
          ],
          "data": "0000000000ea30550000001e4d75af460100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000100000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf010001000000010000000000ea305500000000a8ed3232010000"
        }
      ]
    }
  }
}
