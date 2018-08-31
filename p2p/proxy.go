package p2p

import (
	"fmt"

	"time"

	"log"

	"encoding/hex"

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

		log.Println("Waiting for packet")
		packet, err := sender.Read()
		log.Println("Received for packet")
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

func triggerHandshake(peer *Peer, chainID eos.SHA256Bytes) error {
	log.Println("Sending dummy handshake to: ", peer.Address)
	dummyHandshakeInfo := &HandshakeInfo{
		ChainID:       chainID,
		HeadBlockID:   make([]byte, 32),
		HeadBlockNum:  0,
		HeadBlockTime: time.Now(),
	}
	// Process will resume in handle()

	return peer.SendHandshake(dummyHandshakeInfo)
}

func (p *Proxy) ConnectAndStart(chainID string) error {

	log.Println("Connecting and starting proxy with chain id:", chainID)

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

	return p.Start(chainID)

}

func (p *Proxy) Start(chainID string) error {

	log.Println("Starting readers")
	errorChannel := make(chan error)
	go p.read(p.Peer1, p.Peer2, errorChannel)
	go p.read(p.Peer2, p.Peer1, errorChannel)

	if chainID != "" {
		cID, err := hex.DecodeString(chainID)
		if err != nil {
			return fmt.Errorf("connect and start: parsing chain id: %s", err)
		}

		err = triggerHandshake(p.Peer2, cID)
		if err != nil {
			return fmt.Errorf("connect and start: trigger handshake: %s", err)
		}
	}

	log.Println("Started")
	return <-errorChannel
}
