package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewRentCPU(
	from eos.AccountName,
	receiver eos.AccountName,
	loanPayment eos.Asset,
	loanFund eos.Asset,
) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("rentcpu"),
		Authorization: []eos.PermissionLevel{
			{Actor: from, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(RentCPU{
			From:        from,
			Receiver:    receiver,
			LoanPayment: loanPayment,
			LoanFund:    loanFund,
		}),
	}
}

type RentCPU struct {
	From        eos.AccountName
	Receiver    eos.AccountName
	LoanPayment eos.Asset
	LoanFund    eos.Asset
}
