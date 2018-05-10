package system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	eos "github.com/eoscanada/eos-go"
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
		return nil, fmt.Errorf("unmarshal ABI file:, %s", err)
	}

	actions := []*eos.Action{
		{
			Account: AN("eosio"),
			Name:    ActN("setcode"),
			Authorization: []eos.PermissionLevel{
				{account, eos.PermissionName("active")},
			},
			ActionData: eos.NewActionData(SetCode{
				Account:   account,
				VMType:    0,
				VMVersion: 0,
				Code:      eos.HexBytes(codeContent),
			}),
		},
		{
			Account: AN("eosio"),
			Name:    ActN("setabi"),
			Authorization: []eos.PermissionLevel{
				{account, eos.PermissionName("active")},
			},
			ActionData: eos.NewActionData(SetABI{
				Account: account,
				ABI:     abiDef,
			}),
		},
	}
	return &eos.Transaction{Actions: actions}, nil
}
