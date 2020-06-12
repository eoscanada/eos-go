package boot

import "github.com/eoscanada/eos-go"

//accts/dfu/se44shine/tables/posts/eo/scanadacom.json
//accts/{[acc]/[ountName]}/contract.wasm
//accts/{[acc]/[ountName]}/contract.abi
//accts/{[acc]/[ountName]}/resources.json
//accts/{[acc]/[ountName]}/permissions.json

//accts/{[acc]/[ountName]}/tables/[tableName]/{[sco]/[peName]}.json
//Scope, TableName, Contract
//Payer, Data,

//account dfuse.boot setCode wasm.contract
type TableRow struct {
	Payer string       `json:"payer"`
	Data  eos.HexBytes `json:"data"`
}

// base contract that print data length
