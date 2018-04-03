package token

import eos "github.com/eosioca/eosapi"

func init() {
	eos.RegisterAction(AN("eosio.token"), ActN("transfer"), &Transfer{})
	eos.RegisterAction(AN("eosio.token"), ActN("issue"), &Issue{})
}
