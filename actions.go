package eos

import (
	"bytes"
	"encoding/json"
)

// See: libraries/chain/include/eosio/chain/contracts/types.hpp:203
// See: build/contracts/eosio.system/eosio.system.abi

// SetCode represents the hard-coded `setcode` action.
type SetCode struct {
	Account   AccountName `json:"account"`
	VMType    byte        `json:"vmtype"`
	VMVersion byte        `json:"vmversion"`
	Code      HexBytes    `json:"bytes"`
}

// SetABI represents the hard-coded `setabi` action.
type SetABI struct {
	Account AccountName `json:"account"`
	ABI     ABI         `json:"abi"`
}

// NewAccount represents the hard-coded `newaccount` action.
type NewAccount struct {
	Creator  AccountName `json:"creator"`
	Name     AccountName `json:"name"`
	Owner    Authority   `json:"owner"`
	Active   Authority   `json:"active"`
	Recovery Authority   `json:"recovery"`
}

// Action
type Action struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
	ActionData
}

type ActionData struct {
	HexData  HexBytes    `json:"hex_data,omitempty"`
	Data     interface{} `json:"data,omitempty" eos:"-"`
	abi      []byte      // TBD: we could use the ABI to decode in obj
	toServer bool
}

func NewActionData(obj interface{}) ActionData {
	return ActionData{
		HexData:  []byte(""),
		Data:     obj,
		toServer: true,
	}
}

type jsonActionToServer struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
	Data          HexBytes          `json:"data,omitempty"`
}

type jsonActionFromServer struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
	Data          interface{}       `json:"data,omitempty"`
	HexData       HexBytes          `json:"hex_data,omitempty"`
}

func (a *Action) MarshalJSON() ([]byte, error) {

	if a.toServer { //sending action to server
		var err error
		buf := new(bytes.Buffer)
		encoder := NewEncoder(buf)
		encoder.Encode(a.ActionData.Data)

		if err != nil {
			return nil, err
		}
		data := buf.Bytes()

		return json.Marshal(&jsonActionToServer{
			Account:       a.Account,
			Name:          a.Name,
			Authorization: a.Authorization,
			Data:          HexBytes(data),
		})
	}

	return json.Marshal(&jsonActionFromServer{
		Account:       a.Account,
		Name:          a.Name,
		Authorization: a.Authorization,
		HexData:       a.HexData,
		Data:          a.Data,
	})
}
