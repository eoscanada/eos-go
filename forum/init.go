package forum

import eos "github.com/eoscanada/eos-go"

func init() {
	eos.RegisterAction(AN("eosforumtest"), ActN("post"), Post{})
	eos.RegisterAction(AN("eosforumtest"), ActN("remove"), Remove{})
	eos.RegisterAction(AN("eosforumtest"), ActN("vote"), Vote{})
}

var AN = eos.AN
var PN = eos.PN
var ActN = eos.ActN
