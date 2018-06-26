package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewRemove is an action representing a simple remove to be broadcast
// through the chain network.
func NewRemove(account eos.AccountName, postUUID string) *eos.Action {
	a := &eos.Action{
		Account: AN("eosforumtest"),
		Name:    ActN("remove"),
		Authorization: []eos.PermissionLevel{
			{Actor: account, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Remove{
			Account:  account,
			PostUUID: postUUID,
		}),
	}
	return a
}

// Remove represents the `eosforumtest::remove` action.
type Remove struct {
	Account  eos.AccountName `json:"account"`
	PostUUID string          `json:"post_uuid"`
}
