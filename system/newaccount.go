package system

import (
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

// NewNewAccount returns a `newaccount` action that lives on the
// `eosio.system` contract.
func NewNewAccount(creator, newAccount eos.AccountName, publicKey ecc.PublicKey) *eos.Action {
	return &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("newaccount"),
		Authorization: []eos.PermissionLevel{
			{Actor: creator, Permission: PN("active")},
		},
		Data: eos.NewActionData(NewAccount{
			Creator: creator,
			Name:    newAccount,
			Owner: eos.Authority{
				Threshold: 1,
				Keys: []eos.KeyWeight{
					eos.KeyWeight{
						PublicKey: publicKey,
						Weight:    1,
					},
				},
			},
			Active: eos.Authority{
				Threshold: 1,
				Keys: []eos.KeyWeight{
					eos.KeyWeight{
						PublicKey: publicKey,
						Weight:    1,
					},
				},
			},
			Recovery: eos.Authority{
				Threshold: 1,
				Accounts: []eos.PermissionLevelWeight{
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{creator, PN("active")},
						Weight:     1,
					},
				},
			},
		}),
	}
}

// NewAccount represents a `newaccount` action on the `eosio.system`
// contract. It is one of the rare ones to be hard-coded into the
// blockchain.
type NewAccount struct {
	Creator  eos.AccountName `json:"creator"`
	Name     eos.AccountName `json:"name"`
	Owner    eos.Authority   `json:"owner"`
	Active   eos.Authority   `json:"active"`
	Recovery eos.Authority   `json:"recovery"`
}
