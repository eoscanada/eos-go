package token

import eos "github.com/eoscanada/eos-go"

func init() {
	eos.RegisterAction(AN("eosio.token"), ActN("transfer"), &Transfer{})
	eos.RegisterAction(AN("eosio.token"), ActN("issue"), &Issue{})
}
