package token

import eos "github.com/eosioca/eosapi"

func NewCreate(issuer eos.AccountName, maxSupply eos.Asset, canFreeze, canRecall, canWhitelist bool) *eos.Action {
	return &eos.Action{
		Account: AN("eosio.token"),
		Name:    ActN("create"),
		Authorization: []eos.PermissionLevel{
			{Actor: AN("eosio.token"), Permission: PN("active")},
		},
		Data: Create{
			Issuer:        issuer,
			MaximumSupply: maxSupply,
			CanFreeze:     canFreeze,
			CanRecall:     canRecall,
			CanWhitelist:  canWhitelist,
		},
	}
}
