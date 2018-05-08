package system

import "github.com/eoscanada/eos-go"

// NewUpdateAuth creates an action from the `eosio.system` contract
// called `updateauth`.
//
// usingPermission needs to be `owner` if you want to modify the
// `owner` authorization, otherwise `active` will do for the rest.
func NewUpdateAuth(account eos.AccountName, permission, parent eos.PermissionName, authority eos.Authority, usingPermission eos.PermissionName) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("updateauth"),
		Authorization: []eos.PermissionLevel{
			{account, usingPermission},
		},
		ActionData: eos.NewActionData(UpdateAuth{
			Account:    account,
			Permission: permission,
			Parent:     parent,
			Data:       authority,
		}),
	}

	return a
}

/**

42a3be5a00000100e543ae35
0000
0000
00
02

ACTION 1
0000000000ea3055 eosio
0040cbdaa86c52d5 updateauth
01
0000000000ea3055 eosio
00000000a8ed3232 active
1e len
0000000000ea3055 eosio
00000000a8ed3232 active
0000000080ab26a7 owner
00000000  threshold
00
00

ACTION 2
0000000000ea3055 eosio
0040cbdaa86c52d5 updateauth
01
0000000000ea3055
00000000a8ed3232
1e len
0000000000ea3055 eosio
0000000080ab26a7 owner
0000000000000000 NONE (parent empty)
00000000
00
00


0000000000ea3055 eosio
00000000a8ed3232 active
0000000080ab26a7 owner
00000000
00
00
*/
