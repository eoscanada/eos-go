package p2p

import (
	"bufio"
	"net"

	"fmt"

	"io"

	"runtime"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type Connection struct {
	address        string
	chainID        eos.SHA256Bytes
	agent          string
	nodeConnection net.Conn
	reader         io.Reader
}

func NewConnection(address string, chainID eos.SHA256Bytes, agent string, conn net.Conn) *Connection {
	connection := &Connection{
		address:        address,
		chainID:        chainID,
		agent:          agent,
		nodeConnection: conn,
		reader:         bufio.NewReader(conn),
	}
	return connection
}

func NewOutgoingConnection(address string, chainID eos.SHA256Bytes, agent string) (*Connection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("NewOutgoingConnection: dial %s : %s", address, err)
	}
	return NewConnection(address, chainID, agent, conn), nil
}

func (c *Connection) Read() (*eos.P2PMessageEnvelope, error) {
	envelope, err := eos.ReadP2PMessageData(c.reader)
	if err != nil {
		return nil, fmt.Errorf("connection: read: %s", err)
	}
	return envelope, nil
}

func (c *Connection) SendSyncRequest(startBlockNum uint32, endBlockNumber uint32) (err error) {
	println("SendSyncRequest start [%d] end [%d]\n", startBlockNum, endBlockNumber)
	syncRequest := &eos.SyncRequestMessage{
		StartBlock: startBlockNum,
		EndBlock:   endBlockNumber,
	}

	return c.Write(syncRequest)
}

func (c *Connection) Write(message eos.P2PMessage) (err error) {

	envelope := &eos.P2PMessageEnvelope{
		Type:       message.GetType(),
		P2PMessage: message,
	}

	encoder := eos.NewEncoder(c.nodeConnection)
	err = encoder.Encode(envelope)

	return
}

func (c *Connection) SendHandshake(info *HandshakeInfo) (err error) {

	publicKey, err := ecc.NewPublicKey("EOS1111111111111111111111111111111114T1Anm")
	if err != nil {
		fmt.Println("publicKey : ", err)
		err = fmt.Errorf("sending handshake to %s: create public key: %s", c.address, err)
		return
	}

	tstamp := eos.Tstamp{Time: info.HeadBlockTime}

	signature := ecc.Signature{
		Curve:   ecc.CurveK1,
		Content: make([]byte, 65, 65),
	}

	handshake := &eos.HandshakeMessage{
		NetworkVersion:           1206,
		ChainID:                  c.chainID,
		NodeID:                   make([]byte, 32),
		Key:                      publicKey,
		Time:                     tstamp,
		Token:                    make([]byte, 32, 32), // token[:]
		Signature:                signature,
		P2PAddress:               c.address,
		LastIrreversibleBlockNum: info.LastIrreversibleBlockNum,
		LastIrreversibleBlockID:  info.LastIrreversibleBlockID,
		HeadNum:                  info.HeadBlockNum,
		HeadID:                   info.HeadBlockID,
		OS:                       runtime.GOOS,
		Agent:                    c.agent,
		Generation:               int16(1),
	}

	err = c.Write(handshake)
	if err != nil {
		err = fmt.Errorf("sending handshake to %s: %s", c.address, err)
	}
	return
}
