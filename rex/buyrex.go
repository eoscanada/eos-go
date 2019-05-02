package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewBuyREX(from eos.AccountName, amount eos.AccountName) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("buyrex"),
		Authorization: []eos.PermissionLevel{
			{Actor: from, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(BuyREX{
			From:   from,
			Amount: amount,
		}),
	}
}

type BuyREX struct {
	From   eos.AccountName
	Amount eos.AccountName
}
