package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewDefundCPULoan(from eos.AccountName, loanNumber uint64, amount eos.Asset) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("defcpuloan"),
		Authorization: []eos.PermissionLevel{
			{Actor: from, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(DefundCPULoan{
			From:       from,
			LoanNumber: loanNumber,
			Amount:     amount,
		}),
	}
}

type DefundCPULoan struct {
	From       eos.AccountName
	LoanNumber uint64
	Amount     eos.Asset
}
