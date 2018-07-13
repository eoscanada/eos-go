package system

import (
	eos "github.com/eoscanada/eos-go"
)

// NewSellRAM will sell at current market price a given number of
// bytes of RAM.
func NewSellRAM(account eos.AccountName, bytes uint64) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("sellram"),
		Authorization: []eos.PermissionLevel{
			{Actor: account, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(SellRAM{
			Account: account,
			Bytes:   bytes,
		}),
	}
	return a
}

// SellRAM represents the `eosio.system::sellram` action.
type SellRAM struct {
	Account eos.AccountName `json:"account"`
	Bytes   uint64          `json:"bytes"`
}
