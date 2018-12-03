package p2p

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"time"

	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"

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

// MarshalLogObject calls the underlying function from zap.
func (p Peer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", p.Name)
	enc.AddString("address", p.Address)
	enc.AddString("agent", p.agent)
	return enc.AddObject("handshakeInfo", p.handshakeInfo)
}

type HandshakeInfo struct {
	ChainID                  eos.Checksum256
	HeadBlockNum             uint32
	HeadBlockID              eos.Checksum256
	HeadBlockTime            time.Time
	LastIrreversibleBlockNum uint32
	LastIrreversibleBlockID  eos.Checksum256
}

func (h *HandshakeInfo) String() string {
	return fmt.Sprintf("Handshake Info: HeadBlockNum [%d], LastIrreversibleBlockNum [%d]", h.HeadBlockNum, h.LastIrreversibleBlockNum)
}

// MarshalLogObject calls the underlying function from zap.
func (h HandshakeInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("chainID", h.ChainID.String())
	enc.AddUint32("headBlockNum", h.HeadBlockNum)
	enc.AddString("headBlockID", h.HeadBlockID.String())
	enc.AddTime("headBlockTime", h.HeadBlockTime)
	enc.AddUint32("lastIrreversibleBlockNum", h.LastIrreversibleBlockNum)
	enc.AddString("lastIrreversibleBlockID", h.LastIrreversibleBlockID.String())
	return nil
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
		p2pLog.Error("Connection Read Err", zap.String("address", p.Address), zap.Error(err))
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
		address2log := zap.String("address", p.Address)

		if p.listener {
			p2pLog.Debug("Listening on", address2log)

			ln, err := net.Listen("tcp", p.Address)
			if err != nil {
				errChan <- fmt.Errorf("peer init: listening %s: %s", p.Address, err)
			}

			p2pLog.Debug("Accepting connection on", address2log)
			conn, err := ln.Accept()
			if err != nil {
				errChan <- fmt.Errorf("peer init: accepting connection on %s: %s", p.Address, err)
			}
			p2pLog.Debug("Connected on", address2log)

			p.SetConnection(conn)
			ready <- true

		} else {
			if p.handshakeTimeout > 0 {
				go func(p *Peer) {
					select {
					case <-time.After(p.handshakeTimeout):
						p2pLog.Warn("handshake took too long", address2log)
						errChan <- fmt.Errorf("handshake took too long: %s", p.Address)
					case <-p.cancelHandshakeTimeout:
						p2pLog.Warn("cancelHandshakeTimeout canceled", address2log)
					}
				}(p)
			}

			p2pLog.Info("Dialing", address2log, zap.Duration("timeout", p.connectionTimeout))
			conn, err := net.DialTimeout("tcp", p.Address, p.connectionTimeout)
			if err != nil {
				if p.handshakeTimeout > 0 {
					p.cancelHandshakeTimeout <- true
				}
				errChan <- fmt.Errorf("peer init: dial %s: %s", p.Address, err)
				return
			}
			p2pLog.Info("Connected to", address2log)
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
	p2pLog.Debug("SendSyncRequest",
		zap.String("peer", p.Address),
		zap.Uint32("start", startBlockNum),
		zap.Uint32("end", endBlockNumber))

	syncRequest := &eos.SyncRequestMessage{
		StartBlock: startBlockNum,
		EndBlock:   endBlockNumber,
	}

	return p.WriteP2PMessage(syncRequest)
}
func (p *Peer) SendRequest(startBlockNum uint32, endBlockNumber uint32) (err error) {
	p2pLog.Debug("SendRequest",
		zap.String("peer", p.Address),
		zap.Uint32("start", startBlockNum),
		zap.Uint32("end", endBlockNumber))

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
	p2pLog.Debug("Send Notice",
		zap.String("peer", p.Address),
		zap.Uint32("head", headBlockNum),
		zap.Uint32("lib", libNum),
		zap.Uint8("type", mode))

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
	p2pLog.Debug("SendTime", zap.String("peer", p.Address))

	notice := &eos.TimeMessage{}
	return p.WriteP2PMessage(notice)
}

func (p *Peer) SendHandshake(info *HandshakeInfo) (err error) {

	publicKey, err := ecc.NewPublicKey("EOS1111111111111111111111111111111114T1Anm")
	if err != nil {
		logErr("publicKey err", err)
		err = fmt.Errorf("sending handshake to %s: create public key: %s", p.Address, err)
		return
	}

	p2pLog.Debug("SendHandshake", zap.String("peer", p.Address), zap.Object("info", info))

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
