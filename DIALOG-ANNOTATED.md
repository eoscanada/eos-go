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

github.com/eosioca/eosapi/cmd/eosapi
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
