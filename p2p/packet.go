package p2p

import (
	"github.com/eoscanada/eos-go"
)

type Packet struct {
	Sender   *Peer
	Receiver *Peer
	Envelope *eos.P2PMessageEnvelope `json:"envelope"`
}

func NewPacket(sender *Peer, receiver *Peer, envelope *eos.P2PMessageEnvelope) *Packet {
	return &Packet{
		Sender:   sender,
		Receiver: receiver,
		Envelope: envelope,
	}
}
