package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewWithdraw(owner eos.AccountName, amount eos.Asset) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("withdraw"),
		Authorization: []eos.PermissionLevel{
			{Actor: owner, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Withdraw{
			Owner:  owner,
			Amount: amount,
		}),
	}
}

type Withdraw struct {
	Owner  eos.AccountName
	Amount eos.Asset
}
