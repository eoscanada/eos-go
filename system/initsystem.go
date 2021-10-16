package system

import (
	eos "github.com/eoscanada/eos-go"
)

// NewInitSystem returns a `init` action that lives on the
// `eosio.system` contract.
func NewInitSystem(version eos.Varuint32, core eos.Symbol) *eos.Action {
	return &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("init"),
		Authorization: []eos.PermissionLevel{
			{
				Actor:      AN("eosio"),
				Permission: eos.PermissionName("active"),
			},
		},
		ActionData: eos.NewActionData(Init{
			Version: version,
			Core:    core,
		}),
	}
}

// Init represents the `eosio.system::init` action
type Init struct {
	Version eos.Varuint32 `json:"version"`
	Core    eos.Symbol    `json:"core"`
}
