package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewDeposit(owner eos.AccountName, amount eos.Asset) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("deposit"),
		Authorization: []eos.PermissionLevel{
			{Actor: owner, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Deposit{
			Owner:  owner,
			Amount: amount,
		}),
	}
}

type Deposit struct {
	Owner  eos.AccountName
	Amount eos.Asset
}
