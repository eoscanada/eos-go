package system

import (
	eos "github.com/eoscanada/eos-go"
)

// NewRemoveProducer returns a `rmvproducer` action that lives on the
// `eosio.system` contract.  This is to be called by the consortium of
// BPs, to oust a BP from its place.  If you want to unregister
// yourself as a BP, use `unregprod`.
func NewRemoveProducer(producer eos.AccountName) *eos.Action {
	return &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("rmvproducer"),
		Authorization: []eos.PermissionLevel{
			{Actor: AN("eosio"), Permission: PN("active")},
		},
		ActionData: eos.NewActionData(RemoveProducer{
			Producer: producer,
		}),
	}
}

// RemoveProducer represents the `eosio.system::rmvproducer` action
type RemoveProducer struct {
	Producer eos.AccountName `json:"producer"`
}
