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
			Auth:       authority,
		}),
	}

	return a
}

// UpdateAuth represents the hard-coded `updateauth` action.
//
// If you change the `active` permission, `owner` is the required parent.
//
// If you change the `owner` permission, there should be no parent.
type UpdateAuth struct {
	Account    eos.AccountName    `json:"account"`
	Permission eos.PermissionName `json:"permission"`
	Parent     eos.PermissionName `json:"parent"`
	Auth       eos.Authority      `json:"auth"`
}
