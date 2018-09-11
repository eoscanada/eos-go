package p2p

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"math"

	"runtime"

	"bufio"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type Peer struct {
	Address                string
	agent                  string
	connection             net.Conn
	reader                 io.Reader
	handshake              eos.HandshakeMessage
	catchup                Catchup
	listener               bool
	mockHandshake          bool
	connectionTimeout      time.Duration
	handshakeTimeout       time.Duration
	cancelHandshakeTimeout chan bool
}

type HandshakeInfo struct {
	ChainID                  eos.SHA256Bytes
	HeadBlockNum             uint32
	HeadBlockID              eos.SHA256Bytes
	HeadBlockTime            time.Time
	LastIrreversibleBlockNum uint32
	LastIrreversibleBlockID  eos.SHA256Bytes
}

func (p *Peer) SetHandshakeTimeout(timeout time.Duration) {
	p.handshakeTimeout = timeout
}

func (p *Peer) SetConnectionTimeout(timeout time.Duration) {
	p.connectionTimeout = timeout
}

func newPeer(address string, agent string, listener bool) *Peer {

	return &Peer{
		Address:                address,
		agent:                  agent,
		listener:               listener,
		cancelHandshakeTimeout: make(chan bool),
	}
}

func NewIncommingPeer(address string, agent string) *Peer {
	return newPeer(address, agent, true)
}

func NewOutgoingPeer(address string, agent string) *Peer {
	return newPeer(address, agent, false)
}

func (p *Peer) Read() (*eos.Packet, error) {
	packet, err := eos.ReadPacket(p.reader)
	if err != nil {
		log.Println("Connection Read error:", p.Address, err)
		return nil, fmt.Errorf("connection: read: %s", err)
	}
	if p.handshakeTimeout > 0 {
		p.cancelHandshakeTimeout <- true
	}
	return packet, nil
}

func (p *Peer) SetConnection(conn net.Conn) {
	p.connection = conn
	p.reader = bufio.NewReader(p.connection)
}

func (p *Peer) Connect(errChan chan error) (ready chan bool) {

	ready = make(chan bool, 1)
	go func() {
		if p.listener {
			fmt.Println("Listening on:", p.Address)

			ln, err := net.Listen("tcp", p.Address)
			if err != nil {
				errChan <- fmt.Errorf("peer init: listening %s: %s", p.Address, err)
			}

			fmt.Println("Accepting connection on:\n", p.Address)
			conn, err := ln.Accept()
			if err != nil {
				errChan <- fmt.Errorf("peer init: accepting connection on %s: %s", p.Address, err)
			}
			fmt.Println("Connected on:", p.Address)

			p.SetConnection(conn)
			ready <- true

		} else {
			if p.handshakeTimeout > 0 {
				go func(p *Peer) {
					select {
					case <-time.After(p.handshakeTimeout):
						log.Println("Handshake took too long:", p.Address)
						errChan <- fmt.Errorf("handshake took too long: %s", p.Address)
					case <-p.cancelHandshakeTimeout:
						log.Println("cancelHandshakeTimeout canceled:", p.Address)
					}
				}(p)
			}

			log.Printf("Dialing: %s, timeout: %d\n", p.Address, p.connectionTimeout)
			conn, err := net.DialTimeout("tcp", p.Address, p.connectionTimeout)
			if err != nil {
				errChan <- fmt.Errorf("peer init: dial %s: %s", p.Address, err)
				return
			}
			log.Println("Connected to:", p.Address)
			p.connection = conn
			p.reader = bufio.NewReader(conn)
			ready <- true
		}
	}()

	return
}

func (p *Peer) Write(bytes []byte) (int, error) {

	return p.connection.Write(bytes)
}

func (p *Peer) WriteP2PMessage(message eos.P2PMessage) (err error) {

	packet := &eos.Packet{
		Type:       message.GetType(),
		P2PMessage: message,
	}

	encoder := eos.NewEncoder(p.connection)
	err = encoder.Encode(packet)

	return
}

func (p *Peer) SendSyncRequest(startBlockNum uint32, endBlockNumber uint32) (err error) {
	println("SendSyncRequest start [%d] end [%d]\n", startBlockNum, endBlockNumber)
	syncRequest := &eos.SyncRequestMessage{
		StartBlock: startBlockNum,
		EndBlock:   endBlockNumber,
	}

	return p.WriteP2PMessage(syncRequest)
}
func (p *Peer) SendRequest(startBlockNum uint32, endBlockNumber uint32) (err error) {
	fmt.Printf("SendRequest start [%d] end [%d]\n", startBlockNum, endBlockNumber)
	request := &eos.RequestMessage{
		ReqTrx: eos.OrderedBlockIDs{
			//Unknown: [3]byte{},
			Mode:    [4]byte{0, 0, 0, 2},
			Pending: startBlockNum,
		},
		ReqBlocks: eos.OrderedBlockIDs{
			//Unknown: [3]byte{},
			Mode:    [4]byte{0, 0, 0, 2},
			Pending: endBlockNumber,
		},
	}

	return p.WriteP2PMessage(request)
}

func (p *Peer) SendNotice(headBlockNum uint32, libNum uint32, mode byte) (err error) {
	fmt.Printf("Send Notice head [%d] lib [%d] type[%d]\n", headBlockNum, libNum, mode)

	notice := &eos.NoticeMessage{
		KnownTrx: eos.OrderedBlockIDs{
			Mode:    [4]byte{mode, 0, 0, 0},
			Pending: headBlockNum,
		},
		KnownBlocks: eos.OrderedBlockIDs{
			Mode:    [4]byte{mode, 0, 0, 0},
			Pending: libNum,
		},
	}
	return p.WriteP2PMessage(notice)
}

func (p *Peer) SendTime() (err error) {
	fmt.Printf("SendTime\n")

	notice := &eos.TimeMessage{}
	return p.WriteP2PMessage(notice)
}

func (p *Peer) SendHandshake(info *HandshakeInfo) (err error) {

	publicKey, err := ecc.NewPublicKey("EOS1111111111111111111111111111111114T1Anm")
	if err != nil {
		fmt.Println("publicKey : ", err)
		err = fmt.Errorf("sending handshake to %s: create public key: %s", p.Address, err)
		return
	}

	tstamp := eos.Tstamp{Time: info.HeadBlockTime}

	signature := ecc.Signature{
		Curve:   ecc.CurveK1,
		Content: make([]byte, 65, 65),
	}

	handshake := &eos.HandshakeMessage{
		NetworkVersion:           1206,
		ChainID:                  info.ChainID,
		NodeID:                   make([]byte, 32),
		Key:                      publicKey,
		Time:                     tstamp,
		Token:                    make([]byte, 32, 32), // token[:]
		Signature:                signature,
		P2PAddress:               p.Address,
		LastIrreversibleBlockNum: info.LastIrreversibleBlockNum,
		LastIrreversibleBlockID:  info.LastIrreversibleBlockID,
		HeadNum:                  info.HeadBlockNum,
		HeadID:                   info.HeadBlockID,
		OS:                       runtime.GOOS,
		Agent:                    p.agent,
		Generation:               int16(1),
	}

	err = p.WriteP2PMessage(handshake)
	if err != nil {
		err = fmt.Errorf("sending handshake to %s: %s", p.Address, err)
	}
	return
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
