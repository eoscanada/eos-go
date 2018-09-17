package p2p

import (
	"fmt"
	"log"

	"time"

	"github.com/eoscanada/eos-go"
)

type Client struct {
	peer        *Peer
	handlers    []Handler
	readTimeout time.Duration
}

func NewClient(peer *Peer) *Client {
	return &Client{
		peer: peer,
	}
}

func (c *Client) CloseConnection() error {
	if c.peer.connection == nil {
		return nil
	}
	return c.peer.connection.Close()
}

func (c *Client) SetReadTimeout(readTimeout time.Duration) {
	c.readTimeout = readTimeout
}

func (c *Client) RegisterHandler(handler Handler) {

	c.handlers = append(c.handlers, handler)
}

func (c *Client) read(peer *Peer, errChannel chan error) {
	for {
		packet, err := peer.Read()
		if err != nil {
			errChannel <- fmt.Errorf("read message from %s: %s", peer.Address, err)
			break
		}

		envelope := NewEnvelope(peer, peer, packet)
		for _, handle := range c.handlers {
			handle.Handle(envelope)
		}

		switch m := packet.P2PMessage.(type) {
		case *eos.GoAwayMessage:
			errChannel <- fmt.Errorf("GoAwayMessage reason [%s]: %s", m.Reason, err)

		case *eos.HandshakeMessage:
			fmt.Println("Handshake resent!")
			m.P2PAddress = "localhost:5555"
			m.NodeID = make([]byte, 32)

			err = peer.WriteP2PMessage(m)
			if err != nil {
				errChannel <- fmt.Errorf("HandshakeMessage: %s", err)
				break
			}
		}
	}
}

func (c *Client) Start() error {

	fmt.Println("Starting client")

	errorChannel := make(chan error, 1)

	readyChannel := c.peer.Connect(errorChannel)

	for {
		select {
		case <-readyChannel:
			go c.read(c.peer, errorChannel)
			if c.peer.handshakeInfo != nil {

				err := triggerHandshake(c.peer)
				if err != nil {
					return fmt.Errorf("connect and start: trigger handshake: %s", err)
				}
			}
		case err := <-errorChannel:
			log.Println("Start got ERROR:", err)
			return err
		}
	}
}
