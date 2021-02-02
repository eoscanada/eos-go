package p2p

import (
	"fmt"
	"math"
	"time"

	"github.com/eoscanada/eos-go"
	"go.uber.org/zap"
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
			errChannel <- fmt.Errorf("read message from %s: %w", peer.Address, err)
			break
		}

		envelope := NewEnvelope(peer, peer, packet)
		for _, handle := range c.handlers {
			handle.Handle(envelope)
		}

		switch m := packet.P2PMessage.(type) {
		case *eos.GoAwayMessage:
			errChannel <- fmt.Errorf("GoAwayMessage reason %s: %w", m.Reason, err)

		case *eos.HandshakeMessage:
			if c.catchup == nil {
				m.NodeID = peer.NodeID
				m.P2PAddress = peer.Name
				err = peer.WriteP2PMessage(m)
				if err != nil {
					errChannel <- fmt.Errorf("HandshakeMessage: %w", err)
					break
				}
				zlog.Debug("Handshake resent", zap.String("other", m.P2PAddress))

			} else {

				c.catchup.originHeadBlock = m.HeadNum
				err = c.catchup.sendSyncRequest(peer)
				if err != nil {
					errChannel <- fmt.Errorf("handshake: sending sync request: %w", err)
				}
				c.catchup.IsCatchingUp = true
			}
		case *eos.NoticeMessage:
			if c.catchup != nil {
				pendingNum := m.KnownBlocks.Pending
				if pendingNum > 0 {
					c.catchup.originHeadBlock = pendingNum
					err = c.catchup.sendSyncRequest(peer)
					if err != nil {
						errChannel <- fmt.Errorf("noticeMessage: sending sync request: %w", err)
					}
				}
			}
		case *eos.SignedBlock:

			if c.catchup != nil {

				blockNum := m.BlockNumber()
				c.catchup.headBlock = blockNum
				if c.catchup.requestedEndBlock == blockNum {

					if c.catchup.originHeadBlock <= blockNum {
						zlog.Debug("In sync with last handshake")
						blockID, err := m.BlockID()
						if err != nil {
							errChannel <- fmt.Errorf("getting block id: %w", err)
						}
						peer.handshakeInfo.HeadBlockNum = blockNum
						peer.handshakeInfo.HeadBlockID = blockID
						peer.handshakeInfo.HeadBlockTime = m.SignedBlockHeader.Timestamp.Time
						err = peer.SendHandshake(peer.handshakeInfo)
						if err != nil {
							errChannel <- fmt.Errorf("send handshake: %w", err)
						}
						zlog.Debug("Send new handshake",
							zap.Object("handshakeInfo", peer.handshakeInfo))
					} else {
						err = c.catchup.sendSyncRequest(peer)
						if err != nil {
							errChannel <- fmt.Errorf("signed block: sending sync request: %w", err)
						}
					}
				}
			}
		}
	}
}

func (c *Client) Start() error {
	zlog.Info("Starting client")

	errorChannel := make(chan error, 1)
	readyChannel := c.peer.Connect(errorChannel)

	for {
		select {
		case <-readyChannel:
			go c.read(c.peer, errorChannel)
			if c.peer.handshakeInfo != nil {

				err := triggerHandshake(c.peer)
				if err != nil {
					return fmt.Errorf("connect and start: trigger handshake: %w", err)
				}
			}
		case err := <-errorChannel:
			return fmt.Errorf("start client: %w", err)
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

	c.requestedStartBlock = c.headBlock
	c.requestedEndBlock = c.headBlock + uint32(math.Min(float64(delta), 100))

	zlog.Debug("Sending sync request",
		zap.Uint32("startBlock", c.requestedStartBlock),
		zap.Uint32("endBlock", c.requestedEndBlock))

	err := peer.SendSyncRequest(c.requestedStartBlock, c.requestedEndBlock+1)
	if err != nil {
		return fmt.Errorf("send sync request to %s: %w", peer.Address, err)
	}

	return nil
}
