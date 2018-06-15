package system

import (
	eos "github.com/eoscanada/eos-go"
)

func NewBuyRAM(payer, receiver eos.AccountName, eosQuantity uint64) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("buyram"),
		Authorization: []eos.PermissionLevel{
			{Actor: payer, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(BuyRAM{
			Payer:    payer,
			Receiver: receiver,
			Quantity: eos.NewEOSAsset(int64(eosQuantity)),
		}),
	}
	return a
}

// BuyRAM represents the `eosio.system::buyram` action.
type BuyRAM struct {
	Payer    eos.AccountName `json:"payer"`
	Receiver eos.AccountName `json:"receiver"`
	Quantity eos.Asset       `json:"quant"` // specified in EOS
}
