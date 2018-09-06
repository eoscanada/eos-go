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
	eos.RegisterAction(AN("eosio"), ActN("voteproducer"), VoteProducer{})
	eos.RegisterAction(AN("eosio"), ActN("claimrewards"), ClaimRewards{})
	eos.RegisterAction(AN("eosio"), ActN("buyram"), BuyRAM{})
	eos.RegisterAction(AN("eosio"), ActN("buyrambytes"), BuyRAMBytes{})
	eos.RegisterAction(AN("eosio"), ActN("linkauth"), LinkAuth{})
	eos.RegisterAction(AN("eosio"), ActN("unlinkauth"), UnlinkAuth{})
	eos.RegisterAction(AN("eosio"), ActN("deleteauth"), DeleteAuth{})
	eos.RegisterAction(AN("eosio"), ActN("rmvproducer"), RemoveProducer{})
	eos.RegisterAction(AN("eosio"), ActN("setprods"), SetProds{})
	eos.RegisterAction(AN("eosio"), ActN("setpriv"), SetPriv{})
	eos.RegisterAction(AN("eosio"), ActN("canceldelay"), CancelDelay{})
	eos.RegisterAction(AN("eosio"), ActN("bidname"), Bidname{})
	// eos.RegisterAction(AN("eosio"), ActN("nonce"), &Nonce{})
	eos.RegisterAction(AN("eosio"), ActN("sellram"), SellRAM{})
	eos.RegisterAction(AN("eosio"), ActN("updateauth"), UpdateAuth{})
	eos.RegisterAction(AN("eosio"), ActN("setramrate"), SetRAMRate{})
	eos.RegisterAction(AN("eosio"), ActN("setalimits"), Setalimits{})
}

var AN = eos.AN
var PN = eos.PN
var ActN = eos.ActN
