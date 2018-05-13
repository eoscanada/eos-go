
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


43 01 00 00 00 57 63 00  00 00 00 00 00 00 00 00   C....Wc. ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 b8  13 2b c7 23 0f db ae 71   ........ .+.#...q
ff 7d ba 61 ce f2 bb 00  2e 60 b4 9d 09 d8 49 70   .}.a.... .`....Ip
c6 ab 39 ca 56 3c ee 00  00 00 00 00 00 00 00 00   ..9.V<.. ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 00  00 9c 91 4c 03 0d 03 26   ........ ...L...&
15 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 16 30 2e 30 2e  30 2e 30 3a 39 38 37 36   ....0.0. 0.0:9876
20 2d 20 62 38 31 33 32  62 63 00 00 00 00 00 00    - b8132 bc......
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00   ........ ........
00 00 05 6c 69 6e 75 78  0c 22 45 4f 53 20 43 61   ...linux ."EOS Ca
6e 61 64 61 22 01 00

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









ff6ed75a00009906975014d600000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

ff6ed75a00009906975014d600000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100










DIALOG WITH MY WALLET:

POST /v1/wallet/get_public_keys HTTP/1.0
Host: localhost
content-length: 0
Accept: */*
Connection: close

HTTP/1.0 200 OK
Date: Wed, 18 Apr 2018 17:57:21 GMT
Content-Length: 226
Content-Type: text/plain; charset=utf-8

["EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV","EOS5GNc1NNsChC2URSevTBYhvQxGzjAgFRDi8V992ckjyc5tRzWn4","EOS5Dg9cu3yn5cMpWkrZnhmYk2xDBWmu62Sj2dNrWn6Ui82eoYJQh","EOS71W8hvF43Eq6GQBRhuc5mvWKtknxzmb9NzNwPGpcEm2xAZaG8c"]


POST /v1/wallet/sign_transaction HTTP/1.0
Host: localhost
content-length: 716
Accept: */*
Connection: close

[{"expiration":"2018-04-18T17:57:49","region":0,"ref_block_num":44982,"ref_block_prefix":3032689784,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500002059b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":[],"context_free_data":[]},["EOS5Dg9cu3yn5cMpWkrZnhmYk2xDBWmu62Sj2dNrWn6Ui82eoYJQh"],"0000000000000000000000000000000000000000000000000000000000000000"]HTTP/1.0 201 Created
Date: Wed, 18 Apr 2018 17:57:21 GMT
Content-Length: 415
Content-Type: text/plain; charset=utf-8

{"expiration":"2018-04-18T17:57:49","region":0,"ref_block_num":44982,"ref_block_prefix":3032689784,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":""}],"signatures":["EOSK4CLky4ytx7eTuZJx885CmcgvJwr5oxBCZFzSHbh4uRkXwDGR1paVBtWNgCXxMyTxS5P7pYRSH3K4u5DLsN4HeZpmb7H1F"],"context_free_data":[]}



---- THEN WITH THE KEOSD WALLET:


POST /v1/wallet/get_public_keys HTTP/1.0
Host: localhost
content-length: 0
Accept: */*
Connection: close

HTTP/1.1 200 OK
Content-Length: 14505
Content-type: application/json
Server: WebSocket++/0.7.0

["EOS4uoxMcUkU1MwVzAxi4oKJ5TmxQCXHqz7kruNWXYmAuGgpx7Ez3","EOS4v11rpA2SE7R535QkEoqPPgsNLPdDQAMYvPaBghA7NPwAfkici","EOS4wny9CMBPVstzuRrpke2D65P431dbmsfTX3F9JK6c6csKQ8ei3","EOS4wsTrgioTb4KXCjQPE7iLpACZdy3ANaotjVXHteGyu23TetKJk","EOS4xYKUo5sE9zXy4H1GVieuxVHqsN5MNai42mibMWPXdnTm1Chsk","EOS4zdh1dQcYmVCqnzqqxMGQGLUkZ37bsKHph4Rrh4fMjm97CXJrT","EOS4ziUeQ9AYNWhdEtFnFR2CVdPauBPwMRCL2i3NEyb9FCGFKjBiW","EOS51YeKU2nZLjjeyyDK4NLCmJpgk9xZrenAW4ir6oVa6VYbM14DL","EOS51hCxBK26MoJKWe9Rnzv12Zu1P7ix9vYW4VtkqhoNvZ7qrfkDL","EOS52SDHppWGxQKrJKoiGY1HmvHjqL1aJVcuEpvQuVwcbHxtyKLBX","EOS53CYn2kpPmdau2HrXihbzt3ww1Asfhm4ozduTZwuWFWf6V8iNW","EOS54asM7pJgqHYC7Nyb9SWuWBG7zD8pA8xJvoZPcNVkCu3s1g4Mc","EOS55mwjhh5QzHfcGPvLUnQGfVETiHzUU1ULnQfSc7JokQHe5YQM6","EOS59obWS7geJPfH8cba9yyxxgSXJxsKYQKhPz1uN8iCudNKPBoJz","EOS5BEFpXxvDeBoRpdm52D3Gdp7R67AL13dMA2ktsg4FS39bU9tVU","EOS5Cu9zPdmgK3kgdrZ2x6UdpL788hMPhRBzL8R4p3gYz6WH6X295","EOS5D4j2ConaVjcrzzJpfajEU9e5G29sTeVMBacD177Zr6hCLHKXS","EOS5DCo66YkATu5UmJ1sLH7cXSVPk5tur8ZpHpTFcvHWXdwN5mDpk","EOS5DMzZkHxCRXDQgUR3n5CcT7z863cyBHogPC8ckuKArjfzjAtCi","EOS5DRyMC6EPQeyTV4NUSwQCYarbp1DqDjq46uwZvS9siWbuGGyDE","EOS5Dg9cu3yn5cMpWkrZnhmYk2xDBWmu62Sj2dNrWn6Ui82eoYJQh","EOS5GNc1NNsChC2URSevTBYhvQxGzjAgFRDi8V992ckjyc5tRzWn4","EOS5GzchCUNMW44W8rdQ7XCmQ5gsXXWiN8z9W6QDZp8b2bPCTNeAa","EOS5H6CsmkZVisBXR5CmPpeGQCp3k41PLMf47FyGtDhMrHwkkred4","EOS5KTvQ1Bg9GArqZfTGz4ozGVCFLAYtXQp9uodDyf31oY7cdBmpt","EOS5L2kXb37KXZoMUdvfRuuk49VQVbsayPKP7ZMg73PokpTKrnqza","EOS5LKiWbuhpRELgmb5TmDaEpK4ckgyTuTBysPe1eNjTZtEUtMxDG","EOS5LZmjf6qTRhm1Y1FhF3vb57S1G9whvu8vSS17pvcsRdzkR4Hqq","EOS5M1PC5Xnm4J8oyvFMiCcLc5Gu5xx7pdwE1oGLkScWGDCnmJfif","EOS5MBb8k9ByNUmSz46MfJQ9nUiy6ZZPmTzgsst5e7muznm6FuRSZ","EOS5MDruMhkJLFcPttqNRrLM915C1Pusfgm6AkTsvL5YEZMrbAhtM","EOS5NrR8GXWJqGT2HienFFk2FSkc2v1KuuqmLMCHzxBWNgG6Bx9zL","EOS5PCAWPJ9ARnD7YwZVPN6NDGqYcF31k3ckn9EuGRCkLLdnArRZY","EOS5Pk4M1moyKEcYTQFMztPf4YK2Ej5ZnTSqTqN7wRzbj22jdVEns","EOS5UScwSaYssyJewxnsrDeu9sZ1t9UD68aShXUj91bHHZguneP3o","EOS5UU3KDDRnbFKBXSJzVGnCB3JwZDnia9VwytyWfLphwKMX471F1","EOS5VFJLTs28RVSRBBfUzK9Nm7VqJJ2Ft7kxvFWrjzGAgWuTSfryT","EOS5VRefibP4YmhY4Vx2M7VTD11t49HrCz3jdHjN7Sjq5PE4N8KcU","EOS5W6fSAEHQ4dsmVBQxBM2jibEhFZ2ekceBVLAFSTqHnKvvRBCrb","EOS5Wwepv5soTb3v3KWnuSXieufwjW9yUd8aAU9vWDtJchzdDQDFX","EOS5XJB8gLrhCUQ5JxvaQPRdBgrAXuYUfSyDyrnWzY645Fkv4UdVT","EOS5YYKCCeLfimDbueNgjgxDqTva8w2kJkRZCzZgBbh6BTF6bVDtT","EOS5ZtYjCTfy39tnAZBe895QknNKa4NgVfnVFjXwkdtA3Lth9u2av","EOS5amjQAQN4CsQc2Y4w79cewTuoEWQ2uMDdtaAdgmaEvSw27Jqdc","EOS5b67aLeuDGV1eQQ1PPSxHYczBgYVmpsiqu7HazT8iCGkuJx47u","EOS5dBY6WUYbox6hWZRrxvY92LMbx4U6gFv7TgMjJSiP27XjPyZ37","EOS5e4ZsJZ6oU6rtZcg3AW1Kz4pE6GdnPJJLwiHUe8KPfKMjbrbLc","EOS5e9DUXmmZCaYZhZjbX6t7LmKWRHqHATjx998s3yv5NQuJZHNLi","EOS5eo5GcPog4ofe9ev1wt37PyRdbjUM7SetRL7Gprbh9gNGzZVvQ","EOS5fbnnxWJEX8Bv8obLViaYnxGYUitFGdyqbMWfFyshXYBR99S1g","EOS5fsBqDwWyTqtZQnMmpdmrSkJpfFhQkrzjYGPD19Vt1hT2cMiCG","EOS5hqQqg5Z2S8tzE2zYcmeeuw2UgTVfxFLRZFp62SwaDwYMw5mit","EOS5hqv435sA1sgDk9MWQAxEyc69ZocQ3AByg1L2bUE1YirynHiDo","EOS5iSGBJZBq4c4gie5pFyqknJiukxer3CJjRo6AdbkiV5JfLgSoT","EOS5j4EbH8vR6jTqyKxEzj6YaRw88A5sbVam7peLC8CWGYifiFLX2","EOS5jK4krXq6x1uCiCzFjqGRvN6XneEaRDj5hWneBo1V5Wm9p2mkp","EOS5mT1gbPoxXD4Vrce2hZnyTWNrRffo9r9cDYrGyjoeP125A3gMP","EOS5ntzPJDJjvSHadssSKLkUMicmeEv8nDo9ENrzE7dbn83zw2k86","EOS5oZnmQX8BgYXFwVXvrmaFcr9GqRmJtfr4PY8vnCavCLvLp1Z4w","EOS5obyGC85LSBGuJAtsjheEUpFGABb1NsxcC1iyApUQtsigtdcfR","EOS5oi61Czy7txNNFjuhJhK1aTfFF3jpBMiNxUpv1q8BSCBk8ShYF","EOS5qAuqXNnYSsofL32CEkexzr4qrHBP33BNsPZqPxhEzvJmze7DE","EOS5rhGmpoWmfYwfQtSAnsx82yZaBHP4pVi2tBFJbm1o7n6cjg8ZY","EOS5sWBa9iArpWRSEyGSr4riviaMeMduJgHHxzsrZQY5S6CfBNTLV","EOS5t3Dqt5Ncor5LXiHhLcyUzMT9R4zDfEyy2bKDLQDYKdVGcViBY","EOS5tk4iFjS1vXLfpk3ArDmx8DuH9BqvMkBfcJJ87Q68ZPRHF8gxL","EOS5txcKD8diHt1vPX7nftE2GxVp7PPimfYYEnjeK9Zd5rDsrmgEd","EOS5w2skg5dxhC2sLGquBAzA1FXFCDekxoaudo2ocSueGzevM7Y2C","EOS5wZXDN1UfDPtbKKvFTaBAdnBijsGq1ZgtdMhpoXwphgFDM4tEd","EOS5xWjz6WrtwWt3XHpzZSvRMvu4xnCH2c8ywEJdiXupr7uJpupEv","EOS5xpD2VNdPnYsYr3RcFHNnuoydao3aVva1VvbWe2s8ktPUKrjjc","EOS5z3rCH3p9pEfaJgs9kgTGCwpNBKxoozbcVaEZ7kxwdWV7mSuQN","EOS5zQRkZGu2W9VGbvCEiapWhUbyNHB6xNaBa39F3txduKMPKK8eR","EOS5zsfWa1YAw71mCj2YARuQmoEMbWKxxjMoHAqmtcZz2F76UPsdt","EOS612etfrkFNuMEnzNNzEA5FiuMSA4v7jdndFSL5MFabC7AWKHA1","EOS616K1RvJcyPikmu6yJxKfQdrnPBgfYBdTZXuV2qZPwXicYgrsT","EOS61R4F55XXbpVnxgJxjpwxqAMdiUt2csUMTuw54jtuRsF8ixkZU","EOS64CLdikZTe2ez2Wyz2UqwKES6vLB9JxbBvrAJoMazBEqiR1RF5","EOS65MTqP8qrFQwuAEGAkRTyP9bg1hw9Gz5vNTURJAbj677pLwdEt","EOS65aVawnHMddsaF1EvxGQ1rH1QWqNYBBfKkRGjFwowe5mCPHZmR","EOS66QkdxKvHyVSWTj1AmuZW1YFqHwmi28BZ8FJxThJnoFE6nKQLU","EOS6736GZZtpeRdM8Ajip8mRGoqKQepdigB3STEE3WXkVFSBjcfNW","EOS67CmSjvjddmoStvrWCwuXVVSTPAHp5eggQeWfXMuogDmtnJU2e","EOS67vp4oZw72WuvETpy9qfghPUhKe9PpbhA2nLf5vwa4TUDwTYK3","EOS683CCFWHX2KeaVasHwj1S1bYMQXvSCjrGcuN32ieK5mTgVmGZV","EOS6B97ZuFRMVictWsc5cyv8aA38CmHbiASzeRmovQG9RqJkNEaiu","EOS6Bm5SupBXAF4dZLDV1tXSzjixb9hFMhZBoWT3FLJpjSD1G2k4A","EOS6CeM8t3bZRQvh8RuSrsH2ZmVwqZPnXdjxM84YUaST2NNVkyzqY","EOS6DAqmrmaMVR7FJfnLW8PNvahSTXGBcZn6c92s6N7RfV3DW3WdD","EOS6DHGjZQUaGSVMaVPBuJznFZ45uidmompFfCVMr6C3rcyyBn8ep","EOS6DnJnjLTcUVqHXfKjwLA6DwTq42dgdkDAX4HguZVhUxnBEwQQx","EOS6E5oofyCzPkC1vXKHiyC5RopDEZzkBUxwhquuKAh9BBQ6bJywZ","EOS6EkCBuDGCaSQRKcHX5tgF7j8udFXHmeqpby6P6GcLsxzsqcpWt","EOS6GyMADS2NPqdX3s3Mpg6DjBV7aCXGpy26Ksia36oz2ViePqJUe","EOS6HRWrv32A8k1YqT7pFMYHbW4wqUbDKoMrLdAWbyNC1aAtREaMJ","EOS6KTT7kAgqmx8msCobZkhdY19mKa7UqYvPiXWTcqSXfH7JvPxtC","EOS6LUrPexTfBkx1zPxVsmzhSAdsgbJ32PTyZet4xZHPAtX6hgoDd","EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV","EOS6NofKdArTrXfeGzhuEhew4h8WaL7xSTdCGQFie6yNbVDfjdHzh","EOS6Q4yZe75SWhYWs2hPnrQup53JmyYcy3CR3hPBKfAYrFRXzryGP","EOS6T93qdBRofDpuevtnrTevJvMGxXbnwPak3HPuX5D2F3o9u7tV9","EOS6VUYT7vXziYV4GHuCmXaS6MgzZCD1QBc6qGZVZRCFAWCABiX5x","EOS6XXQbLPLmVZW1h1T5ZfUaYg7NbBBgXLP45vs78bhfz8EEyahNU","EOS6XsteFYyrwuchezg4DVZ4K17LQkhMLSajkURv9Pqp8Sq6Xm2Av","EOS6YeeAc3oTrV5x6D97oHNjrK4xygx3g7wjbKDRoSsHjvGWiQ1VB","EOS6Yng8VUiiN5P5wpSxQzFJiRBEWcrXbijNwQ9VtnXm6M8x7taW1","EOS6aLRp8Zce7NNpiYkktnnng8ucrpaYsb5vBzVbrGpYr2Eyue8xf","EOS6aUydKqq28G1wNYLvVAwe2xSuYU8mWH2VcnmGQ7iPLQQ31xW7Q","EOS6aYyK9ceNC3HzdzxrytqTuBsG4iErteViQTnUWGUcZYqCcCAum","EOS6b2Khw87uvgRHbocnui3wdyMiqQvnxAMjBW3QsBWNu4xAMku3c","EOS6cNXawwqc4pVDNodugLa34cRATbhfxhmiosejjHJKeXEsDXgqj","EOS6ciUrDQnZXqeoHEZdgw2Szi1MpRcHD4wBjJdPEug7MN8s5DxuB","EOS6dMAbmqhpCMBdMdxr88iAeg8vhTNbZAoiwzBzv4LLsWzsJr7vc","EOS6eGQKqA19SQna2o45hxMnUopCLdtY2ZbQNM88bsDj6KaPizjz2","EOS6eswNGePdSka8Wi84rCSvWUExCGCsTwVxLAWD4xqoWmZFW8uTq","EOS6ewf6KPg57X3YBZzbqDewbVijbGsgij4D1RFA7x4x3MrEw99kk","EOS6fYgG9xRcWeUYt8xLfnuAbnYk5H3qgozjvLNgV3L2HjZMfxhX7","EOS6fqhEjVfud6k1PxVDVcdi1LBnkJ5fZQLDCmYjyGnuyyktTkwZz","EOS6g46KBfXfDQXMcT78nLGqR68HKf8zVGD8Z7RoxnJS7v3ykrzgC","EOS6hN3faT6Bwe61YxxuvXa1yopazoJbiLpntT7FnD2vrusFbTd3g","EOS6j27wRtjJZVDnwiQsZJNcw9rLNdh2z2Nq4zJ14vK6LeQAi8Bz4","EOS6j4XH2BJ5RohivbCVWibDCf4EtBp9wY53nhP1Af5FAydXqFwjx","EOS6jYZWY15XuBquxGXcB6ymjGca9B3LijwMwnsMFxwyT27BJnd96","EOS6jcqMV3Cr5qFsHmmuQoNQhokqzMpiCbNvVzXsz45C4vuxh1zYs","EOS6kHLn5nMgHzGtncNnsLecSS5MUXFiFNH4WBguarH5FiWkaTBS2","EOS6m2d6S1J39PTGgYeqpVs5L9av5TaruXW11vkuUYAxmDgs6xAhS","EOS6nFBjisGteCda9DrbqwpaXVGKVMbdSJastT8ip1yAFog9qLpPi","EOS6oAaHLWPQyAUWyBy22D1ATtF4M4Az9HUMA2dzxx3i17ByuS6Gz","EOS6oMiwVDEPABnzLPcJnHLvM47SLcBjBYauCSQtkN2X5kBaCJPNo","EOS6q1qCRcWGBPSv6n66CTxSBNBKVYT9Ay32dmJ2MwccfpGqYeCo9","EOS6q9Nx1x2gRYrtNCdaiqL7MuX5SBixqoemWvZeeeNfT6vUMRzAc","EOS6r27fY6H46gku5gYgc2pSJ7xZH2B9VskfEDR82Cxeuk8ZmARiw","EOS6rz6VGfKrpGXTy5CNgbYAycJjvov644dbqTsWCoi8Km3taDtdv","EOS6v561BuqGkFvF3seS8bigFevVqbmhdzELQNg12perLbpw9aQax","EOS6vBFLTVbKeNdNj2HBPad8QeYKaBbWe5ZiazxFk7tx5BXs3YibJ","EOS6yyCxWK1vLr3UrLKsMTseUEbmEosWNpceLY3uaeT4oQpHASuZN","EOS6zDV3M4NFRw96hsKFpVBPjbCnx5HwV5xNVRua8gpjDk2vUP7BV","EOS71W8hvF43Eq6GQBRhuc5mvWKtknxzmb9NzNwPGpcEm2xAZaG8c","EOS71fc8ZjtxNfEfMkinqwVqAVPxXSu4mUFkh3dLqWCFDk7PfvBmn","EOS71yq9R8Nqsaxw4EZHti57tNyG9KDurwz2zVK7ohA5aBiA2d2Db","EOS74fjF3vTa2HhcibARCCg3NZGsiUqYNSXyFb11NLeY1aaYWxdoh","EOS762m257osRiz1SySmo1oDRbTB1W2SEEiv6vvip1UeaEHcnbRy1","EOS76bKNdKT494wpAVJWqp9EbfcRAke8hVWELAAueY78AiJWz7M33","EOS77dZpx49fo2suPq5goDZA9kRn2FESniU8ajcocdLKncxtfJVjj","EOS78YmC9noFKG1YvpduxSjgQuZGcgbP5oFRbx4yAvr6joryX4q5b","EOS78roAiQ6z4Hg3YoQ96t35vJG4Wa34kL2cUxJcnaH7hi4Fnb2aN","EOS79fvjSFGv4NxeK8N5iCpt9meB1shhdiBDPHJZSddbx1XabKGxp","EOS7A7juZTX2AsH2H6rsWFzmGf7KaJfq4HNPfQ6rjgzSw2Jpq9vYC","EOS7A9nFhaC4iCLEPfGJoYQ55WNXwkjf5mM9aFZ4ir7qSfAPQqV8X","EOS7B7wtYpT1Frs4FgqiyobzZHdj8ZGQtTbcMzK2WidQD9owhpUAD","EOS7B9HsAA8GZHZkhdMJ6E9eGXwJXYQQxv7c7h7SkDcBHzuRWDh5j","EOS7BeZcwLiNKTeBroDbbEqnQSdbrrFU44KcMHfh9EJEMZfSpZVP9","EOS7CVRX64Jk8aqQXn5VCf6K9sdULxUpq7NhVatMxTaRA9xZ4XdZk","EOS7D7iC8qAX9Lfnrgymdu8NZkN2ucVHV7GS7fLLBnVyXtDZ1Mh5X","EOS7DUv5jRMjD5sFHkUgFDQYCsSkJP2yS23pp28QU3g244iJjjsvJ","EOS7Dta2zf4btz5rgjLhpDvEQbQWPJMJGgbYM3iYke7u8dkUAHUAv","EOS7EQaSsCrDV9R66qpdWHhnySh5z3PxM9NwisH9Ue1H4Z1CqkpY2","EOS7FSK2DuHv4aBkrPpshe1GWewVVsXBxUfRfjBMMuqRW6Rz151wb","EOS7FeCKGPz96Z3eYFJpT5E5SMiQFWkuvwRCKWqZtDtGZoHLam2ED","EOS7FgWfdR3yGA64peZvWeX1hm5hW6rhL6mF16zxpeMNQhKhsmRbm","EOS7H1S2uDUfPqxyh4dGpUNMKaJUaytM8otpY3Dqt1uRRzZxmZFuu","EOS7H2iJcDhHNpsp4vaSoTivnmwPJohxjvC6Nko43gUi7oCRPspc9","EOS7JzTuvhSm5rs56gKqH59thfnQUZw7KBboG2vhoU2yVHaVbCPsV","EOS7KuCrZUXWtoNU8KyjrwEuvd19o55DiSGK7eYFuHezrqKKbMhvb","EOS7LmMUYknLYqtzq9moyAdp2QbnWFv3m27PMUnfjdGEcbosSEMzF","EOS7MZw87T2p5M2nBTdtF25Ss8KoovatykbMdA3dFxEPoBcYALYza","EOS7NCN4vAC479WATVJRruKc8DUcJj9pnQH9x6oQdGALxCwGG77mT","EOS7NZ7LD4Po6FCTFBSUgXMZvreBPuAD8oJFJkkgZHHrMf5jKjTiF","EOS7Px7kWRVUkryfU2icqiS2iX33bR3R9KdzwgnhFf4MJAbgGhY24","EOS7SBRsREcrp6q5BtGegFWETxN7SZaTLPyH5XE4vgXHgY3q3KzMo","EOS7UG1XMafvGrFrAHTZP9KtDgASk4aSmJP3TWj787BYyC7MLCGRd","EOS7Ue5ZCFaDXruWimSfkeB9QGopneoUW2MzDyC9eYxEcFiVorPZD","EOS7VqHYsPcaoD9WixvPyC9cVc581tZ6BsUY5VNjFUo1T47nr5FNm","EOS7W3MmjrBvN5vpsjV2ag7n7oubR9PFZF4DGFJCLL6vH3TKqTNV9","EOS7Wv5ZANymRSkVnRQZJ8SNFn6ozFu7wy23oTRyth8i4o9i98y7C","EOS7X5cuvo6H9Nt8pRnb4BvFPhEkcCeoyZDFAwxa2s5b262gKzEdJ","EOS7X7mPvicrZqg5wTiQPThs1kYnKdSU6s7j7Y4GGMSqC85jDVeLh","EOS7Yx95DyxGaKYLfZzhAmJcqii32VmDDXpdq1sHLz8cC1jL6VDbC","EOS7aXQ9wgDBSMbjKHTFLTBmosvXURKVpvRFCxpPRju1XT7qcmPtD","EOS7d4LKqFdXntyHhkP4W88fcZSZ2UvU2QrfyYzogNnuM8iidDEe8","EOS7dZJoNbi31HJJxr4MRC5aEeNkKPZDk6eCdqfA828osKjit8tER","EOS7eD8NfuqJSeJoL1GDcFGxUr8yek6k4hbRw2Pb5XqKaqRLzgDVm","EOS7eWBf5yjia8W8DkU95PKP3D4tZbeUnDkDRLGAFu79bQ37uadme","EOS7fYf789V4bW4CbnLoFvRYiizsorhMgbdRw2Xspa71KUQjpBr1E","EOS7g9ns9s3rgkTzwpkQVQ3wGpErCrwcV9YbeZiB7uEJXTe8Zaufc","EOS7iEqg3d7BFBpwAqAAhMeMMPJ8kpNmbqk87xw8g4ANnkfqAW4Jw","EOS7msxhT54iC3wNWJU5GhSd9otveDCCh4xVBF2yRF7jiXwqUBANx","EOS7nuEzg5uGnabo45P9622HkSWmL26DFsGa79Ys6BWd5evw53dCY","EOS7pAQxsMS5Qp4ZTAGoFQ8Hog2CmXJYxzZwq4aE4dVsTHdGtiUNR","EOS7pERV73nhme3WYgSJ9WUDmrfWT8ipK4foN2ohzZj1Mo6SYf96a","EOS7pXGG4Tfi14Tsz9WgUVUH2Se7mSqgk8won7Fbu3EAjwhqQPhUD","EOS7pue2Fb9vzFouiZDJqH3u8pMasEKjLhvM66g2gmRzynDW6SQYr","EOS7q23WWf3DgETCcafuR1sVKxWBBZTQDP43enXPM9nBVjeQZAiAz","EOS7qvXyxnvtZSqb3zY358E9fdXZSRpts9FJeFNMwyXjQpBjBYUGK","EOS7siknPJXFLBr9aMaW52e5uVLyVKk1ws3XYxJ9GFoFQqVVqZ1DT","EOS7tPgHPFiLDDo8cdiBtbRp9gyjrz4thfPJASAjNqZwJGpczJz7Y","EOS7tXjnaRhJZRyoqJw8zSEugtf7Bd2LusW6qNtbWwC4KarwZeQPj","EOS7tiKPwG5hi82j1iHXYixV7U3E2oVXbxdnuD1k6MhCEGe5Lfyoh","EOS7tniLVBhFrNXa618BUozVztLe215Ft9yJPCgZJtUGoxSDETWwM","EOS7uCcguEwhXizR1MStL2zqtisccAP3QNuoNjmEypT7yFpvLx6Ng","EOS7uXtmWoZ9vTkUVYZa6QyGjTFfRVJNz2EYiGbkqRuHUKHY9tsTo","EOS7unXSN1HYpoTUvEqbDJ8JWU6fN3HYXCjf5kr5qUPisVnHGXUUQ","EOS7vJiQvj1ooZweWRXkhshvn6dDQNfXPx66zDh1DvzvsX9jc6vTi","EOS7vhr5RcuM7C5QPnkatgskMccejEkb5AQRU9htooPM2tQ3EHxaR","EOS7voxniegj5HL3UnnWanB21wZJnUmj2xUUFT4UPToQo7NGPHmHW","EOS7wS9JDedhWLfGaF9uftQvcXq7shaCXoXNDoUsYsJ9gBicbbpYu","EOS7xggspakp6YqNKG3u97y33VFV5LHqWvoUdo6MWSQsrwS2dRb5V","EOS7xpSi2HvMBZdy46xE94oo8ut3j76tfV1gB11GtthDnJLQ4aYxh","EOS7xvHSD9QvFa6BzvZ3qhrfe6nvVhyEW45aXzxBcDg55B6FrLmA6","EOS7y5Cgt1eZsHZn3CKRHbEAbDG4D2AMKpDy3F9WRb59tAkDerJSf","EOS7yuLTVxUr9saFVMEj8tmpZsZY6887jwcDgiAsgk2JZFYesxRaN","EOS7z99Gk6Uut9W7Jt2Vb2pqkNpuWCqF9xwpHfp3ZwF77EQoQkZvb","EOS7zEMhsgqKr9Q4oqES6cMj24wgSszDDTGejvznvmxBzoVWUXeAr","EOS7zuUYLg7bg3waZYjyv42yRvWETBesKpqZnAdE4XZJ7CP6BgzyX","EOS81UnzKPjrC8VD5k1kMFeQax1LE8SkyMbVwBeSHyLRoishJAZJH","EOS82Dhmhsf6ioyBF9SbEP1iPmKm2nzf1CCC2MtcT46JURsjhvAm7","EOS82jNaQCnXKwnQueVALKPhwqrQ3Xu9Q4tqQn6bUy8NjiPKmzGK7","EOS83vsaoj7ShCuj9yYmdsZPJrf64D2gkNMKL6eQUooWSEDkjzWYM","EOS846UhGWBXfZopA7BzsVZUKoVLCdLQkUcJ9VPUweDmLUvKmA3Xm","EOS84ieG9TA8yRvqoxyUhMAbzv7p6NsGk4uGt78qwxWqanRVrxzFA","EOS84oWuad2eb3DhnC5vFWNpeNm2Tr3mtvP3WRpFJMoWbf1Efr1kz","EOS858oSTqxwa8CtZtAxYz6ZiPHN1SxA9WCCk4UXqtEVMGTkPjjKk","EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM","EOS85ToYftzz2jHJpbzAN413CFzLpdu3iY6UehxN9RjTi3jmFmeuX","EOS86qgjrZD66LZqcn7ZWx7jDbaGDSEdWnyxhs4TrQkEYaLvUthtG","EOS88i53KHV5qMCzqkm5G9DRzetR45DzQCLuzb3nFWksEY8MATpj1","EOS89REc7WcGjQen317b8dTuG6bcnhYpcdU8GXSfyQVEcXxhT56qH","EOS89cH461kPszbdxuzkoH4BuwHWQqNGiTsNC92roEh3KyWiNNrY9","EOS8An1ZJ5XD5R94XS5VrA6Ac2a3HgzYCY2G6QRywNiNiCVw12s7h","EOS8B3fhbrUqTEnC27YubiiQ7g8tymhe7LAYuAwSY7GBCDKMvWARC","EOS8FYmNAJcTRbYMwtoZaQUcchiw5pSx2tNJUoMPTV7CL1NqjKv6M","EOS8FeXbqGPbtyyrebtnUu2oZw8HeQroeX7mNFVU4USKFFGfusq3f","EOS8GeZhnz6g1aqX3Cu5YpQbzs6sNXyGvRJqruekhaYsvdy8FsqxM","EOS8HP3Nf2vq3oWHPofC7JXbrfNgt7hU1BWK6RE9pfQd4rzRMRoSi","EOS8K6ySs4TL23hUJ97pUN7NNBGw8rUVEJCvSDDgh4dvZMXPQQK4p","EOS8LhjUcUPQ9X7eZ2LF8HzoAFvuFB1MnHB2aD76HuBg22xvXcwhy","EOS8LsgToQhqAXXq6kUsGGntGtipc65ZxUr7PUXMNSuFKN78nwBV7","EOS8NijGLHT8WyDmt2nqMwfP1hr8EiYx5JCYBWSP9S26WgbeugvSJ","EOS8NvSq3Y5nDbhv6oyQp3aARV8gNJEvcUZma9iVyDfFte5ccGGVb","EOS8P98FQPuh8vD37JWUCEX6JPDzajLd2cS1jCEAS2XZz6NMtSQYQ","EOS8TaFqALkWKUEjstZRzth4TGgrtuvKM3Be6Y7NHRCZWS3Cc7uSd","EOS8TrFqTjVvd1EUM21EVSN31iAL76ZVWiLGVWCm5Yd9Yj6PuVLUo","EOS8UDLmnUNgXSyURG1r6FXGuzom77atmUyiGi55fJ6ur2tDpqrox","EOS8UJXrg8k76tFVYcBr1uYbxzfshGsEph4ugB8Dr6AsUKRiQSbfA","EOS8WfmQMYRiTNLmN2c3BS6XizGAzHW55DjDtfSfzDpfVjfN3N45K","EOS8XWc8LZTZWamDAXbxv9RMTESLb27RYiHf1hccWo6geWqaKhF9G","EOS8Y2TmxSWWToEikNyE5TJsTozuz2TLs5zsbEoUvjnTZc79R5hr5","EOS8ar9v7Jho8FxrJ4eWS3CRDsukqXmm8oQZyCzyYM5Ca7Q19GG8g","EOS8bRN1gVzx2iQ1FovtiPnKJKZ5U9gPDKxYsNppEysXwWSWZNZTY","EOS8bTjvkJUZvqxm6WgArEcXYPX4aGoU2cNAq4dTnQgkVzYikXTjq","EOS8cQUpNxjP8yz8YcMZ9C6Tm7BzMEUr14AywNfFB9cNAQ8UWjZPm","EOS8dsmkC5xbfeL4CAYgZWkrHDtnQ8MAynCaapiaHx9oVcZFJcFrX","EOS8g7HYzdsHnLoWhBaFZdq34e53aiiVQaCN9LPKyopi6eyfhyiQb","EOS8hBENVtaF9fS2TYBz5L14QUAKgyW3dqc8kpDzYK5ExRdELBvNG","EOS8hzQzvDvhyzhdkskPPApGPZpRTf6onSqLZRLYW97hLgBp4ecc1","EOS8iYNCUMgDN1nzWuPaLpWuj1Q4fQKBVh9qzm5NnKCqViQeNzZgS","EOS8kCcwwApzgYJGBKtztJ4a1dPwXcgQXRsYFVBdzizgKPe49XnVu","EOS8mA4GMpJQT9Vko6a9Wut2HyVhwgrTmdNCTKSaPUPZe78cFC4FH","EOS8mADZC74o9MLDm2bdzLGmimDPmA3UeaYjpThBokmrSzdLkzwMu"]



GOWALLET:

POST /v1/wallet/sign_transaction HTTP/1.0
Host: localhost
content-length: 716
Accept: */*
Connection: close

[{"expiration":"2018-04-18T17:57:49","region":0,"ref_block_num":44982,"ref_block_prefix":3032689784,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500002059b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":[],"context_free_data":[]},["EOS5Dg9cu3yn5cMpWkrZnhmYk2xDBWmu62Sj2dNrWn6Ui82eoYJQh"],"0000000000000000000000000000000000000000000000000000000000000000"]

HTTP/1.0 201 Created
Date: Wed, 18 Apr 2018 17:57:21 GMT
Content-Length: 415
Content-Type: text/plain; charset=utf-8

{"expiration":"2018-04-18T17:57:49","region":0,"ref_block_num":44982,"ref_block_prefix":3032689784,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":""}],"signatures":["EOSK4CLky4ytx7eTuZJx885CmcgvJwr5oxBCZFzSHbh4uRkXwDGR1paVBtWNgCXxMyTxS5P7pYRSH3K4u5DLsN4HeZpmb7H1F"],"context_free_data":[]}


KEOSD:


POST /v1/wallet/sign_transaction HTTP/1.0
Host: localhost
content-length: 716
Accept: */*
Connection: close

[{"expiration":"2018-04-18T17:57:49","region":0,"ref_block_num":44982,"ref_block_prefix":3032689784,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500002059b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":[],"context_free_data":[]},["EOS5Dg9cu3yn5cMpWkrZnhmYk2xDBWmu62Sj2dNrWn6Ui82eoYJQh"],"0000000000000000000000000000000000000000000000000000000000000000"]

HTTP/1.1 201 Created
Content-Length: 688
Content-type: application/json
Server: WebSocket++/0.7.0

{"expiration":"2018-04-18T17:57:49","region":0,"ref_block_num":44982,"ref_block_prefix":3032689784,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500002059b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":["EOSKZH6wvyt7SW7wYWrHuBZcXRq3KhV5zhPyYmG6g685swGs1s66sqvDgCRXifswuisAUiiA2JTU4aTzq8KbAsr7zBbqXVMeD"],"context_free_data":[]}


------------



POST /v1/wallet/sign_transaction HTTP/1.0
Host: localhost
content-length: 716
Accept: */*
Connection: close

[{"expiration":"2018-04-18T20:59:43","region":0,"ref_block_num":55374,"ref_block_prefix":4116381623,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500000959b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":[],"context_free_data":[]},["EOS5Dg9cu3yn5cMpWkrZnhmYk2xDBWmu62Sj2dNrWn6Ui82eoYJQh"],"0000000000000000000000000000000000000000000000000000000000000000"]

HTTP/1.0 201 Created
Date: Wed, 18 Apr 2018 20:59:16 GMT
Content-Length: 415
Content-Type: text/plain; charset=utf-8

{"expiration":"2018-04-18T20:59:43","region":0,"ref_block_num":55374,"ref_block_prefix":4116381623,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":""}],"signatures":["EOSKaXs9RjK2eNYPtXsZxW8gwrkVLLE3u5iajGQjmjQqBUVoa5WfBajmnuesf7JhUSRGbKxqSPyAnzr4eGvqN8RintQ36zaPQ"],"context_free_data":[]}




------------

POST /v1/chain/push_transaction HTTP/1.0
Host: localhost
content-length: 285
Accept: */*
Connection: close

{"signatures":["EOSKaXs9RjK2eNYPtXsZxW8gwrkVLLE3u5iajGQjmjQqBUVoa5WfBajmnuesf7JhUSRGbKxqSPyAnzr4eGvqN8RintQ36zaPQ"],"compression":"none","packed_context_free_data":"","packed_trx":"bfb1d75a00004ed8b7ff5af500000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed323200"}


HTTP/1.1 400 Bad Request
Content-Length: 989
Content-type: application/json
Server: WebSocket++/0.7.0

{"code":400,"message":"Bad Request","error":{"code":3030001,"name":"tx_missing_auth","what":"missing required authority","details":[{"message":"missing authority of ","file":"apply_context.cpp","line_number":146,"method":"require_authorization"},{"message":"","file":"eosio_contract.cpp","line_number":104,"method":"apply_eosio_newaccount"},{"message":"","file":"apply_context.cpp","line_number":30,"method":"exec_one"},{"message":"","file":"chain_controller.cpp","line_number":1989,"method":"__apply_transaction"},{"message":"","file":"chain_controller.cpp","line_number":2024,"method":"_apply_transaction"},{"message":"","file":"chain_controller.cpp","line_number":2225,"method":"wrap_transaction_processing"},{"message":"","file":"chain_controller.cpp","line_number":353,"method":"_push_transaction"},{"message":"","file":"chain_controller.cpp","line_number":328,"method":"_push_transaction"},{"message":"","file":"chain_controller.cpp","line_number":284,"method":"push_transaction"}]}}


---------------- WITH THE EOSIO WALLET

POST /v1/wallet/sign_transaction HTTP/1.0
Host: localhost
content-length: 716
Accept: */*
Connection: close

[{"expiration":"2018-04-18T20:59:46","region":0,"ref_block_num":35871,"ref_block_prefix":2098803478,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500000959b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":[],"context_free_data":[]},["EOS71W8hvF43Eq6GQBRhuc5mvWKtknxzmb9NzNwPGpcEm2xAZaG8c"],"0000000000000000000000000000000000000000000000000000000000000000"]HTTP/1.1 201 Created
Content-Length: 688
Content-type: application/json
Server: WebSocket++/0.7.0

{"expiration":"2018-04-18T20:59:46","region":0,"ref_block_num":35871,"ref_block_prefix":2098803478,"max_net_usage_words":0,"max_kcpu_usage":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500000959b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"}],"signatures":["EOSKhJKBad7XFCF1eYnBihGaCs2LepWmUwhvTZPzkezYZMU6v2wp7PcLHY9v6MRLx6tSQkBUvZP6cZpdQ3WS8L6YaM1uQy5w5"],"context_free_data":[]}


POST /v1/chain/push_transaction HTTP/1.0
Host: cbillett.eoscanada.com
content-length: 533
Accept: */*
Connection: close

{"signatures":["EOSKhJKBad7XFCF1eYnBihGaCs2LepWmUwhvTZPzkezYZMU6v2wp7PcLHY9v6MRLx6tSQkBUvZP6cZpdQ3WS8L6YaM1uQy5w5"],"compression":"none","packed_context_free_data":"","packed_trx":"c2b1d75a00001f8c1633197d00000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000959b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"}


HTTP/1.0 202 Accepted
Content-Length: 1517
Content-Type: application/json
Server: WebSocket++/0.7.0
Date: Wed, 18 Apr 2018 20:59:16 GMT

{"transaction_id":"ebd1b9d89e84432e2079e26f0decc980436554b0ff7802d17ae424004afd1608","processed":{"status":"executed","kcpu_usage":102,"net_usage_words":44,"id":"ebd1b9d89e84432e2079e26f0decc980436554b0ff7802d17ae424004afd1608","action_traces":[{"receiver":"eosio","context_free":false,"cpu_usage":2939,"act":{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":{"creator":"eosio","name":"abourget14","owner":{"threshold":1,"keys":[{"key":"EOS71W8hvF43Eq6GQBRhuc5mvWKtknxzmb9NzNwPGpcEm2xAZaG8c","weight":1}],"accounts":[]},"active":{"threshold":1,"keys":[{"key":"EOS71W8hvF43Eq6GQBRhuc5mvWKtknxzmb9NzNwPGpcEm2xAZaG8c","weight":1}],"accounts":[]},"recovery":{"threshold":1,"keys":[],"accounts":[{"permission":{"actor":"eosio","permission":"active"},"weight":1}]}},"hex_data":"0000000000ea305500000959b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"},"console":"","data_access":[{"type":"write","code":"eosio","scope":"eosio.auth","sequence":6}],"_profiling_us":215}],"deferred_transaction_requests":[],"read_locks":[],"write_locks":[{"account":"eosio","scope":"eosio.auth"}],"cpu_usage":104448,"net_usage":352,"packed_trx_digest":"3e398b8278f01a6e31e8ad8f59e1f27dd1f9c0f6b8a08d1fecab9c125418517c","region_id":0,"cycle_index":1,"shard_index":0,"_profiling_us":268,"_setup_profiling_us":174}}






// With my wallet:
{"signatures":["EOSKaXs9RjK2eNYPtXsZxW8gwrkVLLE3u5iajGQjmjQqBUVoa5WfBajmnuesf7JhUSRGbKxqSPyAnzr4eGvqN8RintQ36zaPQ"],"compression":"none","packed_context_free_data":"","packed_trx":"bfb1d75a00004ed8b7ff5af500000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed323200"}
// With the other wallet:
{"signatures":["EOSKhJKBad7XFCF1eYnBihGaCs2LepWmUwhvTZPzkezYZMU6v2wp7PcLHY9v6MRLx6tSQkBUvZP6cZpdQ3WS8L6YaM1uQy5w5"],"compression":"none","packed_context_free_data":"","packed_trx":"c2b1d75a00001f8c1633197d00000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000959b1abe9310100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100"}






// FROm OUR wallet, packed by cleos:
{"signatures":["EOSK4JJ9wWuREGesAZ17NpwbaMoN5RUxJgTaavrmAjHA4jDjJQuf9vf33Jyn2wz2h11Vkc15V2ySYtA7KSQDACubaLcVAPotC"],"compression":"none","packed_context_free_data":"","packed_trx":""}

// From our sigs, packed by ME:
{"signatures":["EOSJv66rZ29iWPzzG2g3bSqFbeZRb1Qb2GHaPQ9JC3rQC8yQLKbawyDa99fcBX5zLvUd1wFwbDGCjDJw9Gc5VJaiFp6JFB6Gx"],"compression":"none","packed_context_free_data":"","packed_trx":""


e2bbd75a 0000 f0e1e13f61ef00000000010000000000ea3055 00409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500c00a59b1abe93101 00000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa01 00000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100

e5bbd75a 0000 f6e1c3f4903200000000010000000000ea3055 00409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000b59b1abe93101 00000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa01 00000100000001000317512c6c36d468953e04b3c638087709576c55828412dba1a5f0aee5065737fa0100000100000000010000000000ea305500000000a8ed32320100





SIGNING THIS payload: ad9bd85a 0000 14c5 1567affe 000000 00 01 0000000000ea3055 00409e9a2264b89a
01 0000000000ea3055 00000000a8ed3232 7c0000000000ea3055 00000041b1abe931 01000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100

SIGNING THIS payload: ae9bd85a 0000 15c5 fa2be2b7 000000 00 01 0000000000ea3055 00409e9a2264b89a
01 0000000000ea3055 00000000a8ed3232 7c0000000000ea3055 00000042b1abe931 01000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100


8e9cd85a0000d6c65b1306df00000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea30550000000000ade93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100
8e9cd85a0000d6c65b1306df00000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea30550000000000ade93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100









FAILED TX via keosd:

{"signatures":["EOSKffxe9NsHW1Q9cYMiQVjNtEB9bWZwcb5stSZHgYo9Dge3Gt11jYZ5t1ptuefTpmH35R4KMCwswDhZtA17U3NYYkzxQJN6t"],"compression":"none","packed_context_free_data":"","packed_trx":"

60a3d85a 0000 7ad4 89fd2237 00 00 00 00 02 (2 actions)
// eosio:setcode, perm-level: 1 eosio:active
0000000000ea3055 00000040258ab2c2 01 0000000000ea3055 00000000a8ed3232

df2a
0000000000ea3055 00 00
d32a
0061736d0100000001581060037f7e7f0060057f7e7e7e7e0060047f7e7e7e0060027f7f0060027f7e0060000060027e7e0060017f006000017f60027f7f017f60037f7f7f017f60017e0060047e7e7e7e0060027e7f0060037e7e7e0060017f017f02ea010c03656e760561626f7274000503656e7610616374696f6e5f646174615f73697a65000803656e760c656f73696f5f617373657274000303656e760a656f73696f5f65786974000703656e76066d656d637079000a03656e76066d656d736574000a03656e7610726561645f616374696f6e5f64617461000903656e760c726571756972655f61757468000b03656e760d726571756972655f6175746832000603656e76147365745f6163746976655f70726f647563657273000303656e760e7365745f70726976696c65676564000d03656e76137365745f7265736f757263655f6c696d697473000c031b1a09070e00090109020903090409090303030f07070a0f090f0705040501700106060503010001077707066d656d6f72790200165f5a6571524b3131636865636b73756d32353653315f000c305f5a4e35656f73696f3132726571756972655f6175746845524b4e535f31367065726d697373696f6e5f6c6576656c45000d056170706c79000e066d656d636d700020066d616c6c6f63002104667265650024090c010041000b06250f111315170ad5251a0b002000200141201020450b0e002000290300200029030810080be60201017f4100410028020441e0006b22033602040240024020012000520d00200320013703580240024002400240200242ffffb7f6a497b2d942570d002002428080b8f6a497b2d942510d01200242808080c093fad6d942510d0220024280808080b6f7d6d942520d05200341003602542003410136025020032003290350370208200341d8006a200341086a10101a0c040b20024280808080daac9bd6ba7f510d022002428080b8f6a4979ad942520d042003410036024c2003410236024820032003290348370210200341d8006a200341106a10121a0c030b200341003602442003410336024020032003290340370218200341d8006a200341186a10141a0c020b2003410036023c2003410436023820032003290338370220200341d8006a200341206a10161a0c010b200341003602342003410536023020032003290330370228200341d8006a200341286a10181a0b4100200341e0006a3602040f0b41001003000b12002000290300100720012002410047100a0bee0103017f017e047f410028020441106b220521074100200536020420012802042102200128020021060240024010012201418104490d002001102121050c010b410020052001410f6a4170716b22053602040b2005200110061a200741003a000820074200370300200141074b4110100220072005410810041a200141084741101002200741086a2204200541086a410110041a02402001418104490d00200510240b200020024101756a210120042d000021052007290300210302402002410171450d00200128020020066a28020021060b20012003200541ff017120061100004100200741106a36020441010b1300200029030010072001200220032004100b0b9e0203027f047e037f410028020441c0006b220a21094100200a36020420012802042102200128020021080240024010012203418104490d002003102121010c010b4100200a2003410f6a4170716b22013602040b2001200310061a2009420037030820094200370300200942003703102009420037031820092001360224200920013602202009200120036a3602282009200941206a36023020092009360238200941386a200941306a101c02402003418104490d00200110240b200020024101756a2101200941186a2903002107200941106a2903002106200941086a29030021052009290300210402402002410171450d00200128020020086a28020021080b2001200420052006200720081101004100200941c0006a36020441010b0900200029030010070ba20203017f037e057f410028020441206b2208210a4100200836020420012802042102200128020021090240024010012201418104490d002001102121080c010b410020082001410f6a4170716b22083602040b2008200110061a200a4200370310200a4200370308200a4200370318200141074b41101002200a41086a2008410810041a2001417871220641084741101002200a41086a41086a2207200841086a410810041a200641104741101002200a41086a41106a2206200841106a410810041a02402001418104490d00200810240b200020024101756a21012006290300210520072903002104200a290308210302402002410171450d00200128020020096a28020021090b200120032004200520091102004100200a41206a36020441010b3a01027f41002802042103200029030010074100200310012200410f6a4170716b22023602042002200010061a200220001009410020033602040b950501067f410028020441306b220721064100200736020420012802042102200128020021050240024010012203418104490d002003102121010c010b410020072003410f6a4170716b22013602040b2001200310061a200641086a22074200370300200642003703002006200120036a36022820062001360220200341034b4110100220062001410410041a2006200141046a360224200641206a200641047210191a02402003418104490d00200110240b410021012006411c6a41003602002006420037021420062006280200360210200728020020062802046b220341306d21070240024002402003450d00200741d6aad52a4f0d01200641106a41086a22042003101d22013602002006411c6a2001200741306c6a36020020062001360214200641086a280200200628020422076b22034101480d0020012007200310041a20042004280200200341306e41306c6a22013602000b200020024101756a210302402002410171450d00200328020020056a28020021050b2006412c6a41003602002006200628021036022020064200370224200120062802146b220141306d210202402001450d00200241d6aad52a4f0d02200641206a41086a22072001101d22013602002006412c6a2001200241306c6a36020020062001360224200641106a41086a280200200628021422006b22024101480d0020012000200210041a20072001200241306e41306c6a3602000b2003200641206a2005110300024020062802242201450d00200641286a20013602002001101e0b024020062802142201450d00200641186a20013602002001101e0b024020062802042201450d00200641086a20013602002001101e0b4100200641306a36020441010f0b200641106a410472101f000b200641206a410472101f000b0600200110070bc30103017f017e037f410028020441106b220421064100200436020420012802042102200128020021050240024010012201418104490d002001102121040c010b410020042001410f6a4170716b22043602040b2004200110061a20064200370308200141074b41101002200641086a2004410810041a2006290308210302402001418104490d00200410240b200020024101756a210102402002410171450d00200128020020056a28020021050b2001200320051104004100200641106a36020441010ba10203047f017e027f4100410028020441106b2208360204200028020421054100210742002106200041086a2102200041046a21030340200520022802004941201002200328020022052d000021042003200541016a2205360200200441ff0071200741ff0171220774ad2006842106200741076a210720044107760d000b0240024002402006a7220720012802042205200128020022046b41306d22034d0d002001200720036b101a20012802002204200141046a2802002205470d010c020b0240200720034f0d00200141046a2004200741306c6a22053602000b20042005460d010b034020082000360200200820043602082008200441086a36020c200841086a2008101b200441306a22042005470d000b0b4100200841106a36020420000bd40201057f0240024002400240024020002802082202200028020422066b41306d20014f0d002006200028020022056b41306d220320016a220441d6aad52a4f0d0241d5aad52a21060240200220056b41306d220241a9d5aa154b0d0020042002410174220620062004491b2206450d020b200641306c101d21020c030b200041046a21000340200641004130100541003602082000200028020041306a22063602002001417f6a22010d000c040b0b41002106410021020c010b2000101f000b2002200641306c6a21042002200341306c6a220521060340200641004130100522064100360208200641306a21062001417f6a22010d000b2005200041046a2203280200200028020022016b220241506d41306c6a2105024020024101480d0020052001200210041a200028020021010b2000200536020020032006360200200041086a20043602002001450d002001101e0f0b0bfe0103037f017e017f200028020021032001280200220628020820062802046b41074b4110100220032006280204410810041a2006200628020441086a360204200028020421022001280200220328020421064100210142002105200341046a210403402006200341086a2802004941201002200428020022062d000021002004200641016a2206360200200041ff0071200141ff0171220174ad2005842105200141076a210120004107760d000b200220053e0200200341086a2104410421000340200428020020064741101002200220006a200341046a2201280200410110041a2001200128020041016a2206360200200041016a22004125470d000b0bdb0101027f200028020021032001280200220228020820022802046b41074b4110100220032002280204410810041a2002200228020441086a360204200028020021002001280200220228020820022802046b41074b41101002200041086a2002280204410810041a2002200228020441086a3602042001280200220228020820022802046b41074b41101002200041106a2002280204410810041a2002200228020441086a3602042001280200220128020820012802046b41074b41101002200041186a2001280204410810041a2001200128020441086a3602040b3701027f02402000410120001b2201102122000d0003404100210041002802242202450d012002110500200110212200450d000b0b20000b0e0002402000450d00200010240b0b05001000000b4901037f4100210502402002450d000240034020002d0000220320012d00002204470d01200141016a2101200041016a21002002417f6a22020d000c020b0b200320046b21050b20050b08004128200010220bcd04010c7f02402001450d00024020002802c041220d0d004110210d200041c0c1006a41103602000b200141086a200141046a41077122026b200120021b210202400240024020002802c441220a200d4f0d002000200a410c6c6a4180c0006a21010240200a0d0020004184c0006a220d2802000d0020014180c000360200200d20003602000b200241046a210a034002402001280208220d200a6a20012802004b0d002001280204200d6a220d200d28020041808080807871200272360200200141086a22012001280200200a6a360200200d200d28020041808080807872360200200d41046a22010d030b2000102322010d000b0b41fcffffff0720026b2104200041c8c1006a210b200041c0c1006a210c20002802c8412203210d03402000200d410c6c6a22014188c0006a28020020014180c0006a2205280200464180c200100220014184c0006a280200220641046a210d0340200620052802006a2107200d417c6a2208280200220941ffffffff07712101024020094100480d000240200120024f0d000340200d20016a220a20074f0d01200a280200220a4100480d012001200a41ffffffff07716a41046a22012002490d000b0b20082001200220012002491b200941808080807871723602000240200120024d0d00200d20026a200420016a41ffffffff07713602000b200120024f0d040b200d20016a41046a220d2007490d000b41002101200b4100200b28020041016a220d200d200c280200461b220d360200200d2003470d000b0b20010f0b2008200828020041808080807872360200200d0f0b41000b870501087f20002802c44121010240024041002d00d642450d0041002802d84221070c010b3f002107410041013a00d6424100200741107422073602d8420b200721030240024002400240200741ffff036a41107622023f0022084d0d00200220086b40001a4100210820023f00470d0141002802d84221030b41002108410020033602d84220074100480d0020002001410c6c6a210220074180800441808008200741ffff037122084181f8034922061b6a2008200741ffff077120061b6b20076b2107024041002d00d6420d003f002103410041013a00d6424100200341107422033602d8420b20024180c0006a210220074100480d01200321060240200741076a417871220520036a41ffff036a41107622083f0022044d0d00200820046b40001a20083f00470d0241002802d84221060b4100200620056a3602d8422003417f460d0120002001410c6c6a22014184c0006a2802002206200228020022086a2003460d020240200820014188c0006a22052802002201460d00200620016a2206200628020041808080807871417c20016b20086a72360200200520022802003602002006200628020041ffffffff07713602000b200041c4c1006a2202200228020041016a220236020020002002410c6c6a22004184c0006a200336020020004180c0006a220820073602000b20080f0b02402002280200220820002001410c6c6a22034188c0006a22012802002207460d0020034184c0006a28020020076a2203200328020041808080807871417c20076b20086a72360200200120022802003602002003200328020041ffffffff07713602000b2000200041c4c1006a220728020041016a22033602c0412007200336020041000f0b2002200820076a36020020020b7b01037f024002402000450d0041002802e84122024101480d0041a8c00021032002410c6c41a8c0006a21010340200341046a2802002202450d010240200241046a20004b0d00200220032802006a20004b0d030b2003410c6a22032001490d000b0b0f0b2000417c6a2203200328020041ffffffff07713602000b0300000b0b7a040041040b04603100000041100b0572656164000041200b0467657400004180c2000b566d616c6c6f635f66726f6d5f6672656564207761732064657369676e656420746f206f6e6c792062652063616c6c6564206166746572205f686561702077617320636f6d706c6574656c7920616c6c6f636174656400

0000000000ea3055 00000000b863b2c2 01 0000000000ea3055 00000000a8ed3232
fc03
0000000000ea3055
00
07
12 7365745f6163636f756e745f6c696d697473
00
04
07 6163636f756e74
0c 6163636f756e745f6e616d650972616d5f627974657305696e7436340a6e65745f77656967687405696e7436340a6370755f77656967687405696e74363407736574707269760002076163636f756e740c6163636f756e745f6e616d650769735f7072697604696e7438117365745f676c6f62616c5f6c696d6974730001136370755f757365635f7065725f706572696f6405696e7436340c70726f64756365725f6b657900020d70726f64756365725f6e616d650c6163636f756e745f6e616d6511626c6f636b5f7369676e696e675f6b65790a7075626c69635f6b65790d7365745f70726f64756365727300020776657273696f6e0675696e7433320970726f6475636572730e70726f64756365725f6b65795b5d0c726571756972655f6175746800010466726f6d0c6163636f756e745f6e616d65056e6f6e636500010576616c756506737472696e67060000ce4eba68b2c2127365745f6163636f756e745f6c696d697473000000ce4ebac8b2c2117365745f676c6f62616c5f6c696d6974730000000060bb5bb3c207736574707269760000000038d15bb3c20d7365745f70726f64756365727300000000a0656dacba0c726571756972655f6175746800000000000085269d056e6f6e6365000000"}


Accepted via local KeyBag signing:

{"signatures":["EOSKhQ6pCAQXrAEjqjkEXRH6jna7Gk2dUK7XCgdG4eHdWF3JiH2ceBBCnw1JMrk8UMYM4Q8fcwAXRfuZnzL6DJa8i4XyDH1w5"],"compression":"none","packed_context_free_data":"","packed_trx":"

66a3d85a 0000 86d4 94c82904 00 00 00 00 02 (2 actions)
// eosio:setcode, perm-level: 1 eosio:active
0000000000ea3055 00000040258ab2c2 01 0000000000ea3055 00000000a8ed3232

df2a len=5471
0000000000ea3055 00 00
d32a len=5459
0061736d0100000001581060037f7e7f0060057f7e7e7e7e0060047f7e7e7e0060027f7f0060027f7e0060000060027e7e0060017f006000017f60027f7f017f60037f7f7f017f60017e0060047e7e7e7e0060027e7f0060037e7e7e0060017f017f02ea010c03656e760561626f7274000503656e7610616374696f6e5f646174615f73697a65000803656e760c656f73696f5f617373657274000303656e760a656f73696f5f65786974000703656e76066d656d637079000a03656e76066d656d736574000a03656e7610726561645f616374696f6e5f64617461000903656e760c726571756972655f61757468000b03656e760d726571756972655f6175746832000603656e76147365745f6163746976655f70726f647563657273000303656e760e7365745f70726976696c65676564000d03656e76137365745f7265736f757263655f6c696d697473000c031b1a09070e00090109020903090409090303030f07070a0f090f0705040501700106060503010001077707066d656d6f72790200165f5a6571524b3131636865636b73756d32353653315f000c305f5a4e35656f73696f3132726571756972655f6175746845524b4e535f31367065726d697373696f6e5f6c6576656c45000d056170706c79000e066d656d636d700020066d616c6c6f63002104667265650024090c010041000b06250f111315170ad5251a0b002000200141201020450b0e002000290300200029030810080be60201017f4100410028020441e0006b22033602040240024020012000520d00200320013703580240024002400240200242ffffb7f6a497b2d942570d002002428080b8f6a497b2d942510d01200242808080c093fad6d942510d0220024280808080b6f7d6d942520d05200341003602542003410136025020032003290350370208200341d8006a200341086a10101a0c040b20024280808080daac9bd6ba7f510d022002428080b8f6a4979ad942520d042003410036024c2003410236024820032003290348370210200341d8006a200341106a10121a0c030b200341003602442003410336024020032003290340370218200341d8006a200341186a10141a0c020b2003410036023c2003410436023820032003290338370220200341d8006a200341206a10161a0c010b200341003602342003410536023020032003290330370228200341d8006a200341286a10181a0b4100200341e0006a3602040f0b41001003000b12002000290300100720012002410047100a0bee0103017f017e047f410028020441106b220521074100200536020420012802042102200128020021060240024010012201418104490d002001102121050c010b410020052001410f6a4170716b22053602040b2005200110061a200741003a000820074200370300200141074b4110100220072005410810041a200141084741101002200741086a2204200541086a410110041a02402001418104490d00200510240b200020024101756a210120042d000021052007290300210302402002410171450d00200128020020066a28020021060b20012003200541ff017120061100004100200741106a36020441010b1300200029030010072001200220032004100b0b9e0203027f047e037f410028020441c0006b220a21094100200a36020420012802042102200128020021080240024010012203418104490d002003102121010c010b4100200a2003410f6a4170716b22013602040b2001200310061a2009420037030820094200370300200942003703102009420037031820092001360224200920013602202009200120036a3602282009200941206a36023020092009360238200941386a200941306a101c02402003418104490d00200110240b200020024101756a2101200941186a2903002107200941106a2903002106200941086a29030021052009290300210402402002410171450d00200128020020086a28020021080b2001200420052006200720081101004100200941c0006a36020441010b0900200029030010070ba20203017f037e057f410028020441206b2208210a4100200836020420012802042102200128020021090240024010012201418104490d002001102121080c010b410020082001410f6a4170716b22083602040b2008200110061a200a4200370310200a4200370308200a4200370318200141074b41101002200a41086a2008410810041a2001417871220641084741101002200a41086a41086a2207200841086a410810041a200641104741101002200a41086a41106a2206200841106a410810041a02402001418104490d00200810240b200020024101756a21012006290300210520072903002104200a290308210302402002410171450d00200128020020096a28020021090b200120032004200520091102004100200a41206a36020441010b3a01027f41002802042103200029030010074100200310012200410f6a4170716b22023602042002200010061a200220001009410020033602040b950501067f410028020441306b220721064100200736020420012802042102200128020021050240024010012203418104490d002003102121010c010b410020072003410f6a4170716b22013602040b2001200310061a200641086a22074200370300200642003703002006200120036a36022820062001360220200341034b4110100220062001410410041a2006200141046a360224200641206a200641047210191a02402003418104490d00200110240b410021012006411c6a41003602002006420037021420062006280200360210200728020020062802046b220341306d21070240024002402003450d00200741d6aad52a4f0d01200641106a41086a22042003101d22013602002006411c6a2001200741306c6a36020020062001360214200641086a280200200628020422076b22034101480d0020012007200310041a20042004280200200341306e41306c6a22013602000b200020024101756a210302402002410171450d00200328020020056a28020021050b2006412c6a41003602002006200628021036022020064200370224200120062802146b220141306d210202402001450d00200241d6aad52a4f0d02200641206a41086a22072001101d22013602002006412c6a2001200241306c6a36020020062001360224200641106a41086a280200200628021422006b22024101480d0020012000200210041a20072001200241306e41306c6a3602000b2003200641206a2005110300024020062802242201450d00200641286a20013602002001101e0b024020062802142201450d00200641186a20013602002001101e0b024020062802042201450d00200641086a20013602002001101e0b4100200641306a36020441010f0b200641106a410472101f000b200641206a410472101f000b0600200110070bc30103017f017e037f410028020441106b220421064100200436020420012802042102200128020021050240024010012201418104490d002001102121040c010b410020042001410f6a4170716b22043602040b2004200110061a20064200370308200141074b41101002200641086a2004410810041a2006290308210302402001418104490d00200410240b200020024101756a210102402002410171450d00200128020020056a28020021050b2001200320051104004100200641106a36020441010ba10203047f017e027f4100410028020441106b2208360204200028020421054100210742002106200041086a2102200041046a21030340200520022802004941201002200328020022052d000021042003200541016a2205360200200441ff0071200741ff0171220774ad2006842106200741076a210720044107760d000b0240024002402006a7220720012802042205200128020022046b41306d22034d0d002001200720036b101a20012802002204200141046a2802002205470d010c020b0240200720034f0d00200141046a2004200741306c6a22053602000b20042005460d010b034020082000360200200820043602082008200441086a36020c200841086a2008101b200441306a22042005470d000b0b4100200841106a36020420000bd40201057f0240024002400240024020002802082202200028020422066b41306d20014f0d002006200028020022056b41306d220320016a220441d6aad52a4f0d0241d5aad52a21060240200220056b41306d220241a9d5aa154b0d0020042002410174220620062004491b2206450d020b200641306c101d21020c030b200041046a21000340200641004130100541003602082000200028020041306a22063602002001417f6a22010d000c040b0b41002106410021020c010b2000101f000b2002200641306c6a21042002200341306c6a220521060340200641004130100522064100360208200641306a21062001417f6a22010d000b2005200041046a2203280200200028020022016b220241506d41306c6a2105024020024101480d0020052001200210041a200028020021010b2000200536020020032006360200200041086a20043602002001450d002001101e0f0b0bfe0103037f017e017f200028020021032001280200220628020820062802046b41074b4110100220032006280204410810041a2006200628020441086a360204200028020421022001280200220328020421064100210142002105200341046a210403402006200341086a2802004941201002200428020022062d000021002004200641016a2206360200200041ff0071200141ff0171220174ad2005842105200141076a210120004107760d000b200220053e0200200341086a2104410421000340200428020020064741101002200220006a200341046a2201280200410110041a2001200128020041016a2206360200200041016a22004125470d000b0bdb0101027f200028020021032001280200220228020820022802046b41074b4110100220032002280204410810041a2002200228020441086a360204200028020021002001280200220228020820022802046b41074b41101002200041086a2002280204410810041a2002200228020441086a3602042001280200220228020820022802046b41074b41101002200041106a2002280204410810041a2002200228020441086a3602042001280200220128020820012802046b41074b41101002200041186a2001280204410810041a2001200128020441086a3602040b3701027f02402000410120001b2201102122000d0003404100210041002802242202450d012002110500200110212200450d000b0b20000b0e0002402000450d00200010240b0b05001000000b4901037f4100210502402002450d000240034020002d0000220320012d00002204470d01200141016a2101200041016a21002002417f6a22020d000c020b0b200320046b21050b20050b08004128200010220bcd04010c7f02402001450d00024020002802c041220d0d004110210d200041c0c1006a41103602000b200141086a200141046a41077122026b200120021b210202400240024020002802c441220a200d4f0d002000200a410c6c6a4180c0006a21010240200a0d0020004184c0006a220d2802000d0020014180c000360200200d20003602000b200241046a210a034002402001280208220d200a6a20012802004b0d002001280204200d6a220d200d28020041808080807871200272360200200141086a22012001280200200a6a360200200d200d28020041808080807872360200200d41046a22010d030b2000102322010d000b0b41fcffffff0720026b2104200041c8c1006a210b200041c0c1006a210c20002802c8412203210d03402000200d410c6c6a22014188c0006a28020020014180c0006a2205280200464180c200100220014184c0006a280200220641046a210d0340200620052802006a2107200d417c6a2208280200220941ffffffff07712101024020094100480d000240200120024f0d000340200d20016a220a20074f0d01200a280200220a4100480d012001200a41ffffffff07716a41046a22012002490d000b0b20082001200220012002491b200941808080807871723602000240200120024d0d00200d20026a200420016a41ffffffff07713602000b200120024f0d040b200d20016a41046a220d2007490d000b41002101200b4100200b28020041016a220d200d200c280200461b220d360200200d2003470d000b0b20010f0b2008200828020041808080807872360200200d0f0b41000b870501087f20002802c44121010240024041002d00d642450d0041002802d84221070c010b3f002107410041013a00d6424100200741107422073602d8420b200721030240024002400240200741ffff036a41107622023f0022084d0d00200220086b40001a4100210820023f00470d0141002802d84221030b41002108410020033602d84220074100480d0020002001410c6c6a210220074180800441808008200741ffff037122084181f8034922061b6a2008200741ffff077120061b6b20076b2107024041002d00d6420d003f002103410041013a00d6424100200341107422033602d8420b20024180c0006a210220074100480d01200321060240200741076a417871220520036a41ffff036a41107622083f0022044d0d00200820046b40001a20083f00470d0241002802d84221060b4100200620056a3602d8422003417f460d0120002001410c6c6a22014184c0006a2802002206200228020022086a2003460d020240200820014188c0006a22052802002201460d00200620016a2206200628020041808080807871417c20016b20086a72360200200520022802003602002006200628020041ffffffff07713602000b200041c4c1006a2202200228020041016a220236020020002002410c6c6a22004184c0006a200336020020004180c0006a220820073602000b20080f0b02402002280200220820002001410c6c6a22034188c0006a22012802002207460d0020034184c0006a28020020076a2203200328020041808080807871417c20076b20086a72360200200120022802003602002003200328020041ffffffff07713602000b2000200041c4c1006a220728020041016a22033602c0412007200336020041000f0b2002200820076a36020020020b7b01037f024002402000450d0041002802e84122024101480d0041a8c00021032002410c6c41a8c0006a21010340200341046a2802002202450d010240200241046a20004b0d00200220032802006a20004b0d030b2003410c6a22032001490d000b0b0f0b2000417c6a2203200328020041ffffffff07713602000b0300000b0b7a040041040b04603100000041100b0572656164000041200b0467657400004180c2000b566d616c6c6f635f66726f6d5f6672656564207761732064657369676e656420746f206f6e6c792062652063616c6c6564206166746572205f686561702077617320636f6d706c6574656c7920616c6c6f636174656400

0000000000ea3055 00000000b863b2c2 01 0000000000ea3055 00000000a8ed3232
fc03 len=508
0000000000ea3055
00 []types
07 []structs
12 7365745f6163636f756e745f6c696d697473 name="set_accounts_limits"
00 base=""
04 []fields
07 6163636f756e74 name="account"
0c6163636f756e745f6e616d650972616d5f627974657305696e7436340a6e65745f77656967687405696e7436340a6370755f77656967687405696e74363407736574707269760002076163636f756e740c6163636f756e745f6e616d650769735f7072697604696e7438117365745f676c6f62616c5f6c696d6974730001136370755f757365635f7065725f706572696f6405696e7436340c70726f64756365725f6b657900020d70726f64756365725f6e616d650c6163636f756e745f6e616d6511626c6f636b5f7369676e696e675f6b65790a7075626c69635f6b65790d7365745f70726f64756365727300020776657273696f6e0675696e7433320970726f6475636572730e70726f64756365725f6b65795b5d0c726571756972655f6175746800010466726f6d0c6163636f756e745f6e616d65056e6f6e636500010576616c756506737472696e67060000ce4eba68b2c2127365745f6163636f756e745f6c696d697473000000ce4ebac8b2c2117365745f676c6f62616c5f6c696d6974730000000060bb5bb3c207736574707269760000000038d15bb3c20d7365745f70726f64756365727300000000a0656dacba0c726571756972655f6175746800000000000085269d056e6f6e6365000000"}


0000000000ea3055
00 07 12 7365745f6163636f756e745f6c696d697473
00
04
07 6163636f756e74
0c6163636f756e745f6e616d650972616d5f627974657305696e7436340a6e65745f77656967687405696e7436340a6370755f77656967687405696e74363407736574707269760002076163636f756e740c6163636f756e745f6e616d650769735f7072697604696e7438117365745f676c6f62616c5f6c696d6974730001136370755f757365635f7065725f706572696f6405696e7436340c70726f64756365725f6b657900020d70726f64756365725f6e616d650c6163636f756e745f6e616d6511626c6f636b5f7369676e696e675f6b65790a7075626c69635f6b65790d7365745f70726f64756365727300020776657273696f6e0675696e7433320970726f6475636572730e70726f64756365725f6b65795b5d0c726571756972655f6175746800010466726f6d0c6163636f756e745f6e616d65056e6f6e636500010576616c756506737472696e67060000ce4eba68b2c2127365745f6163636f756e745f6c696d697473000000ce4ebac8b2c2117365745f676c6f62616c5f6c696d6974730000000060bb5bb3c207736574707269760000000038d15bb3c20d7365745f70726f64756365727300000000a0656dacba0c726571756972655f6175746800000000000085269d056e6f6e6365000000





a5c9d85a 0000 0421 067cf72e
00 00 00 00 01
0000000000ea3055
00000040258ab2c2
01
0000000000ea3055
00000000a8ed3232

df2a
0000000000ea3055
00
00
d32a

0061736d0100000001581060037f7e7f0060057f7e7e7e7e0060047f7e7e7e0060027f7f0060027f7e0060000060027e7e0060017f006000017f60027f7f017f60037f7f7f017f60017e0060047e7e7e7e0060027e7f0060037e7e7e0060017f017f02ea010c03656e760561626f7274000503656e7610616374696f6e5f646174615f73697a65000803656e760c656f73696f5f617373657274000303656e760a656f73696f5f65786974000703656e76066d656d637079000a03656e76066d656d736574000a03656e7610726561645f616374696f6e5f64617461000903656e760c726571756972655f61757468000b03656e760d726571756972655f6175746832000603656e76147365745f6163746976655f70726f647563657273000303656e760e7365745f70726976696c65676564000d03656e76137365745f7265736f757263655f6c696d697473000c031b1a09070e00090109020903090409090303030f07070a0f090f0705040501700106060503010001077707066d656d6f72790200165f5a6571524b3131636865636b73756d32353653315f000c305f5a4e35656f73696f3132726571756972655f6175746845524b4e535f31367065726d697373696f6e5f6c6576656c45000d056170706c79000e066d656d636d700020066d616c6c6f63002104667265650024090c010041000b06250f111315170ad5251a0b002000200141201020450b0e002000290300200029030810080be60201017f4100410028020441e0006b22033602040240024020012000520d00200320013703580240024002400240200242ffffb7f6a497b2d942570d002002428080b8f6a497b2d942510d01200242808080c093fad6d942510d0220024280808080b6f7d6d942520d05200341003602542003410136025020032003290350370208200341d8006a200341086a10101a0c040b20024280808080daac9bd6ba7f510d022002428080b8f6a4979ad942520d042003410036024c2003410236024820032003290348370210200341d8006a200341106a10121a0c030b200341003602442003410336024020032003290340370218200341d8006a200341186a10141a0c020b2003410036023c2003410436023820032003290338370220200341d8006a200341206a10161a0c010b200341003602342003410536023020032003290330370228200341d8006a200341286a10181a0b4100200341e0006a3602040f0b41001003000b12002000290300100720012002410047100a0bee0103017f017e047f410028020441106b220521074100200536020420012802042102200128020021060240024010012201418104490d002001102121050c010b410020052001410f6a4170716b22053602040b2005200110061a200741003a000820074200370300200141074b4110100220072005410810041a200141084741101002200741086a2204200541086a410110041a02402001418104490d00200510240b200020024101756a210120042d000021052007290300210302402002410171450d00200128020020066a28020021060b20012003200541ff017120061100004100200741106a36020441010b1300200029030010072001200220032004100b0b9e0203027f047e037f410028020441c0006b220a21094100200a36020420012802042102200128020021080240024010012203418104490d002003102121010c010b4100200a2003410f6a4170716b22013602040b2001200310061a2009420037030820094200370300200942003703102009420037031820092001360224200920013602202009200120036a3602282009200941206a36023020092009360238200941386a200941306a101c02402003418104490d00200110240b200020024101756a2101200941186a2903002107200941106a2903002106200941086a29030021052009290300210402402002410171450d00200128020020086a28020021080b2001200420052006200720081101004100200941c0006a36020441010b0900200029030010070ba20203017f037e057f410028020441206b2208210a4100200836020420012802042102200128020021090240024010012201418104490d002001102121080c010b410020082001410f6a4170716b22083602040b2008200110061a200a4200370310200a4200370308200a4200370318200141074b41101002200a41086a2008410810041a2001417871220641084741101002200a41086a41086a2207200841086a410810041a200641104741101002200a41086a41106a2206200841106a410810041a02402001418104490d00200810240b200020024101756a21012006290300210520072903002104200a290308210302402002410171450d00200128020020096a28020021090b200120032004200520091102004100200a41206a36020441010b3a01027f41002802042103200029030010074100200310012200410f6a4170716b22023602042002200010061a200220001009410020033602040b950501067f410028020441306b220721064100200736020420012802042102200128020021050240024010012203418104490d002003102121010c010b410020072003410f6a4170716b22013602040b2001200310061a200641086a22074200370300200642003703002006200120036a36022820062001360220200341034b4110100220062001410410041a2006200141046a360224200641206a200641047210191a02402003418104490d00200110240b410021012006411c6a41003602002006420037021420062006280200360210200728020020062802046b220341306d21070240024002402003450d00200741d6aad52a4f0d01200641106a41086a22042003101d22013602002006411c6a2001200741306c6a36020020062001360214200641086a280200200628020422076b22034101480d0020012007200310041a20042004280200200341306e41306c6a22013602000b200020024101756a210302402002410171450d00200328020020056a28020021050b2006412c6a41003602002006200628021036022020064200370224200120062802146b220141306d210202402001450d00200241d6aad52a4f0d02200641206a41086a22072001101d22013602002006412c6a2001200241306c6a36020020062001360224200641106a41086a280200200628021422006b22024101480d0020012000200210041a20072001200241306e41306c6a3602000b2003200641206a2005110300024020062802242201450d00200641286a20013602002001101e0b024020062802142201450d00200641186a20013602002001101e0b024020062802042201450d00200641086a20013602002001101e0b4100200641306a36020441010f0b200641106a410472101f000b200641206a410472101f000b0600200110070bc30103017f017e037f410028020441106b220421064100200436020420012802042102200128020021050240024010012201418104490d002001102121040c010b410020042001410f6a4170716b22043602040b2004200110061a20064200370308200141074b41101002200641086a2004410810041a2006290308210302402001418104490d00200410240b200020024101756a210102402002410171450d00200128020020056a28020021050b2001200320051104004100200641106a36020441010ba10203047f017e027f4100410028020441106b2208360204200028020421054100210742002106200041086a2102200041046a21030340200520022802004941201002200328020022052d000021042003200541016a2205360200200441ff0071200741ff0171220774ad2006842106200741076a210720044107760d000b0240024002402006a7220720012802042205200128020022046b41306d22034d0d002001200720036b101a20012802002204200141046a2802002205470d010c020b0240200720034f0d00200141046a2004200741306c6a22053602000b20042005460d010b034020082000360200200820043602082008200441086a36020c200841086a2008101b200441306a22042005470d000b0b4100200841106a36020420000bd40201057f0240024002400240024020002802082202200028020422066b41306d20014f0d002006200028020022056b41306d220320016a220441d6aad52a4f0d0241d5aad52a21060240200220056b41306d220241a9d5aa154b0d0020042002410174220620062004491b2206450d020b200641306c101d21020c030b200041046a21000340200641004130100541003602082000200028020041306a22063602002001417f6a22010d000c040b0b41002106410021020c010b2000101f000b2002200641306c6a21042002200341306c6a220521060340200641004130100522064100360208200641306a21062001417f6a22010d000b2005200041046a2203280200200028020022016b220241506d41306c6a2105024020024101480d0020052001200210041a200028020021010b2000200536020020032006360200200041086a20043602002001450d002001101e0f0b0bfe0103037f017e017f200028020021032001280200220628020820062802046b41074b4110100220032006280204410810041a2006200628020441086a360204200028020421022001280200220328020421064100210142002105200341046a210403402006200341086a2802004941201002200428020022062d000021002004200641016a2206360200200041ff0071200141ff0171220174ad2005842105200141076a210120004107760d000b200220053e0200200341086a2104410421000340200428020020064741101002200220006a200341046a2201280200410110041a2001200128020041016a2206360200200041016a22004125470d000b0bdb0101027f200028020021032001280200220228020820022802046b41074b4110100220032002280204410810041a2002200228020441086a360204200028020021002001280200220228020820022802046b41074b41101002200041086a2002280204410810041a2002200228020441086a3602042001280200220228020820022802046b41074b41101002200041106a2002280204410810041a2002200228020441086a3602042001280200220128020820012802046b41074b41101002200041186a2001280204410810041a2001200128020441086a3602040b3701027f02402000410120001b2201102122000d0003404100210041002802242202450d012002110500200110212200450d000b0b20000b0e0002402000450d00200010240b0b05001000000b4901037f4100210502402002450d000240034020002d0000220320012d00002204470d01200141016a2101200041016a21002002417f6a22020d000c020b0b200320046b21050b20050b08004128200010220bcd04010c7f02402001450d00024020002802c041220d0d004110210d200041c0c1006a41103602000b200141086a200141046a41077122026b200120021b210202400240024020002802c441220a200d4f0d002000200a410c6c6a4180c0006a21010240200a0d0020004184c0006a220d2802000d0020014180c000360200200d20003602000b200241046a210a034002402001280208220d200a6a20012802004b0d002001280204200d6a220d200d28020041808080807871200272360200200141086a22012001280200200a6a360200200d200d28020041808080807872360200200d41046a22010d030b2000102322010d000b0b41fcffffff0720026b2104200041c8c1006a210b200041c0c1006a210c20002802c8412203210d03402000200d410c6c6a22014188c0006a28020020014180c0006a2205280200464180c200100220014184c0006a280200220641046a210d0340200620052802006a2107200d417c6a2208280200220941ffffffff07712101024020094100480d000240200120024f0d000340200d20016a220a20074f0d01200a280200220a4100480d012001200a41ffffffff07716a41046a22012002490d000b0b20082001200220012002491b200941808080807871723602000240200120024d0d00200d20026a200420016a41ffffffff07713602000b200120024f0d040b200d20016a41046a220d2007490d000b41002101200b4100200b28020041016a220d200d200c280200461b220d360200200d2003470d000b0b20010f0b2008200828020041808080807872360200200d0f0b41000b870501087f20002802c44121010240024041002d00d642450d0041002802d84221070c010b3f002107410041013a00d6424100200741107422073602d8420b200721030240024002400240200741ffff036a41107622023f0022084d0d00200220086b40001a4100210820023f00470d0141002802d84221030b41002108410020033602d84220074100480d0020002001410c6c6a210220074180800441808008200741ffff037122084181f8034922061b6a2008200741ffff077120061b6b20076b2107024041002d00d6420d003f002103410041013a00d6424100200341107422033602d8420b20024180c0006a210220074100480d01200321060240200741076a417871220520036a41ffff036a41107622083f0022044d0d00200820046b40001a20083f00470d0241002802d84221060b4100200620056a3602d8422003417f460d0120002001410c6c6a22014184c0006a2802002206200228020022086a2003460d020240200820014188c0006a22052802002201460d00200620016a2206200628020041808080807871417c20016b20086a72360200200520022802003602002006200628020041ffffffff07713602000b200041c4c1006a2202200228020041016a220236020020002002410c6c6a22004184c0006a200336020020004180c0006a220820073602000b20080f0b02402002280200220820002001410c6c6a22034188c0006a22012802002207460d0020034184c0006a28020020076a2203200328020041808080807871417c20076b20086a72360200200120022802003602002003200328020041ffffffff07713602000b2000200041c4c1006a220728020041016a22033602c0412007200336020041000f0b2002200820076a36020020020b7b01037f024002402000450d0041002802e84122024101480d0041a8c00021032002410c6c41a8c0006a21010340200341046a2802002202450d010240200241046a20004b0d00200220032802006a20004b0d030b2003410c6a22032001490d000b0b0f0b2000417c6a2203200328020041ffffffff07713602000b0300000b0b7a040041040b04603100000041100b0572656164000041200b0467657400004180c2000b566d616c6c6f635f66726f6d5f6672656564207761732064657369676e656420746f206f6e6c792062652063616c6c6564206166746572205f686561702077617320636f6d706c6574656c7920616c6c6f636174656400










{"signatures":["EOSK4mXdcAF56SCFoKW3XSwm1NP92jadrRt9EFbtQQ71B4gpN7kYZhJiWpZFuanAjKJdeTFM9MgWxoWWNhyWufxMTfoXm3wP7"],"compression":"none","packed_context_free_data":"","packed_trx":"eed8d85a0000963f644ca7cc00000000010000000000ea305500000000b863b2c2010000000000ea305500000000a8ed3232fc030000000000ea30550007127365745f6163636f756e745f6c696d6974730004076163636f756e740c6163636f756e745f6e616d650972616d5f627974657305696e7436340a6e65745f77656967687405696e7436340a6370755f77656967687405696e74363407736574707269760002076163636f756e740c6163636f756e745f6e616d650769735f7072697604696e7438117365745f676c6f62616c5f6c696d6974730001136370755f757365635f7065725f706572696f6405696e7436340c70726f64756365725f6b657900020d70726f64756365725f6e616d650c6163636f756e745f6e616d6511626c6f636b5f7369676e696e675f6b65790a7075626c69635f6b65790d7365745f70726f64756365727300020776657273696f6e0675696e7433320970726f6475636572730e70726f64756365725f6b65795b5d0c726571756972655f6175746800010466726f6d0c6163636f756e745f6e616d65056e6f6e636500010576616c756506737472696e67060000ce4eba68b2c2127365745f6163636f756e745f6c696d697473000000ce4ebac8b2c2117365745f676c6f62616c5f6c696d6974730000000060bb5bb3c207736574707269760000000038d15bb3c20d7365745f70726f64756365727300000000a0656dacba0c726571756972655f6175746800000000000085269d056e6f6e6365000000"}


UNAUTHORIZED




0f010000
06
000330907937f721e855bb1ff0b6143af8bd82cba550d1c20438b29dbb30b2a7
4170e244 = 2004 ou 2008, si c'est un nom de secondes depuis l'EPOCH (1970)
d756399dd08666bbbe2e44c0f0967b9dfaa02b33234a2a31e52319961cf2de52
9eb4bc72f5b7b28b86848e837ce30805ae57e81023189eabd77d516a12ac02a7
e9f40f2a2d66d2199a8bc9983e30066f61cbf24230653bb74b9c91e8d7b972ab
0000000000ea3055 eosio
00000000
00 optional not there !!
00205c87c866b524d313866c8c0c036518ee4e898b133d570b41403e31bb072f909d1b8e40775e261e24ae787aeeb127e5f25725ab684f382f7f92a68ee4672e1fac producer_signature 66 bytes
01 []region
0000 region
01 []cycles
01 []shardsummary
00 []readlocks
01 []writelocks
0000000000ea3055 "eosio"
0000000000ea3055 "eosio"
01 []transactions
00 status = executed
13 KCPUUsage = 19
26 NetusageWords = 38
fd25e5ad35924d41b9611304b682dcfca0bf1154008f97268d43bfaba999b104  transaction ID

00010000
07
000002b6fe7206df847da1e7e7691bd4ad25d0998f4fe226a2fbbd901221657d
8988e644
51b3b73face644d29603d2dcebf00af121c4af5e4ed86f5e80844fb4832dbeb0
3394ae8dde1ee9b0df4edbcd4cd7478768e97ca72801699ac07e9614b4470717
2ed846471b8f9d46914be0822d4d43dc2eebe08f759af1b6261af05107d23b88
0000000000ea3055
00000000
00 optional flag
0020783d2ff33f718c6cef10d5d0875ba99177259e184ae7ac92022f09302693b8925ed5429d7ec474cea6dec0fe4f5ab2e3da55188515c6a914e6e97a30d82ee000
01 []Regions
0000
01 []CyclesSummary
01 []Cycles
00 []ReadLocks
00 []WriteLocks
01 []Transactions
00 status
02
26
8149cec70c1ed3fafc6f3ebe90e4c7af423d9b1edd627a4fe1a8f54d9c90f81e
00


3300000003
00000000000000000001000000123403000100033412ead7e5d4ebc2f2a47edbbbbde260f2f13a29668dc0366aa72a8e665d




13000000 len
03 Message Type

02000000bd24????
00 mode
00000200 pending
00 []ids
00 mode
be240000 pending
00 ids[]



-------------

00096e88
0000
00000000
00
00
00

00 cfa
01 actions[] = 1
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
7c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100


00096e88
0000
00000000
00
00
00
00
01
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
7c // ACTION
0000000000ea3055
00000059b1abe931
owner
01000000 threshold
01 keys[]
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf pubkey (curve=00)
0100 weight
00 accounts[]

active
01000000
01
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf pubkey
0100 weight
00 accounts[]
0100000000010000000000ea305500000000a8ed32320100
00


from CLEOS:
0000000000ea3055
20e1738d2d95d031
owner:
01000000 threshold
01 keys
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf pubkey (curve=00)
0100 weight
00 accounts[]
00 wait_weight[] (!!)

active
01000000
01 keys[]
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
0100 weight
00 accounts[]
00 wait_weight[]









---- that TRANSACTION from cleos:

0d38f85a
5207
26b26cc1
00
00
00
00 cfa[]
01  actions[]
0000000000ea3055
00409e9a2264b89a
01
0000000000ea3055
00000000a8ed3232
66
0000000000ea3055
20e1738d2d95d031
owner
01000000
01
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
0100
00
00 wait_weight[]
active
01000000
01
0002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf
0100
00
00
00 []extensions..







FROM CLEOS:
{"signatures":["SIG_K1_KBt2hPWGeZKNnfjgZUsKUcF3toGhvCfpZJY3S2DoMPBhsuGDJt3s7L4b7eMnaBzAPT6VMXW8jbwfq3KF1RswixMvdbufJ7"],"compression":"none","packed_context_free_data":"","packed_trx":"523df85add116f5b6df200000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232660000000000ea3055708c31c618638c3101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000000"}

FROM eos-go:
{"signatures":["SIG_K1_KApv9d2E2hMGQ74fsPz63gSzuB9zHxXw72pFWqoLgxVhxWuVG5Z3e1hzXtAdDbjicghSUQnqKZSJKopJ752dw77PjFTM5v"],"compression":"none","packed_context_free_data":"","packed_trx":"073ef85a4713ec1f28d000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32326d0000000000ea3055708c31c618638c3101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf010000000000000000000000"}

{"signatures":["SIG_K1_K8tnFMSwdkRjxGwrxpTeWZQdASqBdYZR74vPqmntQmBkUijZtFD3iyy5GF68HRPhaG8hXCYjrTr49zkqSdiaEaLgNsyBjb"],"compression":"none","packed_context_free_data":"","packed_trx":"533df85ade11bc12389900000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327f0000000000ea3055708c31c618638c3101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf010000000100000000010000000000ea305500000000a8ed323201000000"}
