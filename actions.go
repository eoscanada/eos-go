package eos

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		HexData:  []byte{},
		Data:     obj,
		toServer: true,
	}
}
func (a *ActionData) SetToServer(toServer bool) {
	// FIXME: let's clarify this design, make it simpler..
	// toServer doesn't speak of the intent..
	a.toServer = toServer
}

//  jsonActionToServer represents what /v1/chain/push_transaction
//  expects, which isn't allllways the same everywhere.
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
	println(fmt.Sprintf("MarshalJSON toServer? %t", a.toServer))

	if a.toServer {
		buf := new(bytes.Buffer)
		encoder := NewEncoder(buf)

		println("MarshalJSON, encoding action data to binary")
		if err := encoder.Encode(a.ActionData.Data); err != nil {
			return nil, err
		}
		data := buf.Bytes()
		println("MarshalJSON data length : ", len(data)) /**/

		return json.Marshal(&jsonActionToServer{
			Account:       a.Account,
			Name:          a.Name,
			Authorization: a.Authorization,
			Data:          data,
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
