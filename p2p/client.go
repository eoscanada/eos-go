package p2p

import (
	"fmt"
	"log"
	"math"

	"time"

	"github.com/eoscanada/eos-go"
)

type Client struct {
	peer        *Peer
	handlers    []Handler
	readTimeout time.Duration
	catchup     *Catchup
}

func NewClient(peer *Peer, catchup bool) *Client {
	client := &Client{
		peer: peer,
	}
	if catchup {
		client.catchup = &Catchup{
			headBlock: peer.handshakeInfo.HeadBlockNum,
		}
	}
	return client
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

		case *eos.SignedBlock:

		case *eos.HandshakeMessage:
			fmt.Println("Handshake resent!")
			if c.catchup == nil {
				m.NodeID = peer.NodeID
				m.P2PAddress = peer.Name
				err = peer.WriteP2PMessage(m)
				if err != nil {
					errChannel <- fmt.Errorf("HandshakeMessage: %s", err)
					break
				}
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

type Catchup struct {
	IsCatchingUp        bool
	requestedStartBlock uint32
	requestedEndBlock   uint32
	headBlock           uint32
	originHeadBlock     uint32
}

func (c *Catchup) sendSyncRequest(peer *Peer) error {

	c.IsCatchingUp = true

	delta := c.originHeadBlock - c.headBlock

	c.requestedStartBlock = c.headBlock + 1
	c.requestedEndBlock = c.headBlock + uint32(math.Min(float64(delta), 250))

	fmt.Printf("Sending sync request to origin: start block [%d] end block [%d]\n", c.requestedStartBlock, c.requestedEndBlock)
	err := peer.SendSyncRequest(c.requestedStartBlock, c.requestedEndBlock+1)

	if err != nil {
		return fmt.Errorf("send sync request: %s", err)
	}

	return nil

}
