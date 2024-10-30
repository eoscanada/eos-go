package p2p

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"runtime"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		zlog.Error("Connection Read Err", zap.String("address", p.Address), zap.Error(err))
		return nil, fmt.Errorf("connection: read %s err: %w", p.Address, err)
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
		errChan <- fmt.Errorf("generating random node id: %w", err)
	}

	p.NodeID = nodeID
	hexNodeID := hex.EncodeToString(p.NodeID)
	p.Name = fmt.Sprintf("Client Peer - %s", hexNodeID[0:8])

	ready = make(chan bool, 1)
	go func() {
		address2log := zap.String("address", p.Address)

		if p.listener {
			zlog.Debug("Listening on", address2log)

			ln, err := net.Listen("tcp", p.Address)
			if err != nil {
				errChan <- fmt.Errorf("peer init: listening %s: %w", p.Address, err)
			}

			zlog.Debug("Accepting connection on", address2log)
			conn, err := ln.Accept()
			if err != nil {
				errChan <- fmt.Errorf("peer init: accepting connection on %s: %w", p.Address, err)
			}
			zlog.Debug("Connected on", address2log)

			p.SetConnection(conn)
			ready <- true

		} else {
			if p.handshakeTimeout > 0 {
				go func(p *Peer) {
					select {
					case <-time.After(p.handshakeTimeout):
						zlog.Warn("handshake took too long", address2log)
						errChan <- fmt.Errorf("handshake took too long: %s: %w", p.Address, err)
					case <-p.cancelHandshakeTimeout:
						zlog.Warn("cancelHandshakeTimeout canceled", address2log)
					}
				}(p)
			}

			zlog.Info("Dialing", address2log, zap.Duration("timeout", p.connectionTimeout))
			conn, err := net.DialTimeout("tcp", p.Address, p.connectionTimeout)
			if err != nil {
				if p.handshakeTimeout > 0 {
					p.cancelHandshakeTimeout <- true
				}
				errChan <- fmt.Errorf("peer init: dial %s: %w", p.Address, err)
				return
			}
			zlog.Info("Connected to", address2log)
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

	buff := bytes.NewBuffer(make([]byte, 0, 512))

	encoder := eos.NewEncoder(buff)
	err = encoder.Encode(packet)
	if err != nil {
		return fmt.Errorf("unable to encode message %s: %w", message, err)
	}

	_, err = p.Write(buff.Bytes())
	if err != nil {
		return fmt.Errorf("write msg to %s: %w", p.Address, err)
	}

	return nil
}

func (p *Peer) SendSyncRequest(startBlockNum uint32, endBlockNumber uint32) (err error) {
	zlog.Debug("SendSyncRequest",
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
	zlog.Debug("SendRequest",
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

func (p *Peer) SendNotice(headBlockNum uint32, libNum uint32, mode byte) error {
	zlog.Debug("Send Notice",
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

func (p *Peer) SendTime() error {
	zlog.Debug("SendTime", zap.String("peer", p.Address))

	notice := &eos.TimeMessage{}
	return p.WriteP2PMessage(notice)
}

func (p *Peer) SendHandshake(info *HandshakeInfo) error {
	publicKey, err := ecc.NewPublicKey("PUB_K1_1111111111111111111111111111111114T1Anm")
	if err != nil {
		return fmt.Errorf("sending handshake to %s: create public key: %w", p.Address, err)
	}

	zlog.Debug("SendHandshake", zap.String("peer", p.Address), zap.Object("info", info))

	tstamp := eos.Tstamp{Time: info.HeadBlockTime}

	data := make([]byte, 66)

	data[0] = byte(ecc.CurveK1)

	signature, _ := ecc.NewSignatureFromData(data)

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
		return fmt.Errorf("sending handshake to %s: %w", p.Address, err)
	}

	return nil
}
