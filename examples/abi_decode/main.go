package main

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
)

func main() {
	api := eos.New("https://mainnet.eoscanada.com")
	abiResponse, err := api.GetABI("betdiceadmin")
	if err != nil {
		panic(fmt.Errorf("get ABI: %s", err))
	}

	abi := &abiResponse.ABI()
	tableDef := abi.TableForName(eos.TableName("bet"))

}

func data() string {
	return `3cd003010000000010c4a519ab3ddd3b0b656f73696f2e746f6b656ea08601000000000004454f530000000060000000a0bc62d329859d331232656f776176646c37764b7335597475354524333339373a626a696e766573746f7232313a32656f776176646c37764b7335597475354535f83aa9200e730bfaa085b5aa1df7513b28c52c4ce3e892835c224d778585a8`
}

func abi() string {
	return `
	{
		"account_name": "betdiceadmin",
		"abi": {
		  "version": "eosio::abi/1.0",
		  "structs": [
			{
			  "name": "premiumrefer",
			  "base": "",
			  "fields": [
				{
				  "name": "username",
				  "type": "name"
				}
			  ]
			},
			{
			  "name": "dicenumstats",
			  "base": "",
			  "fields": [
				{
				  "name": "diceNumber",
				  "type": "uint64"
				},
				{
				  "name": "count",
				  "type": "uint64"
				}
			  ]
			},
			{
			  "name": "player",
			  "base": "",
			  "fields": [
				{
				  "name": "accountName",
				  "type": "uint64"
				},
				{
				  "name": "playCount",
				  "type": "uint32"
				}
			  ]
			},
			{
			  "name": "bet",
			  "base": "",
			  "fields": [
				{
				  "name": "gameId",
				  "type": "uint64"
				},
				{
				  "name": "accountName",
				  "type": "uint64"
				},
				{
				  "name": "contractName",
				  "type": "string"
				},
				{
				  "name": "betAsset",
				  "type": "asset"
				},
				{
				  "name": "rollUnder",
				  "type": "uint32"
				},
				{
				  "name": "referer",
				  "type": "name"
				},
				{
				  "name": "seed",
				  "type": "string"
				},
				{
				  "name": "hashSeed",
				  "type": "string"
				},
				{
				  "name": "hashSeedHash",
				  "type": "checksum256"
				}
			  ]
			},
			{
			  "name": "autoid",
			  "base": "",
			  "fields": [
				{
				  "name": "key",
				  "type": "uint64"
				},
				{
				  "name": "nextId",
				  "type": "uint64"
				}
			  ]
			},
			{
			  "name": "symbolvar",
			  "base": "",
			  "fields": [
				{
				  "name": "id",
				  "type": "uint64"
				},
				{
				  "name": "contractName",
				  "type": "uint64"
				},
				{
				  "name": "symbolType",
				  "type": "uint64"
				},
				{
				  "name": "fee",
				  "type": "uint32"
				},
				{
				  "name": "min",
				  "type": "uint64"
				},
				{
				  "name": "max",
				  "type": "uint64"
				}
			  ]
			},
			{
			  "name": "globalvar",
			  "base": "",
			  "fields": [
				{
				  "name": "id",
				  "type": "uint64"
				},
				{
				  "name": "value",
				  "type": "string"
				}
			  ]
			},
			{
			  "name": "account",
			  "base": "",
			  "fields": [
				{
				  "name": "balance",
				  "type": "asset"
				}
			  ]
			},
			{
			  "name": "setglobalvar",
			  "base": "",
			  "fields": [
				{
				  "name": "id",
				  "type": "uint64"
				},
				{
				  "name": "value",
				  "type": "string"
				}
			  ]
			},
			{
			  "name": "betreceipt",
			  "base": "",
			  "fields": [
				{
				  "name": "gamename",
				  "type": "name"
				},
				{
				  "name": "gameId",
				  "type": "uint64"
				},
				{
				  "name": "accountName",
				  "type": "name"
				},
				{
				  "name": "contractName",
				  "type": "string"
				},
				{
				  "name": "betAsset",
				  "type": "asset"
				},
				{
				  "name": "payoutAsset",
				  "type": "asset"
				},
				{
				  "name": "result",
				  "type": "string"
				},
				{
				  "name": "rollUnder",
				  "type": "uint32"
				},
				{
				  "name": "referer",
				  "type": "name"
				},
				{
				  "name": "seed",
				  "type": "string"
				},
				{
				  "name": "hashSeed",
				  "type": "string"
				},
				{
				  "name": "hashSeedHash",
				  "type": "checksum256"
				},
				{
				  "name": "encryptedHash",
				  "type": "string"
				},
				{
				  "name": "numberHash",
				  "type": "string"
				},
				{
				  "name": "diceNumber",
				  "type": "uint32"
				},
				{
				  "name": "now",
				  "type": "uint32"
				}
			  ]
			},
			{
			  "name": "setsymbolvar",
			  "base": "",
			  "fields": [
				{
				  "name": "asset",
				  "type": "asset"
				},
				{
				  "name": "contractName",
				  "type": "name"
				},
				{
				  "name": "fee",
				  "type": "float64"
				},
				{
				  "name": "min",
				  "type": "uint64"
				},
				{
				  "name": "max",
				  "type": "uint64"
				}
			  ]
			},
			{
			  "name": "cancelbet",
			  "base": "",
			  "fields": [
				{
				  "name": "gameId",
				  "type": "uint64"
				},
				{
				  "name": "message",
				  "type": "string"
				}
			  ]
			},
			{
			  "name": "refundbet",
			  "base": "",
			  "fields": [
				{
				  "name": "gameId",
				  "type": "uint64"
				},
				{
				  "name": "message",
				  "type": "string"
				}
			  ]
			},
			{
			  "name": "reveal2",
			  "base": "",
			  "fields": [
				{
				  "name": "hashSeedHash",
				  "type": "checksum256"
				},
				{
				  "name": "encryptedHash",
				  "type": "string"
				},
				{
				  "name": "numberHash",
				  "type": "checksum256"
				}
			  ]
			},
			{
			  "name": "addpremium",
			  "base": "",
			  "fields": [
				{
				  "name": "username",
				  "type": "name"
				}
			  ]
			},
			{
			  "name": "rmpremium",
			  "base": "",
			  "fields": [
				{
				  "name": "username",
				  "type": "name"
				}
			  ]
			}
		  ],
		  "actions": [
			{
			  "name": "setglobalvar",
			  "type": "setglobalvar",
			  "ricardian_contract": ""
			},
			{
			  "name": "betreceipt",
			  "type": "betreceipt",
			  "ricardian_contract": ""
			},
			{
			  "name": "setsymbolvar",
			  "type": "setsymbolvar",
			  "ricardian_contract": ""
			},
			{
			  "name": "cancelbet",
			  "type": "cancelbet",
			  "ricardian_contract": ""
			},
			{
			  "name": "refundbet",
			  "type": "refundbet",
			  "ricardian_contract": ""
			},
			{
			  "name": "reveal2",
			  "type": "reveal2",
			  "ricardian_contract": ""
			},
			{
			  "name": "addpremium",
			  "type": "addpremium",
			  "ricardian_contract": ""
			},
			{
			  "name": "rmpremium",
			  "type": "rmpremium",
			  "ricardian_contract": ""
			}
		  ],
		  "tables": [
			{
			  "name": "premiumrefer",
			  "index_type": "i64",
			  "key_names": [
				"username"
			  ],
			  "key_types": [
				"name"
			  ],
			  "type": "premiumrefer"
			},
			{
			  "name": "dicenumstats",
			  "index_type": "i64",
			  "key_names": [
				"diceNumber"
			  ],
			  "key_types": [
				"uint64"
			  ],
			  "type": "dicenumstats"
			},
			{
			  "name": "player",
			  "index_type": "i64",
			  "key_names": [
				"accountName"
			  ],
			  "key_types": [
				"uint64"
			  ],
			  "type": "player"
			},
			{
			  "name": "bet",
			  "index_type": "i64",
			  "key_names": [
				"gameId"
			  ],
			  "key_types": [
				"uint64"
			  ],
			  "type": "bet"
			},
			{
			  "name": "autoid",
			  "index_type": "i64",
			  "key_names": [
				"key"
			  ],
			  "key_types": [
				"uint64"
			  ],
			  "type": "autoid"
			},
			{
			  "name": "symbolvar",
			  "index_type": "i64",
			  "key_names": [
				"id"
			  ],
			  "key_types": [
				"uint64"
			  ],
			  "type": "symbolvar"
			},
			{
			  "name": "globalvar",
			  "index_type": "i64",
			  "key_names": [
				"id"
			  ],
			  "key_types": [
				"uint64"
			  ],
			  "type": "globalvar"
			},
			{
			  "name": "accounts",
			  "index_type": "i64",
			  "key_names": [
				"balance"
			  ],
			  "key_types": [
				"asset"
			  ],
			  "type": "account"
			}
		  ]
		}
	  }
	`
}
