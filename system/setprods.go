package system

import eos "github.com/eoscanada/eos-go"

// NewSetPriv returns a `setpriv` action that lives on the
// `eosio.bios` contract. It should exist only when booting a new
// network, as it is replaced using the `eos-bios` boot process by the
// `eosio.system` contract.
func NewSetProds(producers []ProducerKey) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("setprods"),
		Authorization: []eos.PermissionLevel{
			{Actor: AN("eosio"), Permission: PN("active")},
		},
		ActionData: eos.NewActionData(SetProds{
			Schedule: producers,
		}),
	}
	return a
}
