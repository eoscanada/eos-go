package unregd

import "github.com/eoscanada/eos-go"

func NewAdd(ethAccount string, balance eos.Asset) *eos.Action {
	action := &eos.Action{
		Account: eos.AccountName("eosio.unregd"),
		Name:    eos.ActionName("add"),
		Authorization: []eos.PermissionLevel{
			{eos.AccountName("eosio.unregd"), eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Add{
			EthereumAddress: ethAccount,
			Balance:         balance,
		}),
	}
	return action
}

type Add struct {
	EthereumAddress string    `json:"ethereum_account"`
	Balance         eos.Asset `json:"balance"`
}
