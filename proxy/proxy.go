package proxy

import "github.com/eoscanada/eos-go"

type P2PMessageChannel chan eos.P2PMessage

type PostProcessorPlugin interface {
	Start()
	Channel() P2PMessageChannel
}
