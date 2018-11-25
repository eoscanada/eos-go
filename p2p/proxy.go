package p2p

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/eoscanada/eos-go"
)

type Proxy struct {
	Peer1                       *Peer
	Peer2                       *Peer
	handlers                    []Handler
	waitingOriginHandShake      bool
	waitingDestinationHandShake bool
}

func NewProxy(peer1 *Peer, peer2 *Peer) *Proxy {
	return &Proxy{
		Peer1: peer1,
		Peer2: peer2,
	}
}

func (p *Proxy) RegisterHandler(handler Handler) {
	p.handlers = append(p.handlers, handler)
}

func (p *Proxy) RegisterHandlers(handlers []Handler) {
	p.handlers = append(p.handlers, handlers...)
}

func (p *Proxy) read(sender *Peer, receiver *Peer, errChannel chan error) {
	for {

		//p2pLog.Debug("Waiting for packet")
		packet, err := sender.Read()
		//p2pLog.Debug("Received for packet")
		if err != nil {
			errChannel <- fmt.Errorf("read message from %s: %s", sender.Address, err)
			return
		}
		err = p.handle(packet, sender, receiver)
		if err != nil {
			errChannel <- err
		}
	}
}

func (p *Proxy) handle(packet *eos.Packet, sender *Peer, receiver *Peer) error {

	_, err := receiver.Write(packet.Raw)
	if err != nil {
		return fmt.Errorf("handleDefault: %s", err)
	}

	switch m := packet.P2PMessage.(type) {
	case *eos.GoAwayMessage:
		return fmt.Errorf("handling message: go away: reason [%d]", m.Reason)
	}

	envelope := NewEnvelope(sender, receiver, packet)

	for _, handle := range p.handlers {
		handle.Handle(envelope)
	}

	return nil
}

func triggerHandshake(peer *Peer) error {
	return peer.SendHandshake(peer.handshakeInfo)
}

func (p *Proxy) ConnectAndStart() error {

	p2pLog.Info("Connecting and starting proxy")

	errorChannel := make(chan error)

	peer1ReadyChannel := p.Peer1.Connect(errorChannel)
	peer2ReadyChannel := p.Peer2.Connect(errorChannel)

	peer1Ready := false
	peer2Ready := false
	for {

		select {
		case <-peer1ReadyChannel:
			peer1Ready = true
		case <-peer2ReadyChannel:
			peer2Ready = true
		case err := <-errorChannel:
			return err
		}
		if peer1Ready && peer2Ready {
			break
		}
	}

	return p.Start()

}

func (p *Proxy) Start() error {
	p2pLog.Info("Starting readers",
		zap.String("peer1", p.Peer1.Address),
		zap.String("peer1", p.Peer2.Address))
	errorChannel := make(chan error)
	go p.read(p.Peer1, p.Peer2, errorChannel)
	go p.read(p.Peer2, p.Peer1, errorChannel)

	if p.Peer2.handshakeInfo != nil {

		err := triggerHandshake(p.Peer2)
		if err != nil {
			return fmt.Errorf("connect and start: trigger handshake: %s", err)
		}
	}

	//p2pLog.Info("Started")
	return <-errorChannel
}
