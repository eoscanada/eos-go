package system

import (
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

// NewSetPriv returns a `setpriv` action that lives on the
// `eosio.bios` contract. It should exist only when booting a new
// network, as it is replaced using the `eos-bios` boot process by the
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
