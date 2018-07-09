package system

import "github.com/eoscanada/eos-go"

// NewLinkAuth creates an action from the `eosio.system` contract
// called `linkauth`.
//
// `linkauth` allows you to attach certain permission to the given
// `code::actionName`. With this set on-chain, you can use the
// `requiredPermission` to sign transactions for `code::actionName`
// and not rely on your `active` (which might be more sensitive as it
// can sign anything) for the given operation.
func NewLinkAuth(account, code eos.AccountName, actionName eos.ActionName, requiredPermission eos.PermissionName) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("linkauth"),
		Authorization: []eos.PermissionLevel{
			{account, eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(LinkAuth{
			Account:     account,
			Code:        code,
			Type:        actionName,
			Requirement: requiredPermission,
		}),
	}

	return a
}

// LinkAuth represents the native `linkauth` action, through the
// system contract.
type LinkAuth struct {
	Account     eos.AccountName    `json:"account"`
	Code        eos.AccountName    `json:"code"`
	Type        eos.ActionName     `json:"type"`
	Requirement eos.PermissionName `json:"requirement"`
}
