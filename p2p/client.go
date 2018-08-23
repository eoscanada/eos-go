package p2p

import (
	"fmt"
	"log"
	"sync"

	"encoding/hex"

	"github.com/eoscanada/eos-go"
)

type Client struct {
	peer         *Peer
	handlers     []Handler
	handlersLock sync.Mutex
}

func NewClient(peer *Peer) *Client {
	return &Client{
		peer: peer,
	}
}

func (c *Client) RegisterHandler(handler Handler) {
	c.handlersLock.Lock()
	defer c.handlersLock.Unlock()

	c.handlers = append(c.handlers, handler)
}

func (c *Client) read(peer *Peer, errChannel chan error) {
	for {
		packet, err := peer.Read()
		if err != nil {
			errChannel <- fmt.Errorf("read message from %s: %s", peer.Address, err)
		}

		envelope := NewEnvelope(peer, peer, packet)
		c.handlersLock.Lock()
		for _, handle := range c.handlers {
			handle.Handle(envelope)
		}
		c.handlersLock.Unlock()

		switch m := packet.P2PMessage.(type) {
		case *eos.GoAwayMessage:
			log.Fatalf("handling message: go away: reason [%d]", m.Reason)

		case *eos.HandshakeMessage:
			if err != nil {
				log.Fatal(fmt.Errorf("nodeID: %s", err))
			}
			fmt.Println("Handshake resent!")
			m.P2PAddress = "localhost:5555"
			m.NodeID = make([]byte, 32)

			err = peer.WriteP2PMessage(m)
			if err != nil {
				log.Fatal(fmt.Errorf("HandshakeMessage: %s", err))
			}
		}
	}
}

func (c *Client) Start(chainID string) error {

	errorChannel := make(chan error)

	readyChannel := c.peer.Connect(errorChannel)

	for {

		select {
		case <-readyChannel:
			go c.read(c.peer, errorChannel)
			if chainID != "" {
				cID, err := hex.DecodeString(chainID)
				if err != nil {
					return fmt.Errorf("connect and start: parsing chain id: %s", err)
				}

				err = triggerHandshake(c.peer, cID)
				if err != nil {
					return fmt.Errorf("connect and start: trigger handshake: %s", err)
				}
			}
		case err := <-errorChannel:
			return err
		}
	}
}
