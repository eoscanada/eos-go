package rex

import (
	eos "github.com/eoscanada/eos-go"
)

func NewUnstakeToREX(
	owner eos.AccountName,
	receiver eos.AccountName,
	fromNet eos.Asset,
	fromCPU eos.Asset,
) *eos.Action {
	return &eos.Action{
		Account: REXAN,
		Name:    ActN("unstaketorex"),
		Authorization: []eos.PermissionLevel{
			{Actor: owner, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(UnstakeToREX{
			Owner:    owner,
			Receiver: receiver,
			FromNet:  fromNet,
			FromCPU:  fromCPU,
		}),
	}
}

type UnstakeToREX struct {
	Owner    eos.AccountName
	Receiver eos.AccountName
	FromNet  eos.Asset
	FromCPU  eos.Asset
}
