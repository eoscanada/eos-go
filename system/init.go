package system

import (
	"github.com/eoscanada/eos-go"
)

func init() {
	eos.RegisterAction(AN("eosio"), ActN("setcode"), &SetCode{})
	eos.RegisterAction(AN("eosio"), ActN("setabi"), &SetABI{})
	eos.RegisterAction(AN("eosio"), ActN("newaccount"), &NewAccount{})
	// eos.RegisterAction(AN("eosio"), ActN("delegatebw"), &DelegateBW{})
	// eos.RegisterAction(AN("eosio"), ActN("undelegatebw"), &UndelegateBW{})
	// eos.RegisterAction(AN("eosio"), ActN("refund"), &Refund{})
	// eos.RegisterAction(AN("eosio"), ActN("regproducer"), &RegisterProducer{})
	// eos.RegisterAction(AN("eosio"), ActN("unregprod"), &UnregisterProducer{})
	// eos.RegisterAction(AN("eosio"), ActN("regproxy"), &RegisterProxy{})
	// eos.RegisterAction(AN("eosio"), ActN("unregproxy"), &UnregisterProxy{})
	// eos.RegisterAction(AN("eosio"), ActN("voteproducer"), &VoteProducer{})
	// eos.RegisterAction(AN("eosio"), ActN("claimrewards"), &ClaimRewards{})
	// eos.RegisterAction(AN("eosio"), ActN("nonce"), &Nonce{})
}
