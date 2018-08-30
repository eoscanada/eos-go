package p2p

import (
	"github.com/eoscanada/eos-go"
)

type Envelope struct {
	Sender   *Peer
	Receiver *Peer
	Packet   *eos.Packet `json:"envelope"`
}

func NewEnvelope(sender *Peer, receiver *Peer, packet *eos.Packet) *Envelope {
	return &Envelope{
		Sender:   sender,
		Receiver: receiver,
		Packet:   packet,
	}
}
