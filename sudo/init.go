package sudo

import eos "github.com/eoscanada/eos-go"

func init() {
	eos.RegisterAction(AN("eosio.sudo"), ActN("exec"), Exec{})
}

var AN = eos.AN
var ActN = eos.ActN
