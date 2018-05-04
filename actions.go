package eos

import (
	"bytes"
	"encoding/hex"
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
	Data          ActionData        `json:"data,omitempty"` // as HEX when we receive it.. FIXME: decode from hex directly.. and encode back plz!
}

func (a Action) Obj() interface{} { // Payload ? ActionData ? GetData ?
	return a.Data.Obj
}

type ActionData struct {
	HexBytes
	Obj interface{} // potentially unpacked from the Actions registry mapped through `RegisterAction`.
	abi []byte      // TBD: we could use the ABI to decode in obj
}

func NewActionData(obj interface{}) ActionData {
	return ActionData{
		HexBytes: []byte(""),
		Obj:      obj,
	}
}

func (a *ActionData) UnmarshalJSON(v []byte) (err error) {
	// Unmarshal from the JSON format ?  We'd need it to be registered.. but we can't hook into the JSON
	// lib to read the current action above.. we'll need to defer loading
	// Either keep as json.RawMessage, or as map[string]interface{}
	a.HexBytes = v
	return nil
}

func (a ActionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Obj)
}

type jsonAction struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
	Data          HexBytes          `json:"data"`
}

func (a *Action) UnmarshalJSON(v []byte) (err error) {
	// load Account, Name, Authorization, Data
	// and then unpack other fields in a struct based on `Name` and `AccountName`..
	var newAct jsonAction
	if err = json.Unmarshal(v, &newAct); err != nil {
		return
	}

	a.Account = newAct.Account
	a.Name = newAct.Name
	a.Authorization = newAct.Authorization
	a.Data.HexBytes = newAct.Data

	// err = UnmarshalBinaryWithAction([]byte(newAct.Data), &a.Data, *a)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (a *Action) MarshalJSON() ([]byte, error) {
	var data HexBytes
	if a.Data.Obj == nil {
		data = a.Data.HexBytes
	} else {
		var err error

		buf := new(bytes.Buffer)
		encoder := NewEncoder(buf)
		fmt.Println("Will encode action.Data.obj")
		encoder.Encode(a.Data.Obj)

		if err != nil {
			return nil, err
		}
		data = buf.Bytes()
		fmt.Print("-------->", hex.EncodeToString(data))
	}

	return json.Marshal(&jsonAction{
		Account:       a.Account,
		Name:          a.Name,
		Authorization: a.Authorization,
		Data:          HexBytes(data),
	})
}
