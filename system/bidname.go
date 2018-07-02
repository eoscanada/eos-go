package system

import (
	eos "github.com/eoscanada/eos-go"
)

func NewBidname(bidder, newname eos.AccountName, bid eos.Asset) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("bidname"),
		Authorization: []eos.PermissionLevel{
			{Actor: bidder, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(Bidname{
			Bidder:  bidder,
			Newname: newname,
			Bid:     bid,
		}),
	}
	return a
}

// Bidname represents the `eosio.system_contract::bidname` action.
type Bidname struct {
	Bidder  eos.AccountName `json:"bidder"`
	Newname eos.AccountName `json:"newname"`
	Bid     eos.Asset       `json:"bid"` // specified in EOS
}
