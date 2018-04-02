package system

import (
	"encoding/json"
	"io/ioutil"

	eos "github.com/eosioca/eosapi"
)

func NewSetCodeTx(account eos.AccountName, wasmPath, abiPath string) (out *eos.Transaction, err error) {
	codeContent, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return nil, err
	}

	abiContent, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return nil, err
	}

	var abiDef eos.ABI
	if err := json.Unmarshal(abiContent, &abiDef); err != nil {
		return nil, err
	}

	actions := []*eos.Action{
		{
			Account: AN("eosio"),
			Name:    ActN("setcode"),
			Authorization: []eos.PermissionLevel{
				{account, PN("active")},
			},
			Data: SetCode{
				Account:   account,
				VMType:    0,
				VMVersion: 0,
				Code:      eos.HexBytes(codeContent),
			},
		},
		{
			Account: AN("eosio"),
			Name:    ActN("setabi"),
			Authorization: []eos.PermissionLevel{
				{account, eos.PermissionName("active")},
			},
			Data: SetABI{
				Account: account,
				ABI:     abiDef,
			},
		},
	}
	return &eos.Transaction{Actions: actions}, nil
}
