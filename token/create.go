package token

import eos "github.com/eoscanada/eos-go"

func NewCreate(issuer eos.AccountName, maxSupply eos.Asset, canFreeze, canRecall, canWhitelist bool) *eos.Action {
	return &eos.Action{
		Account: AN("eosio.token"),
		Name:    ActN("create"),
		Authorization: []eos.PermissionLevel{
			{Actor: AN("eosio.token"), Permission: PN("active")},
		},
		Data: eos.NewActionData(Create{
			Issuer:        issuer,
			MaximumSupply: maxSupply,
			CanFreeze:     canFreeze,
			CanRecall:     canRecall,
			CanWhitelist:  canWhitelist,
		}),
	}
}

// Create represents the `create` struct on the `eosio.token` contract.
type Create struct {
	Issuer        eos.AccountName `json:"issuer"`
	MaximumSupply eos.Asset       `json:"maximum_supply"`
	CanFreeze     bool            `json:"can_freeze"`
	CanRecall     bool            `json:"can_recall"`
	CanWhitelist  bool            `json:"can_whitelist"`
}
