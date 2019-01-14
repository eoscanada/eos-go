package forum

import eos "github.com/eoscanada/eos-go"

func init() {
	eos.RegisterAction(ForumAN, ActN("clnproposal"), CleanProposal{})
	eos.RegisterAction(ForumAN, ActN("expire"), Expire{})
	eos.RegisterAction(ForumAN, ActN("post"), Post{})
	eos.RegisterAction(ForumAN, ActN("propose"), Propose{})
	eos.RegisterAction(ForumAN, ActN("status"), Status{})
	eos.RegisterAction(ForumAN, ActN("unpost"), UnPost{})
	eos.RegisterAction(ForumAN, ActN("unvote"), UnVote{})
	eos.RegisterAction(ForumAN, ActN("vote"), Vote{})
}

var AN = eos.AN
var PN = eos.PN
var ActN = eos.ActN

var ForumAN = AN("eosio.forum")
