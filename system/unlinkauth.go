package system

import "github.com/eoscanada/eos-go"

// NewUnlinkAuth creates an action from the `eosio.system` contract
// called `unlinkauth`.
//
// `unlinkauth` detaches a previously set permission from a
// `code::actionName`. See `linkauth`.
func NewUnlinkAuth(account, code eos.AccountName, actionName eos.ActionName) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("unlinkauth"),
		Authorization: []eos.PermissionLevel{
			{account, eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(UnlinkAuth{
			Account: account,
			Code:    code,
			Type:    actionName,
		}),
	}

	return a
}

// UnlinkAuth represents the native `unlinkauth` action, through the
// system contract.
type UnlinkAuth struct {
	Account eos.AccountName `json:"account"`
	Code    eos.AccountName `json:"code"`
	Type    eos.ActionName  `json:"type"`
}
