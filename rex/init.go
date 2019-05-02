package rex

import eos "github.com/eoscanada/eos-go"

func init() {
	eos.RegisterAction(REXAN, ActN("buyrex"), BuyREX{})
	eos.RegisterAction(REXAN, ActN("closerex"), CloseREX{})
	eos.RegisterAction(REXAN, ActN("cnclrexorder"), CancelREXorder{})
	eos.RegisterAction(REXAN, ActN("consolidate"), Consolidate{})
	eos.RegisterAction(REXAN, ActN("defcpuloan"), DefundCPULoan{})
	eos.RegisterAction(REXAN, ActN("defnetloan"), DefundNetLoan{})
	eos.RegisterAction(REXAN, ActN("deposit"), Deposit{})
	eos.RegisterAction(REXAN, ActN("fundcpuloan"), FundCPULoan{})
	eos.RegisterAction(REXAN, ActN("fundnetloan"), FundNetLoan{})
	eos.RegisterAction(REXAN, ActN("mvfrsavings"), MoveFromSavings{})
	eos.RegisterAction(REXAN, ActN("mvtosavings"), MoveToSavings{})
	eos.RegisterAction(REXAN, ActN("rentcpu"), RentCPU{})
	eos.RegisterAction(REXAN, ActN("rentnet"), RentNet{})
	eos.RegisterAction(REXAN, ActN("rexexec"), REXExec{})
	eos.RegisterAction(REXAN, ActN("sellrex"), SellREX{})
	eos.RegisterAction(REXAN, ActN("unstaketorex"), UnstakeToREX{})
	eos.RegisterAction(REXAN, ActN("updaterex"), UpdateREX{})
	eos.RegisterAction(REXAN, ActN("withdraw"), Withdraw{})
}

var AN = eos.AN
var PN = eos.PN
var ActN = eos.ActN

var REXAN = AN("eosio")
