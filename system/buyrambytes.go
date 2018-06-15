package system

import (
	eos "github.com/eoscanada/eos-go"
)

// NewBuyRAMBytes will buy at current market price a given number of
// bytes of RAM, and grant them to the `receiver` account.
func NewBuyRAMBytes(payer, receiver eos.AccountName, bytes uint32) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("buyrambytes"),
		Authorization: []eos.PermissionLevel{
			{Actor: payer, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(BuyRAMBytes{
			Payer:    payer,
			Receiver: receiver,
			Bytes:    bytes,
		}),
	}
	return a
}

// BuyRAMBytes represents the `eosio.system::buyrambytes` action.
type BuyRAMBytes struct {
	Payer    eos.AccountName `json:"payer"`
	Receiver eos.AccountName `json:"receiver"`
	Bytes    uint32          `json:"bytes"`
}
