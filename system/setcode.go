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
		return nil, fmt.Errorf("unmarshal ABI file: %s", err)
	}

	abiPacked, err := eos.MarshalBinary(abiDef)
	if err != nil {
		return nil, fmt.Errorf("packing ABI: %s", err)
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
				ABI:     eos.HexBytes(abiPacked),
			}),
		},
	}
	return &eos.Transaction{Actions: actions}, nil
}

// SetCode represents the hard-coded `setcode` action.
type SetCode struct {
	Account   eos.AccountName `json:"account"`
	VMType    byte            `json:"vmtype"`
	VMVersion byte            `json:"vmversion"`
	Code      eos.HexBytes    `json:"code"`
}

// SetABI represents the hard-coded `setabi` action.
type SetABI struct {
	Account eos.AccountName `json:"account"`
	ABI     eos.HexBytes    `json:"abi"`
}
