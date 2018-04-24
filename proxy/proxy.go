package proxy

import (
	"github.com/eoscanada/eos-go"
)

type Route struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type PostProcessable struct {
	Route              *Route                  `json:"route"`
	P2PMessageEnvelope *eos.P2PMessageEnvelope `json:"p2p_message_envelope"`
	P2PMessage         *eos.P2PMessage         `json:"p2p_message"`
}

type PostProcessChannel chan PostProcessable
type Handler func(channel PostProcessChannel)

type PostProcessorPlugin interface {
	Start()
	Channel() PostProcessChannel
	Handler() Handler
}
