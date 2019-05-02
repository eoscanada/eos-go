package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewCancelREXorder(owner eos.AccountName) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("cnclrexorder"),
		Authorization: []eos.PermissionLevel{
			{Actor: owner, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(CancelREXorder{
			Owner: owner,
		}),
	}
}

type CancelREXorder struct {
	Owner eos.AccountName
}
