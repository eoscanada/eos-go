package token

import eos "github.com/eosioca/eosapi"

func NewTransfer(from, to eos.AccountName, quantity eos.Asset, memo string) *eos.Action {
	return &eos.Action{
		Account: AN("eosio.token"),
		Name:    ActN("transfer"),
		Authorization: []eos.PermissionLevel{
			{Actor: from, Permission: PN("active")},
		},
		Data: Transfer{
			From:     from,
			To:       to,
			Quantity: quantity,
			Memo:     memo,
		},
	}
}
