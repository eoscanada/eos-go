package eosapi

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"reflect"
)

// See the other ones here: /home/abourget/build/eos3/libraries/chain/include/eosio/chain/contracts/types.hpp:203

type Transfer struct {
	From     AccountName `json:"from"`
	To       AccountName `json:"to"`
	Quantity uint64      `json:"quantity"`
	Memo     string      `json:"memo"`
}

type Issue struct {
	To       AccountName `json:"to"`
	Quantity uint64      `json:"quantity" struc:"uint64,little"`
}

type SetCode struct {
	Account   AccountName `json:"account"`
	VMType    byte        `json:"vmtype"`
	VMVersion byte        `json:"vmversion"`
	Code      HexBytes    `json:"bytes"`
}

type NewAccount struct {
	Creator  AccountName `json:"creator"`
	Name     AccountName `json:"name"`
	Owner    Authority   `json:"owner"`
	Active   Authority   `json:"active"`
	Recovery Authority   `json:"recovery"`
}

type Action struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
	Data          ActionData        `json:"data,omitempty"` // as HEX when we receive it.. FIXME: decode from hex directly.. and encode back plz!
}

type action struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
}

type ActionData interface{}

func (a *Action) UnmarshalBinaryRead(r io.Reader) error {
	length, err := binary.ReadUvarint(&ByteReader{r})
	if err != nil {
		return err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}

	actionMap := registeredActions[a.Account]
	if actionMap == nil {
		return nil
	}

	objMap := actionMap[a.Name]
	if objMap == nil {
		return nil
	}

	obj := reflect.New(reflect.TypeOf(objMap))

	err = UnmarshalBinary(data, &obj)
	if err != nil {
		return err
	}

	a.Data = obj.Elem().Interface()

	return nil
}

var registeredActions = map[AccountName]map[ActionName]reflect.Type{}

func init() {
	RegisterAction(AccountName("eosio"), ActionName("transfer"), &Transfer{})
	RegisterAction(AccountName("eosio"), ActionName("issue"), &Issue{})
}

// Registers Action objects..
func RegisterAction(accountName AccountName, actionName ActionName, obj interface{}) {
	// TODO: lock or som'th.. unless we never call after boot time..
	if registeredActions[accountName] == nil {
		registeredActions[accountName] = make(map[ActionName]reflect.Type)
	}
	registeredActions[accountName][actionName] = reflect.ValueOf(obj).Type()
}

func (a *Action) UnmarshalJSON(v []byte) (err error) {
	// load Account, Name, Authorization, Data
	// and then unpack other fields in a struct based on `Name` and `AccountName`..
	var newAct *action
	if err = json.Unmarshal(v, &newAct); err != nil {
		return
	}

	a.Account = newAct.Account
	a.Name = newAct.Name
	a.Authorization = newAct.Authorization

	actionMap := registeredActions[a.Account]
	if actionMap == nil {
		return nil
	}

	objMap := actionMap[a.Name]
	if objMap == nil {
		return nil
	}

	obj := reflect.New(reflect.TypeOf(objMap))
	err = json.Unmarshal(v, &obj)
	if err != nil {
		return err
	}

	a.Data = obj.Elem().Interface()

	return nil
}

func (a *Action) MarshalBinary() ([]byte, error) {
	// marshal binary a short action, then data
	common, err := MarshalBinary(&action{
		Account:       a.Account,
		Name:          a.Name,
		Authorization: a.Authorization,
	})
	if err != nil {
		return nil, err
	}

	var data []byte
	if a.Data != nil {
		data, err = MarshalBinary(a.Data)
		if err != nil {
			return nil, err
		}
	}

	varint := make([]byte, 4, 4)
	varintLen := binary.PutUvarint(varint, uint64(len(data)))
	common = append(common, varint[:varintLen]...)
	common = append(common, data...)
	return common, nil
}

func (a *Action) MarshalJSON() ([]byte, error) {
	// Start with the base-line Action fields.

	cnt, err := json.Marshal(&action{
		Account:       a.Account,
		Name:          a.Name,
		Authorization: a.Authorization,
	})
	if err != nil {
		return nil, err
	}

	var keys1 map[string]interface{}
	err = json.Unmarshal(cnt, &keys1)
	if err != nil {
		return nil, err
	}

	data, err := MarshalBinary(a)

	keys1["data"] = hex.EncodeToString(data)

	return json.Marshal(keys1)
}
