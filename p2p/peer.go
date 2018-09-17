package p2p

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"runtime"

	"bufio"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type Peer struct {
	Address                string
	Name                   string
	agent                  string
	NodeID                 []byte
	connection             net.Conn
	reader                 io.Reader
	listener               bool
	handshakeInfo          *HandshakeInfo
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

func (h *HandshakeInfo) String() string {
	return fmt.Sprintf("Handshake Info: HeadBlockNum [%d], LastIrreversibleBlockNum [%d] ConnectingPeerAddress [%s]", h.HeadBlockNum, h.LastIrreversibleBlockNum)
}

func (p *Peer) SetHandshakeTimeout(timeout time.Duration) {
	p.handshakeTimeout = timeout
}

func (p *Peer) SetConnectionTimeout(timeout time.Duration) {
	p.connectionTimeout = timeout
}

func newPeer(address string, agent string, listener bool, handshakeInfo *HandshakeInfo) *Peer {

	return &Peer{
		Address:                address,
		agent:                  agent,
		listener:               listener,
		handshakeInfo:          handshakeInfo,
		cancelHandshakeTimeout: make(chan bool),
	}
}

func NewIncommingPeer(address string, agent string) *Peer {
	return newPeer(address, agent, true, nil)
}

func NewOutgoingPeer(address string, agent string, handshakeInfo *HandshakeInfo) *Peer {
	return newPeer(address, agent, false, handshakeInfo)
}

func (p *Peer) Read() (*eos.Packet, error) {
	packet, err := eos.ReadPacket(p.reader)
	if p.handshakeTimeout > 0 {
		p.cancelHandshakeTimeout <- true
	}
	if err != nil {
		log.Println("Connection Read error:", p.Address, err)
		return nil, fmt.Errorf("connection: read: %s", err)
	}
	return packet, nil
}

func (p *Peer) SetConnection(conn net.Conn) {
	p.connection = conn
	p.reader = bufio.NewReader(p.connection)
}

func (p *Peer) Connect(errChan chan error) (ready chan bool) {

	nodeID := make([]byte, 32)
	_, err := rand.Read(nodeID)
	if err != nil {
		errChan <- fmt.Errorf("generating random node id: %s", err)
	}

	p.NodeID = nodeID
	hexNodeID := hex.EncodeToString(p.NodeID)
	p.Name = fmt.Sprintf("Client Peer - %s", hexNodeID[0:8])

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
				if p.handshakeTimeout > 0 {
					p.cancelHandshakeTimeout <- true
				}
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
			Mode:    [4]byte{0, 0, 0, 0},
			Pending: startBlockNum,
		},
		ReqBlocks: eos.OrderedBlockIDs{
			Mode:    [4]byte{0, 0, 0, 0},
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
		NodeID:                   p.NodeID,
		Key:                      publicKey,
		Time:                     tstamp,
		Token:                    make([]byte, 32, 32),
		Signature:                signature,
		P2PAddress:               p.Name,
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
