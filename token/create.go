package token

import eos "github.com/eoscanada/eos-go"

func NewCreate(issuer eos.AccountName, maxSupply eos.Asset) *eos.Action {
	return &eos.Action{
		Account: AN("eosio.token"),
		Name:    ActN("create"),
		Authorization: []eos.PermissionLevel{
			{Actor: AN("eosio.token"), Permission: PN("active")},
		},
		ActionData: eos.NewActionData(Create{
			Issuer:        issuer,
			MaximumSupply: maxSupply,
		}),
	}
}

// Create represents the `create` struct on the `eosio.token` contract.
type Create struct {
	Issuer        eos.AccountName `json:"issuer"`
	MaximumSupply eos.Asset       `json:"maximum_supply"`
}
