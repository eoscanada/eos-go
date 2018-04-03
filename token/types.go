package token

import eos "github.com/eosioca/eosapi"

// Create represents the `create` struct on the `eosio.token` contract.
type Create struct {
	Issuer        eos.AccountName `json:"issuer"`
	MaximumSupply eos.Asset       `json:"maximum_supply"`
	CanFreeze     bool            `json:"can_freeze"`
	CanRecall     bool            `json:"can_recall"`
	CanWhitelist  bool            `json:"can_whitelist"`
}

// Issue represents the `issue` struct on the `eosio.token` contract.
type Issue struct {
	To       eos.AccountName `json:"to"`
	Quantity eos.Asset       `json:"quantity"`
	Memo     string          `json:"memo"`
}

// Transfer represents the `transfer` struct on `eosio.token` contract.
type Transfer struct {
	From     eos.AccountName `json:"from"`
	To       eos.AccountName `json:"to"`
	Quantity eos.Asset       `json:"quantity"`
	Memo     string          `json:"memo"`
}
