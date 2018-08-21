package p2p

import (
	"fmt"
	"time"

	"math"

	"github.com/eoscanada/eos-go"
)

type Route struct {
	Origin      *Peer
	Destination *Peer
}

type Peer struct {
	Connection *Connection
	handshake  eos.HandshakeMessage
	catchup    Catchup
}

type HandshakeInfo struct {
	HeadBlockNum             uint32
	HeadBlockID              eos.SHA256Bytes
	HeadBlockTime            time.Time
	LastIrreversibleBlockNum uint32
	LastIrreversibleBlockID  eos.SHA256Bytes
}

func NewPeer(connection *Connection) *Peer {

	return &Peer{
		Connection: connection,
	}

}

type Catchup struct {
	IsCatchingUp        bool
	requestedStartBlock uint32
	requestedEndBlock   uint32
	headBlock           uint32
	originHeadBlock     uint32
}

func (c *Catchup) sendSyncRequestTo(peer *Peer) error {

	c.IsCatchingUp = true

	delta := c.originHeadBlock - c.headBlock

	c.requestedStartBlock = c.headBlock + 1
	c.requestedEndBlock = c.headBlock + uint32(math.Min(float64(delta), 250))

	fmt.Printf("Sending sync request to origin: start block [%d] end block [%d]\n", c.requestedStartBlock, c.requestedEndBlock)
	err := peer.Connection.SendSyncRequest(c.requestedStartBlock, c.requestedEndBlock+1)

	if err != nil {
		return fmt.Errorf("send sync request: %s", err)
	}

	return nil

}
