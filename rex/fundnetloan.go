package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewFundNetLoan(from eos.AccountName, loanNumber uint64, payment eos.Asset) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("fundnetloan"),
		Authorization: []eos.PermissionLevel{
			{Actor: from, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(FundNetLoan{
			From:       from,
			LoanNumber: loanNumber,
			Payment:    payment,
		}),
	}
}

type FundNetLoan struct {
	From       eos.AccountName
	LoanNumber uint64
	Payment    eos.Asset
}
