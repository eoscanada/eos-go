package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewDefundNetLoan(from eos.AccountName, loanNumber uint64, amount eos.Asset) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("defnetloan"),
		Authorization: []eos.PermissionLevel{
			{Actor: from, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(DefundNetLoan{
			From:       from,
			LoanNumber: loanNumber,
			Amount:     amount,
		}),
	}
}

type DefundNetLoan struct {
	From       eos.AccountName
	LoanNumber uint64
	Amount     eos.Asset
}
