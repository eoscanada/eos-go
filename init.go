package eos

func init() {

	RegisterAction(AN("eosio"), ActN("setcode"), &SetCode{})
	RegisterAction(AN("eosio"), ActN("setabi"), &SetABI{})
	RegisterAction(AN("eosio"), ActN("newaccount"), &NewAccount{})

}
