package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewFundCPULoan(from eos.AccountName, loanNumber uint64, payment eos.Asset) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("fundcpuloan"),
		Authorization: []eos.PermissionLevel{
			{Actor: from, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(FundCPULoan{
			From:       from,
			LoanNumber: loanNumber,
			Payment:    payment,
		}),
	}
}

type FundCPULoan struct {
	From       eos.AccountName
	LoanNumber uint64
	Payment    eos.Asset
}
