
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












-------------------------


Successful newaccount with my program:

github.com/eoscanada/eos-go/cmd/eosapi
POST /v1/chain/get_info HTTP/1.1
Host: localhost:8888


refblockprefix: 3600688091



POST /v1/wallet/get_public_keys HTTP/1.1
Host: localhost:6666




POST /v1/chain/get_required_keys HTTP/1.1
Host: localhost:8888

{"available_keys":["EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV","EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM"],"transaction":{"expiration":"2018-03-19T15:07:46","ref_block_num":8182,"ref_block_prefix":3600688091,"actions":[{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232a4010000000000ea305500000059b1abe931010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed32320100","name":"newaccount"}]}}


OUTPUT:
2018/03/19 11:07:16 GetRequiredKeys &{[EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV]} <nil>





POST /v1/wallet/sign_transaction HTTP/1.1
Host: localhost:6666

[{"expiration":"2018-03-19T15:07:46","ref_block_num":8182,"ref_block_prefix":3600688091,"actions":[{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232a4010000000000ea305500000059b1abe931010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed32320100","name":"newaccount"}]},["EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"],"0000000000000000000000000000000000000000000000000000000000000000"]


were able to SIGN the thing:


ERROR calling NewAccount: Sign: status code=201, body={"expiration":"2018-03-19T15:07:46","region":0,"ref_block_num":8182,"ref_block_prefix":3600688091,"packed_bandwidth_words":0,"context_free_cpu_bandwidth":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232a4010000000000ea305500000059b1abe931010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":["EOSKgQmc9WbHff9WmzJNETjCzxbJroEqxf4mj95eYHaJRk2ZRnY2rj91JEiE15Jf9qKFGqYGHvC7H2CAG9KgzLb6VLVqatGi2"],"context_free_data":[]}





----------------------

now fails at `push_transaction` (from a second call:)

POST /v1/chain/push_transaction HTTP/1.1
Host: localhost:8888

{"compression":"none","data":"fcd2af5a00006a21887363b50000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32320000000000ea305500000859b1abe931010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed32320100","signatures":["EOS32JRPThMCYkmgfdKSnxGGjpbLN6quxYyfvzZLjKmrUqVb5x1GWojCoPwHSMAhvSVADLeGAgBfzgtuW6avuQFoq5g8SAyHpCo4JCQ"]}


ERROR calling NewAccount: status code=500, body={"code":500,"message":"Internal Service Error","error":{"code":10,"name":"assert_exception","message":"Assert Exception","details":"false: No matching prefix for 32JRPThMCYkmgfdKSnxGGjpbLN6quxYyfvzZLjKmrUqVb5x1GWojCoPwHSMAhvSVADLeGAgBfzgtuW6avuQFoq5g8SAyHpCo4JCQ","stack_trace":[{"level":"error","file":"common.hpp","line":60,"method":"apply","hostname":"","thread_name":"thread-0","timestamp":"2018-03-19T15:10:22.151"},{"level":"error","file":"abi_serializer.hpp","line":428,"method":"from_variant","hostname":"","thread_name":"thread-0","timestamp":"2018-03-19T15:10:22.151"}]}}


exp: abd6af5a
0000
c828
90ec4c36
0000
0000
00
01
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232

0000000000ea305500000859b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

expiration: fcd2af5a
region: 0000
refblock: 6a21
prefix: 887363b5
packed bandwidth: 0000
contexcpuband: 0000
contextfreeactions: 00
actions:01
- account: 0000000000ea3055
  name: 00409e9a2264b89a
  authorization: 01
  - actor: 0000000000ea3055
    perm: 00000000a8ed3232
  data: NO PREFIX!
    0000000000ea305500000859b1abe931010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed32320100


---------------

A transaction from `eosioc`:

POST /v1/chain/push_transaction HTTP/1.0
Host: localhost
content-length: 499
Accept: */*
Connection: close

{"signatures":["EOSK5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnNK"],"compression":"none","data":"20d8af5a0000b32bcc0e37eb0000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500001059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100"}

ACCEPTED 202:

{"transaction_id":"1f5e90b39175258ab507e57264a636436bc14b0f3e907a086b7a617473e55eb4","processed":{"status":"executed","id":"1f5e90b39175258ab507e57264a636436bc14b0f3e907a086b7a617473e55eb4","action_traces":[{"receiver":"eosio","act":{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500001059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100"},"console":"","region_id":0,"cycle_index":0,"data_access":[{"type":"write","code":"eosio","scope":"eosio.auth","sequence":2}]}],"deferred_transactions":[]}}


20d8af5a
0000
b32b
cc0e37eb
0000
0000
00
01
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
LENGTH: 7c
0000000000ea3055
00001059b1abe931
01000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

My tx:

0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000859b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

Returned tx:

d2d8af5a0000172dfd85a6840000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32320000000000ea305500000859b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100



8af0af5a
0000
875c
016c94f0
0000
0000
00
01
first action:
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
0000000000ea305500000859b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

02f1af5a0000765d57c697460000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232
0000000000ea305500000859b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100


f0f2af5a
0000
5261
0b9e01d4
0000
0000
00
01
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
7c
0000000000ea3055
00001859b1abe931
01000000
01
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
0100
00

01000000
01
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
0100
00

01000000
00
01
0000000000ea3055
00000000a8ed3232
0100


00096e88
0000
0000
00000000
0000
0000
00
01

ACTION 1:

0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
length: 7c
0000000000ea3055
00000059b1abe931
01000000
01 keys
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
0100 weight
00 accounts
01000000 thres
01 keys
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
0100 weight
00 accounts
01000000
00 keys
01
0000000000ea3055
00000000a8ed3232
01
00




--------------------------

From a special built wallet, spittin the serialized transaction when I send:

110db05a 0000 6a82 ef3dfbfd 0000 0000 00 01
ACTION 1:
0000000000ea3055 00409e9a2264b89a 01 0000000000ea3055 00000000a8ed3232
LEN: 7c
0000000000ea3055
00000859b1abe931
01
000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

SERVER internal representation:

110db05a 0000 6a82 ef3dfbfd 0000 0000 00 01
ACTION 1: 0000000000ea3055 00409e9a2264b89a 01 0000000000ea3055 00000000a8ed3232
LEN: 9e
01
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
7c
0000000000ea305500000859b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed323201000000



// FROM EOSIO same newaccount:

620eb05a 0000 0d85 438f55d2 0000 0000 00 01
ACTION 1:
0000000000ea3055 00409e9a2264b89a 01 0000000000ea3055 00000000a8ed3232
7c
0000000000ea3055 00002059b1abe931
01000000 01 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf 0100 00 (
01000000 01 0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf 0100 00
01000000 00 0100 00000000ea305500000000a8ed323201000000 (recovery)






MY SIGN TX JSON:

[
  {
    "expiration": "2018-03-19T15:07:46",
    "ref_block_num": 8182,
    "ref_block_prefix": 3600688091,
    "actions": [
      {
        "account": "eosio",
        "authorization": [
          {
            "actor": "eosio",
            "permission": "active"
          }
        ],
        "data": "0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232a4010000000000ea305500000059b1abe931010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed32320100",
        "name": "newaccount"
      }
    ]
  },
  [
    "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"
  ],
  "0000000000000000000000000000000000000000000000000000000000000000"
]
// [{"expiration":"2018-03-19T15:07:46",  "ref_block_num":8182,"ref_block_prefix":3600688091,"actions":[{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232a4010000000000ea305500000059b1abe931010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed32320100","name":"newaccount"}]},["EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"],"0000000000000000000000000000000000000000000000000000000000000000"]



EOSIOC SIGN TX JSON:

[
  {
    "expiration": "2018-03-19T19:38:44",
    "region": 0,
    "ref_block_num": 35331,
    "ref_block_prefix": 1388235515,
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
        "data": "0000000000ea305500003059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100"
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


[{"expiration":"2018-03-19T19:38:44","region":0,"ref_block_num":35331,"ref_block_prefix":1388235515,"packed_bandwidth_words":0,"context_free_cpu_bandwidth":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500003059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":[],"context_free_data":[]},["EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"],"0000000000000000000000000000000000000000000000000000000000000000"]



-----------

From `eosapi` `newaccount`:

c411b05a 0000 038a fbcabe52 0000 0000 00 01
Action 1: 0000000000ea3055 00409e9a2264b89a 01 0000000000ea3055 00000000a8ed3232
len: c701
0000000000ea3055
00409e9a2264b89a (new account)
01
0000000000ea3055
00000000a8ed3232
len: a401 AGAIN
0000000000ea3055
00000059b1abe931 (abourgetXX)
01
0000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f384268744875475971455435474457354356010000010000000135454f53364d5279416a51713875643768564e5963666e56504a7163567073634e35536f3842687448754759714554354744573543560100000100000000010000000000ea305500000000a8ed323201
000000


From `eosioc` `newaccount`:

c411b05a 0000 038a fbcabe52 0000 0000 00 01
action 1: 0000000000ea3055 00409e9a2264b89a 010000000000ea3055 00000000a8ed3232
len: 7c
0000000000ea3055
00003059b1abe931 (abourgetXX)
01
000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed323201000000


----




`eosapi`, who WE serialize it for signature (thus what we send on `push_transaction` later):

6c1bb05a0000549dad9012f90000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232
7c
0000000000ea3055
00000859b1abe931
01000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

how `wallet` sees it based on our `JSON`, and thus signs:

6c1bb05a0000549dad9012f90000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232
9e01
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
7c
0000000000ea3055
00000859b1abe931
01000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100
0000


When we run with `eosioc`, the data that ends up on the blockchain (through getting the tx ID) is:

381cb05a0000e99e0aae55180000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500004059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

which is identical to what the `wallet` sees, and thus signs and thus sends to `push_transaction`:

381cb05a0000e99e0aae55180000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500004059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100 0000
(these last bytes are the signature and context free data lists)


So there must be a difference between OUR JSON, and `eosioc`'s JSON:

A network trace tries them both, when calling `/v1/wallet/sign_transaction`:

this is `eosapi`:

[{"expiration":"2018-03-19T20:26:30","region":0,"ref_block_num":41061,"ref_block_prefix":1229959025,"actions":[{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000859b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100","name":"newaccount"}],"signatures":[],"context_free_data":[]},["EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"],"0000000000000000000000000000000000000000000000000000000000000000"]

with data being:
0000000000ea3055 00409e9a2264b89a 01 0000000000ea3055 00000000a8ed3232
7c
0000000000ea305500000859b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

this is `eosioc`:

[{"expiration":"2018-03-19T19:38:44","region":0,"ref_block_num":35331,"ref_block_prefix":1388235515,"packed_bandwidth_words":0,"context_free_cpu_bandwidth":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500003059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":[],"context_free_data":[]},["EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"],"0000000000000000000000000000000000000000000000000000000000000000"]

with data being:
0000000000ea3055 00003059b1abe931 01 000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100










------------------------

GO LANG:

POST /v1/chain/push_transaction HTTP/1.1
Host: localhost:8889
User-Agent: Go-http-client/1.1
Content-Length: 421
Accept-Encoding: gzip
Connection: close

{"signatures":[],"context_free_data":[],"compression":"none","data":"70e1c15a00000200e2976d0300000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100"}HTTP/1.1 401 Unauthorized
Content-Length: 553
Content-type: application/json
Server: WebSocket++/0.7.0

{"code":401,"message":"UnAuthorized","error":{"code":3030002,"name":"tx_missing_sigs","what":"signatures do not satisfy declared authorizations","details":[{"message":"transaction declares authority '{\"actor\":\"eosio\",\"permission\":\"active\"}', but does not have signatures for it.","file":"chain_controller.cpp","line_number":972,"method":"check_authorization"},{"message":"","file":"chain_controller.cpp","line_number":346,"method":"_push_transaction"},{"message":"","file":"chain_controller.cpp","line_number":271,"method":"push_transaction"}]}}


CLEOS:


POST /v1/chain/push_transaction HTTP/1.0
Host: localhost
content-length: 522
Accept: */*
Connection: close

{"signatures":["EOSKigRNhVVat7ZpumFgbqnRmFPJisJNVr2zvYWqcfop2u9QJ2KGViZGYakMyUz2UsxndoCB9ZKCBwk5tdtydQLcZ73XpVrQy"],"context_free_data":[],"compression":"none","data":"81e1c15a00002300be63087b21e8070000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"}HTTP/1.1 202 Accepted
Content-Length: 956
Content-type: application/json
Server: WebSocket++/0.7.0

{"transaction_id":"10768cb66940db434a23b477675ec884b5c835a43bbeca6a09b5fcacbe82cf2a","processed":{"status":"executed","id":"10768cb66940db434a23b477675ec884b5c835a43bbeca6a09b5fcacbe82cf2a","action_traces":[{"receiver":"eosio","context_free":false,"cpu_usage":0,"act":{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"},"console":"","region_id":0,"cycle_index":0,"data_access":[{"type":"write","code":"eosio","scope":"eosio.auth","sequence":0}],"_profiling_us":87}],"deferred_transaction_requests":[],"read_locks":[],"write_locks":[{"account":"eosio","scope":"eosio.auth"}],"cpu_usage":1000,"net_usage":364,"_profiling_us":188,"_setup_profiling_us":0}}





70e1c15a 0000 0200 e2976d03 00000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100

91e2c15a 0000 0200 8c87d0de 17 02 00 00 01 0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100

81e1c15a 0000 2300 be63087b 21 e807 00 00 010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100




---------------

GO

POST /v1/chain/push_transaction HTTP/1.1
Host: localhost:8889
User-Agent: Go-http-client/1.1
Content-Length: 520
Accept-Encoding: gzip
Connection: close

{"signatures":["EOSJxw7WFxgM5xTCfU79wkk8RCy1b8u5D5GWE4VX7fx8ThTB2tBLktXaZgoYeKD9xjX5n3jydvy2NJGuJAJu2ccnxBCXgJzny"],"context_free_data":[],"compression":"none","data":"86e4c15a0000020019e992a700000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100"}HTTP/1.1 401 Unauthorized
Content-Length: 553
Content-type: application/json
Server: WebSocket++/0.7.0

{"code":401,"message":"UnAuthorized","error":{"code":3030002,"name":"tx_missing_sigs","what":"signatures do not satisfy declared authorizations","details":[{"message":"transaction declares authority '{\"actor\":\"eosio\",\"permission\":\"active\"}', but does not have signatures for it.","file":"chain_controller.cpp","line_number":972,"method":"check_authorization"},{"message":"","file":"chain_controller.cpp","line_number":346,"method":"_push_transaction"},{"message":"","file":"chain_controller.cpp","line_number":271,"method":"push_transaction"}]}}


CLEOS


POST /v1/chain/push_transaction HTTP/1.0
Host: localhost
content-length: 522
Accept: */*
Connection: close

{"signatures":["EOSKfA5xarhtSop6v5s2Cse1QdBxTcGufvpzx7wZ4bNrp6s2yVU5gcRPzafyoEDKxEpzFyaQzFAYDUzDAqxMXx5jUc5z6tgEN"],"context_free_data":[],"compression":"none","data":"8ae4c15a00000a00f150f7ca21e8070000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"}HTTP/1.1 202 Accepted
Content-Length: 955
Content-type: application/json
Server: WebSocket++/0.7.0

{"transaction_id":"3908e02af8f2b829b22daa5e4ff55bc1f83d301147667d56b03def2a9665421b","processed":{"status":"executed","id":"3908e02af8f2b829b22daa5e4ff55bc1f83d301147667d56b03def2a9665421b","action_traces":[{"receiver":"eosio","context_free":false,"cpu_usage":0,"act":{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"},"console":"","region_id":0,"cycle_index":0,"data_access":[{"type":"write","code":"eosio","scope":"eosio.auth","sequence":0}],"_profiling_us":24}],"deferred_transaction_requests":[],"read_locks":[],"write_locks":[{"account":"eosio","scope":"eosio.auth"}],"cpu_usage":1000,"net_usage":364,"_profiling_us":49,"_setup_profiling_us":0}}


86e4c15a 0000 0200 19e992a7 00 00 00 00 01 0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100

8ae4c15a 0000 0a00 f150f7ca 21 e807 00 00 01 0000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100

-------------------

POST /v1/wallet/sign_transaction HTTP/1.1
Host: localhost:6667
User-Agent: Go-http-client/1.1
Content-Length: 679
Accept-Encoding: gzip
Connection: close

[{"expiration":"2018-04-02T08:06:30","region":0,"ref_block_num":2,"ref_block_prefix":2811423001,"net_usage_words":23,"kcpu_usage":2,"delay_sec":0,"actions":[{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100","name":"newaccount"}],"signatures":[],"context_free_data":[]},["EOS5qAuqXNnYSsofL32CEkexzr4qrHBP33BNsPZqPxhEzvJmze7DE"],"0000000000000000000000000000000000000000000000000000000000000000"]HTTP/1.1 201 Created
Content-Length: 677
Content-type: application/json
Server: WebSocket++/0.7.0

{"expiration":"2018-04-02T08:06:30","region":0,"ref_block_num":2,"ref_block_prefix":2811423001,"net_usage_words":23,"kcpu_usage":2,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":["EOSJxw7WFxgM5xTCfU79wkk8RCy1b8u5D5GWE4VX7fx8ThTB2tBLktXaZgoYeKD9xjX5n3jydvy2NJGuJAJu2ccnxBCXgJzny"],"context_free_data":[]}


POST /v1/wallet/sign_transaction HTTP/1.0
Host: localhost
content-length: 709
Accept: */*
Connection: close

[{"expiration":"2018-04-02T08:06:34","region":0,"ref_block_num":10,"ref_block_prefix":3405205745,"net_usage_words":33,"kcpu_usage":1000,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":[],"context_free_data":[]},["EOS5qAuqXNnYSsofL32CEkexzr4qrHBP33BNsPZqPxhEzvJmze7DE"],"0000000000000000000000000000000000000000000000000000000000000000"]HTTP/1.1 201 Created
Content-Length: 681
Content-type: application/json
Server: WebSocket++/0.7.0

{"expiration":"2018-04-02T08:06:34","region":0,"ref_block_num":10,"ref_block_prefix":3405205745,"net_usage_words":33,"kcpu_usage":1000,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":["EOSKfA5xarhtSop6v5s2Cse1QdBxTcGufvpzx7wZ4bNrp6s2yVU5gcRPzafyoEDKxEpzFyaQzFAYDUzDAqxMXx5jUc5z6tgEN"],"context_free_data":[]}

IN:
0000000000ea3055 000030c94c833055 01000000 01 00 02 3bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a 01000001000000010002 3bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a
01000001 00000000 01 0000000000ea3055 00000000a8ed3232 0100

0000000000ea3055 000030c94c833055 01000000 01 00 02 c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e 01000001000000010002 c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e
01000001 00000000 01 0000000000ea3055 00000000a8ed3232 0100


---------



POST /v1/chain/push_transaction HTTP/1.1
Host: localhost:8889
User-Agent: Go-http-client/1.1
Content-Length: 520
Accept-Encoding: gzip
Connection: close

{"signatures":["EOSJxw7WFxgM5xTCfU79wkk8RCy1b8u5D5GWE4VX7fx8ThTB2tBLktXaZgoYeKD9xjX5n3jydvy2NJGuJAJu2ccnxBCXgJzny"],"context_free_data":[],"compression":"none","data":"86e4c15a0000020019e992a700000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100"}HTTP/1.1 401 Unauthorized
Content-Length: 553
Content-type: application/json
Server: WebSocket++/0.7.0

{"code":401,"message":"UnAuthorized","error":{"code":3030002,"name":"tx_missing_sigs","what":"signatures do not satisfy declared authorizations","details":[{"message":"transaction declares authority '{\"actor\":\"eosio\",\"permission\":\"active\"}', but does not have signatures for it.","file":"chain_controller.cpp","line_number":972,"method":"check_authorization"},{"message":"","file":"chain_controller.cpp","line_number":346,"method":"_push_transaction"},{"message":"","file":"chain_controller.cpp","line_number":271,"method":"push_transaction"}]}}



POST /v1/chain/push_transaction HTTP/1.0
Host: localhost
content-length: 522
Accept: */*
Connection: close

{"signatures":["EOSKfA5xarhtSop6v5s2Cse1QdBxTcGufvpzx7wZ4bNrp6s2yVU5gcRPzafyoEDKxEpzFyaQzFAYDUzDAqxMXx5jUc5z6tgEN"],"context_free_data":[],"compression":"none","data":"8ae4c15a00000a00f150f7ca21e8070000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"}HTTP/1.1 202 Accepted
Content-Length: 955
Content-type: application/json
Server: WebSocket++/0.7.0

{"transaction_id":"3908e02af8f2b829b22daa5e4ff55bc1f83d301147667d56b03def2a9665421b","processed":{"status":"executed","id":"3908e02af8f2b829b22daa5e4ff55bc1f83d301147667d56b03def2a9665421b","action_traces":[{"receiver":"eosio","context_free":false,"cpu_usage":0,"act":{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"},"console":"","region_id":0,"cycle_index":0,"data_access":[{"type":"write","code":"eosio","scope":"eosio.auth","sequence":0}],"_profiling_us":24}],"deferred_transaction_requests":[],"read_locks":[],"write_locks":[{"account":"eosio","scope":"eosio.auth"}],"cpu_usage":1000,"net_usage":364,"_profiling_us":49,"_setup_profiling_us":0}}

86e4c15a 0000 0200 19e992a7 00 00 00 00 01 0000000000ea3055 00409e9a2264b89a 01 0000000000ea3055 00000000a8ed3232 7c0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100

8ae4c15a 0000 0a00 f150f7ca 21 e807 00 00 01 0000000000ea3055 00409e9a2264b89a 01 0000000000ea3055 00000000a8ed3232 7c0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100



------------


47010000 (len)
00 (msg type)
f8f6 network_version
0000000000000000000000000000000000000000000000000000000000000000 chain_id
c1a8e76747f237675e0b663f6f589379689d6006fc02bcd3436f616e7c0f2049 node_id
00000000000000000000000000000000000000000000000000000000000000000000a23fcae11ca921150000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 public_key, time, token(sha256), signature
16302e302e302e303a39383736202d2063316138653736 (p2p address) "0.0.0.0:9876 - c1a8e76"
50030000 last irreversible block
0000035033892147f50cb068b842457e3754c245d7da2149f5e632cff685efce5103000000000351fba35fca451d33228041e9e2c42aab0e7a9424cc07006524e4658532056c696e75781022454f532054657374204167656e74220100


   struct handshake_message {
      int16_t                    network_version = 0; ///< derived from git commit hash, not sequential
      chain_id_type              chain_id; ///< used to identify chain
      fc::sha256                 node_id; ///< used to identify peers and prevent self-connect
      chain::public_key_type     key; ///< authentication key; may be a producer or peer key, or empty
      tstamp                     time;
      fc::sha256                 token; ///< digest of time to prove we own the private key of the key above
      chain::signature_type      sig; ///< signature for the digest
      string                     p2p_address;
      uint32_t                   last_irreversible_block_num = 0;
      block_id_type              last_irreversible_block_id;
      uint32_t                   head_num = 0;
      block_id_type              head_id;
      string                     os;
      string                     agent;
      int16_t                    generation;
   };

Reply:

25000000 (len)
01 (msgtype)
05000000 (enum go away reason)
0000000000000000000000000000000000000000000000000000000000000000 (chain_id)




------------------------------


9e96c25a00000300d847cf541780100000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c8330550100000001
9e96c25a00000300d847cf541f80100000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c8330550100000001

00023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f
00023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f

953f9a0100000100000000010000000000ea305500000000a8ed323201000000
953f9a0100000100000000010000000000ea305500000000a8ed32320100



-----

Mine to get_required_keys

{"expiration":"2018-04-02T20:52:30","region":0,"ref_block_num":2,"ref_block_prefix":837290549,"net_usage_words":0,"kcpu_usage":0,"delay_sec":0,"actions":[{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100","name":"newaccount"}]}


cleos to get_required_keys

{"transaction":{"expiration":"2018-04-02T20:52:39","region":0,"ref_block_num":19,"ref_block_prefix":153469732,"net_usage_words":0,"kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"}]}

0000000000ea3055000030c94c833055 0100000001000 23bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a  0100000100000001000  23bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a  0100000100000000010000000000ea3055 00000000a8ed32320100

0000000000ea3055000030c94c833055 0100000001000 2c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e  0100000100000001000  2c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e  0100000100000000010000000000ea3055 00000000a8ed32320100

required: EOS616K1RvJcyPikmu6yJxKfQdrnPBgfYBdTZXuV2qZPwXicYgrsT

----

US:
POST /v1/wallet/sign_transaction HTTP/1.1
Host: localhost:6667
User-Agent: Go-http-client/1.1
Content-Length: 681
Accept-Encoding: gzip
Connection: close

[{"expiration":"2018-04-02T20:52:30","region":0,"ref_block_num":2,"ref_block_prefix":837290549,"net_usage_words":23,"kcpu_usage":2048,"delay_sec":0,"actions":[{"account":"eosio","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100","name":"newaccount"}],"signatures":[],"context_free_data":[]},["EOS616K1RvJcyPikmu6yJxKfQdrnPBgfYBdTZXuV2qZPwXicYgrsT"],"0000000000000000000000000000000000000000000000000000000000000000"]HTTP/1.1 201 Created
Content-Length: 679
Content-type: application/json
Server: WebSocket++/0.7.0

{"expiration":"2018-04-02T20:52:30","region":0,"ref_block_num":2,"ref_block_prefix":837290549,"net_usage_words":23,"kcpu_usage":2048,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":["EOSK5NbevJaECybrWQcHotbyFsYJfYJPcZtqxnxN38baN9GXNn9MFtr2RK9t9yyogpyjkZQ1HsCtXurkR6q5Z9xwHuPaPeyWD"],"context_free_data":[]}

POST /v1/wallet/sign_transaction HTTP/1.0
Host: localhost
content-length: 708
Accept: */*
Connection: close

[{"expiration":"2018-04-02T20:52:39","region":0,"ref_block_num":19,"ref_block_prefix":153469732,"net_usage_words":33,"kcpu_usage":1000,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":[],"context_free_data":[]},["EOS616K1RvJcyPikmu6yJxKfQdrnPBgfYBdTZXuV2qZPwXicYgrsT"],"0000000000000000000000000000000000000000000000000000000000000000"]HTTP/1.1 201 Created
Content-Length: 680
Content-type: application/json
Server: WebSocket++/0.7.0

{"expiration":"2018-04-02T20:52:39","region":0,"ref_block_num":19,"ref_block_prefix":153469732,"net_usage_words":33,"kcpu_usage":1000,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":["EOSKjwLgpMmJfGiVaf3x8SMjkBHkNFTsLBbiYgcXtCZ2eDv316noccL3WEst87jFUaLEKyaofzVGyQjgNhs3sRxEbUS9vbq5x"],"context_free_data":[]}



--------------------

POST /v1/chain/push_transaction HTTP/1.1
Host: localhost:8889
User-Agent: Go-http-client/1.1
Content-Length: 522
Accept-Encoding: gzip
Connection: close

{"signatures":["EOSK5NbevJaECybrWQcHotbyFsYJfYJPcZtqxnxN38baN9GXNn9MFtr2RK9t9yyogpyjkZQ1HsCtXurkR6q5Z9xwHuPaPeyWD"],"context_free_data":[],"compression":"none","data":"0e98c25a00000200350ae8311f80100000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c833055010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a010000010000000100023bf0afb1a36116a70276d69920d4b8a50c039af08aafe6e096be61328f953f9a0100000100000000010000000000ea305500000000a8ed32320100"}HTTP/1.1 401 Unauthorized
Content-Length: 553
Content-type: application/json
Server: WebSocket++/0.7.0

{"code":401,"message":"UnAuthorized","error":{"code":3030002,"name":"tx_missing_sigs","what":"signatures do not satisfy declared authorizations","details":[{"message":"transaction declares authority '{\"actor\":\"eosio\",\"permission\":\"active\"}', but does not have signatures for it.","file":"chain_controller.cpp","line_number":972,"method":"check_authorization"},{"message":"","file":"chain_controller.cpp","line_number":346,"method":"_push_transaction"},{"message":"","file":"chain_controller.cpp","line_number":271,"method":"push_transaction"}]}}

---

POST /v1/chain/push_transaction HTTP/1.0
Host: localhost
content-length: 522
Accept: */*
Connection: close

{"signatures":["EOSKjwLgpMmJfGiVaf3x8SMjkBHkNFTsLBbiYgcXtCZ2eDv316noccL3WEst87jFUaLEKyaofzVGyQjgNhs3sRxEbUS9vbq5x"],"context_free_data":[],"compression":"none","data":"1798c25a0000130024c3250921e8070000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"}HTTP/1.1 202 Accepted
Content-Length: 956
Content-type: application/json
Server: WebSocket++/0.7.0

{"transaction_id":"92d851839a79b615ecd78fc4f4ef9797ba814cf71d4ba7741000fd8696fb6535","processed":{"status":"executed","id":"92d851839a79b615ecd78fc4f4ef9797ba814cf71d4ba7741000fd8696fb6535","action_traces":[{"receiver":"eosio","context_free":false,"cpu_usage":0,"act":{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea3055000030c94c83305501000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e01000001000000010002c3fd81035735eb1685705524a59f6dd4f1c799012736d219f790ad5b7641ba4e0100000100000000010000000000ea305500000000a8ed32320100"},"console":"","region_id":0,"cycle_index":0,"data_access":[{"type":"write","code":"eosio","scope":"eosio.auth","sequence":0}],"_profiling_us":91}],"deferred_transaction_requests":[],"read_locks":[],"write_locks":[{"account":"eosio","scope":"eosio.auth"}],"cpu_usage":1000,"net_usage":364,"_profiling_us":213,"_setup_profiling_us":0}}

--------------

Pushing ABI for eosio.msig

0000735802ea3055010c6163636f756e745f6e616d65046e616d65050770726f706f736500040870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65037472780b7472616e73616374696f6e09726571756573746564127065726d697373696f6e5f6c6576656c5b5d07617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c09756e617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c0663616e63656c00030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650863616e63656c65720c6163636f756e745f6e616d65046578656300030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650863616e63656c65720865786563757465720500000040615ae9ad0770726f706f7365000000406d7a6b3507617070726f76650000509bde5acdd409756e617070726f7665000000004485a6410663616e63656c0000000000


0000735802ea3055010c6163636f756e745f6e616d65046e616d65050770726f706f736500040870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65037472780b7472616e73616374696f6e09726571756573746564127065726d697373696f6e5f6c6576656c5b5d07617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c09756e617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c0663616e63656c00030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650863616e63656c65720c6163636f756e745f6e616d65046578656300030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650863616e63656c65720865786563757465720500000040615ae9ad0770726f706f7365000000406d7a6b3507617070726f76650000509bde5acdd409756e617070726f7665000000004485a6410663616e63656c0000000000

8054570465786563010870726f706f73616c03693634010d70726f706f73616c5f6e616d6501046e616d650870726f706f73616c

805457046578656301000000d1605ae9ad03693634010d70726f706f73616c5f6e616d6501046e616d650870726f706f73616c

0000735802ea3055010c6163636f756e745f6e616d65046e616d65050770726f706f736500040870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65037472780b7472616e73616374696f6e09726571756573746564127065726d697373696f6e5f6c6576656c5b5d07617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c09756e617070726f766500030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d65056c6576656c107065726d697373696f6e5f6c6576656c0663616e63656c00030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650863616e63656c65720c6163636f756e745f6e616d65046578656300030870726f706f7365720c6163636f756e745f6e616d650d70726f706f73616c5f6e616d65046e616d650863616e63656c65720865786563757465720500000040615ae9ad0770726f706f7365000000406d7a6b3507617070726f76650000509bde5acdd409756e617070726f7665000000004485a6410663616e63656c0000000000

805457046578656301000000d1605ae9ad03693634010d70726f706f73616c5f6e616d6501046e616d650870726f706f73616c




------------------------


$ ec transfer eoscanada eosarctic 1000


POST /v1/wallet/sign_transaction HTTP/1.0
Host: localhost
content-length: 534
Accept: */*
Connection: close

[{"expiration":"2018-04-06T19:17:18","region":0,"ref_block_num":350,"ref_block_prefix":2277385774,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"transfer","authorization":[{"actor":"eoscanada","permission":"active"}],"data":"000030c94c8330550000402ea36b3055e80300000000000004454f530000000000"}],"signatures":[],"context_free_data":[]},["EOS8NijGLHT8WyDmt2nqMwfP1hr8EiYx5JCYBWSP9S26WgbeugvSJ"],"0000000000000000000000000000000000000000000000000000000000000000"]HTTP/1.1 201 Created
Content-Length: 506
Content-type: application/json
Server: WebSocket++/0.7.0

{"expiration":"2018-04-06T19:17:18","region":0,"ref_block_num":350,"ref_block_prefix":2277385774,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"transfer","authorization":[{"actor":"eoscanada","permission":"active"}],"data":"000030c94c8330550000402ea36b3055e80300000000000004454f530000000000"}],"signatures":["EOSKj6qNjC2K75hj81P8RVuM6wyA3fP8SD1URafucMS7LfX7XwDVDNb8jWxo21jixNJqr7UaYBrB3wERJXxUNnAvs9jyTAc7j"],"context_free_data":[]}


--- Present in the wallet:

    "EOS8NijGLHT8WyDmt2nqMwfP1hr8EiYx5JCYBWSP9S26WgbeugvSJ",
    "5KWBFG1co7XsaKCHBND9r5RdVDUqBWdPsfqy3xxpZiRSs5kHgF4"

---

POST /v1/chain/push_transaction HTTP/1.0
Host: localhost
content-length: 351
Accept: */*
Connection: close

{"signatures":["EOSKj6qNjC2K75hj81P8RVuM6wyA3fP8SD1URafucMS7LfX7XwDVDNb8jWxo21jixNJqr7UaYBrB3wERJXxUNnAvs9jyTAc7j"],"compression":"none","packed_context_free_data":"","packed_trx":"bec7c75a00005e012e26be8700000000010000000000ea3055000000572d3ccdcd01000030c94c83305500000000a8ed323221000030c94c8330550000402ea36b3055e80300000000000004454f530000000000"}HTTP/1.1 400 Bad Request
Content-Length: 858
Content-type: application/json
Server: WebSocket++/0.7.0

{"code":400,"message":"Bad Request","error":{"code":3030000,"name":"transaction_exception","what":"transaction validation exception","details":[{"message":"condition: assertion failed: integer underflow subtracting token balance","file":"wasm_interface.cpp","line_number":805,"method":"eosio_assert"},{"message":"","file":"apply_context.cpp","line_number":30,"method":"exec_one"},{"message":"","file":"chain_controller.cpp","line_number":1989,"method":"__apply_transaction"},{"message":"","file":"chain_controller.cpp","line_number":2024,"method":"_apply_transaction"},{"message":"","file":"chain_controller.cpp","line_number":2225,"method":"wrap_transaction_processing"},{"message":"","file":"chain_controller.cpp","line_number":353,"method":"_push_transaction"},{"message":"","file":"chain_controller.cpp","line_number":328,"method":"_push_transaction"}]}}
