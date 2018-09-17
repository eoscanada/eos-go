package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// Status is an action to set a status update for a given account on the forum contract.
func NewStatus(account eos.AccountName, content string) *eos.Action {
	a := &eos.Action{
		Account: ForumAN,
		Name:    ActN("status"),
		Authorization: []eos.PermissionLevel{
			{Actor: account, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Status{
			Account: account,
			Content: content,
		}),
	}
	return a
}

// Status represents the `eosio.forum::status` action.
type Status struct {
	Account eos.AccountName `json:"account_name"`
	Content string          `json:"content"`
}
