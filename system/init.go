package system

import (
	"github.com/eoscanada/eos-go"
)

func init() {
	eos.RegisterAction(AN("eosio"), ActN("setcode"), SetCode{})
	eos.RegisterAction(AN("eosio"), ActN("setabi"), SetABI{})
	eos.RegisterAction(AN("eosio"), ActN("newaccount"), NewAccount{})
	eos.RegisterAction(AN("eosio"), ActN("delegatebw"), DelegateBW{})
	eos.RegisterAction(AN("eosio"), ActN("undelegatebw"), UndelegateBW{})
	eos.RegisterAction(AN("eosio"), ActN("refund"), Refund{})
	eos.RegisterAction(AN("eosio"), ActN("regproducer"), RegProducer{})
	eos.RegisterAction(AN("eosio"), ActN("unregprod"), UnregProducer{})
	eos.RegisterAction(AN("eosio"), ActN("regproxy"), RegProxy{})
	eos.RegisterAction(AN("eosio"), ActN("unregproxy"), UnregProxy{})
	eos.RegisterAction(AN("eosio"), ActN("voteproducer"), VoteProducer{})
	eos.RegisterAction(AN("eosio"), ActN("claimrewards"), ClaimRewards{})
	eos.RegisterAction(AN("eosio"), ActN("buyram"), BuyRAM{})
	eos.RegisterAction(AN("eosio"), ActN("buyrambytes"), BuyRAMBytes{})
	// eos.RegisterAction(AN("eosio"), ActN("nonce"), &Nonce{})
}

var AN = eos.AN
var PN = eos.PN
var ActN = eos.ActN
