package system

import (
	eos "github.com/eoscanada/eos-go"
)

func NewSetRAMRate(bytesPerBlock uint16) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("setram"),
		Authorization: []eos.PermissionLevel{
			{AN("eosio"), eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(SetRAMRate{
			BytesPerBlock: bytesPerBlock,
		}),
	}
	return a
}

// SetRAMRate represents the system contract's `setramrate` action.
type SetRAMRate struct {
	BytesPerBlock uint16 `json:"bytes_per_block"`
}
