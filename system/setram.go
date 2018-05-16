package system

import (
	eos "github.com/eoscanada/eos-go"
)

func NewSetRAM(maxRAMSize uint64) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("setram"),
		Authorization: []eos.PermissionLevel{
			{AN("eosio"), eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(SetRAM{
			MaxRAMSize: maxRAMSize,
		}),
	}
	return a
}

// SetRAM represents the hard-coded `setram` action.
type SetRAM struct {
	MaxRAMSize uint64 `json:"max_ram_size"`
}
