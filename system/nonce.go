package system

import "github.com/eoscanada/eos-go"

// NewNonce returns a `nonce` action that lives on the
// `eosio.bios` contract. It should exist only when booting a new
// network, as it is replaced using the `eos-bios` boot process by the
// `eosio.system` contract.
func NewNonce(nonce string) *eos.Action {
	a := &eos.Action{
		Account:       AN("eosio"),
		Name:          ActN("nonce"),
		Authorization: []eos.PermissionLevel{
			//{Actor: AN("eosio"), Permission: PN("active")},
		},
		ActionData: eos.NewActionData(Nonce{
			Value: nonce,
		}),
	}
	return a
}
