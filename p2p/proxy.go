package p2p

import (
	"fmt"

	"log"

	"time"

	"github.com/eoscanada/eos-go"
)

type Proxy struct {
	route                       *Route
	handlers                    []Handler
	waitingOriginHandShake      bool
	waitingDestinationHandShake bool
}

func NewProxy(route *Route) *Proxy {
	return &Proxy{
		route: route,
	}
}

func (p *Proxy) registerHandler(handler Handler) {
	p.handlers = append(p.handlers, handler)
}

func readMessage(connection *Connection, msgChannel chan *eos.P2PMessageEnvelope, errChannel chan error) {
	for {
		msg, err := connection.Read()
		if err != nil {
			errChannel <- fmt.Errorf("read message from %s: %s", connection.address, err)
		}
		msgChannel <- msg
	}
}

func (p *Proxy) handleOrigin(envelope *eos.P2PMessageEnvelope) error {

	destination := p.route.Destination

	switch m := envelope.P2PMessage.(type) {
	case *eos.HandshakeMessage:

		//do not forward handshake to destination to prevent a handshake storm

		fmt.Printf("Received handshake from origin: head block [%d]\n", m.HeadNum)
		destination.catchup.originHeadBlock = m.HeadNum

		if destination.handshake.HeadNum < m.HeadNum {
			fmt.Println("Destination need to catchup")
			err := destination.catchup.sendSyncRequestTo(p.route.Origin)
			if err != nil {
				return fmt.Errorf("handling origin handshake: %s", err)
			}
		}

		return nil

	case *eos.SignedBlock:

		if destination.catchup.IsCatchingUp {
			fmt.Printf("\rBlock: %d / %d", m.BlockNumber(), destination.catchup.requestedEndBlock)
		} else {
			fmt.Printf("\rBlock: %d (live)", m.BlockNumber())
		}

		err := p.handleDefault(envelope, p.route.Destination)
		if err != nil {
			return fmt.Errorf("handling origin signed block [%d]: %s", m.BlockNumber(), err)
		}

		if destination.catchup.IsCatchingUp {
			destination.catchup.headBlock = m.BlockNumber()

			if destination.catchup.headBlock == destination.catchup.requestedEndBlock {

				fmt.Println(" Sync request completed")
				destination.catchup.IsCatchingUp = false

				if destination.catchup.headBlock == destination.catchup.originHeadBlock {
					fmt.Println("Destination has catchup with origin last handshake sending handshake to origin")
					blockID, err := m.BlockID()
					if err != nil {
						return fmt.Errorf("handling origin signed block [%d]: block id: %s", m.BlockNumber(), err)
					}
					return p.route.Origin.Connection.SendHandshake(
						&HandshakeInfo{
							HeadBlockNum:  destination.catchup.headBlock,
							HeadBlockTime: m.Timestamp.Time,
							HeadBlockID:   blockID,
							//LastIrreversibleBlockID:  m.L,
							//LastIrreversibleBlockNum: m.LastIrreversibleBlockNum,
						})

				} else {
					fmt.Println("Need more block")
					err := destination.catchup.sendSyncRequestTo(p.route.Origin)
					if err != nil {
						return fmt.Errorf("handling origin signed block [%d] need more block: %s", m.BlockNumber(), err)
					}
				}
			}
		}

	default:
		return p.handleDefault(envelope, p.route.Destination)
	}
	return nil
}

func (p *Proxy) handleDestination(envelope *eos.P2PMessageEnvelope) error {
	switch m := envelope.P2PMessage.(type) {
	case *eos.HandshakeMessage:

		fmt.Printf("Received handshake from destination: head block [%d]\n", m.HeadNum)

		//Forward handshake from destination to origin. Process will resume in handleOrigin()
		p.route.Destination.handshake = *m
		p.route.Destination.catchup.headBlock = m.HeadNum
		return p.route.Origin.Connection.SendHandshake(
			&HandshakeInfo{
				HeadBlockNum:             m.HeadNum,
				HeadBlockTime:            m.Time.Time,
				HeadBlockID:              m.HeadID,
				LastIrreversibleBlockID:  m.LastIrreversibleBlockID,
				LastIrreversibleBlockNum: m.LastIrreversibleBlockNum,
			})
	default:
		return p.handleDefault(envelope, p.route.Origin)
	}
}

func (p *Proxy) handleDefault(envelope *eos.P2PMessageEnvelope, peer *Peer) error {

	switch m := envelope.P2PMessage.(type) {
	case *eos.GoAwayMessage:
		return fmt.Errorf("handling message: go away: reason [%d]", m.Reason)
	default:
		err := p.send(envelope.P2PMessage, peer)
		if err != nil {
			return fmt.Errorf("handleDefault: %s", err)
		}
	}
	return nil
}

func (p *Proxy) send(message eos.P2PMessage, peer *Peer) error {

	err := peer.Connection.Write(message)
	if err != nil {
		return fmt.Errorf("writing message to %s: %s", peer.Connection.address, err)
	}
	return nil
}

func (p *Proxy) triggerHandshake() error {
	dummyHandshakeInfo := &HandshakeInfo{
		HeadBlockID:   make([]byte, 32),
		HeadBlockNum:  0,
		HeadBlockTime: time.Now(),
	}
	fmt.Println("Sending dummy handshake to: ", p.route.Destination.Connection.address)
	// Process will resume in handleDestination()
	return p.route.Destination.Connection.SendHandshake(dummyHandshakeInfo)
}

func (p *Proxy) Start() {
	originChannel := make(chan *eos.P2PMessageEnvelope)
	destinationChannel := make(chan *eos.P2PMessageEnvelope)
	errorChannel := make(chan error)

	go readMessage(p.route.Origin.Connection, originChannel, errorChannel)
	go readMessage(p.route.Destination.Connection, destinationChannel, errorChannel)

	err := p.triggerHandshake()
	if err != nil {
		log.Fatal("proxy start: trigger handshake:", err)
	}

	for {
		select {
		case msg := <-originChannel:
			err := p.handleOrigin(msg)
			if err != nil {
				log.Fatal("proxy: handle origin:", err)
			}
		case msg := <-destinationChannel:
			err := p.handleDestination(msg)
			if err != nil {
				log.Fatal("proxy: handle destination:", err)
			}
		case err := <-errorChannel:
			log.Fatal(err)
		}
	}
}
