package system

import (
	eos "github.com/eosioca/eosapi"
	"github.com/eosioca/eosapi/ecc"
)

// NewSetPriv returns a `setpriv` action that lives on the
// `eosio.bios` contract. It should exist only when booting a new
// network, as it is replaced using the `eos-bios` boot process by the
// `eosio.system` contract.
func NewRegProducer(producer eos.AccountName, producerKey ecc.PublicKey, params eos.EOSIOParameters) *eos.Action {
	return &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("regproducer"),
		Authorization: []eos.PermissionLevel{
			{Actor: producer, Permission: PN("active")},
		},
		Data: RegProducer{
			Producer:    producer,
			ProducerKey: []byte(producerKey),
			Prefs:       params,
		},
	}
}
