package system

import eos "github.com/eosioca/eosapi"

// NewSetPriv returns a `setpriv` action that lives on the
// `eosio.bios` contract. It should exist only when booting a new
// network, as it is replaced using the `eos-bios` boot process by the
// `eosio.system` contract.
func NewSetPriv(account eos.AccountName) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("setpriv"),
		Authorization: []eos.PermissionLevel{
			{AN("eosio"), PN("active")},
		},
		Data: SetPriv{
			Account: account,
			IsPriv:  0x01,
		},
	}
	return a
}
