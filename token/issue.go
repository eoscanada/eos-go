package token

import eos "github.com/eoscanada/eos-go"

func NewIssue(to eos.AccountName, quantity eos.Asset, memo string) *eos.Action {
	return &eos.Action{
		Account: AN("eosio.token"),
		Name:    ActN("issue"),
		Authorization: []eos.PermissionLevel{
			{Actor: AN("eosio"), Permission: PN("active")},
		},
		ActionData: eos.NewActionData(Issue{
			To:       to,
			Quantity: quantity,
			Memo:     memo,
		}),
	}
}

// Issue represents the `issue` struct on the `eosio.token` contract.
type Issue struct {
	To       eos.AccountName `json:"to"`
	Quantity eos.Asset       `json:"quantity"`
	Memo     string          `json:"memo"`
}
