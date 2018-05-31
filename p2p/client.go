package p2p

import (
	"bufio"
	"fmt"
	"net"
	"runtime"
	"sync"

	"encoding/hex"
	"log"

	"time"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type loggerWriter struct {
}

func (l loggerWriter) Write(p []byte) (n int, err error) {

	length := len(p)

	fmt.Printf("\t\t[%d] data [%s]\n", length, hex.EncodeToString(p))

	return length, nil
}

func NewClient(p2pAddr string, chainID eos.SHA256Bytes, networkVersion uint16) *Client {
	c := &Client{
		p2pAddress:     p2pAddr,
		ChainID:        chainID,
		NetworkVersion: networkVersion,
		AgentName:      "eos-go client",
		// by default, fake being a peer at the same level as the other..
	}
	c.registerInitHandler()
	c.NodeID = chainID
	return c
}

type Client struct {
	handlers       []Handler
	handlersLock   sync.Mutex
	p2pAddress     string
	ChainID        eos.SHA256Bytes
	NetworkVersion uint16
	Conn           net.Conn
	NodeID         eos.SHA256Bytes
	SigningKey     *ecc.PrivateKey
	AgentName      string
	ShowEmptyChain bool

	LastHandshakeReceived *eos.HandshakeMessage
}

func (c *Client) Connect() (err error) {
	conn, err := net.Dial("tcp", c.p2pAddress)
	if err != nil {
		return err
	}

	c.Conn = conn

	println("Connecting to: ", c.p2pAddress)
	ready := make(chan bool)
	go c.handleConnection(&Route{From: c.p2pAddress}, ready)
	<-ready

	println("Connected")

	if err := c.SendHandshake(&HandshakeInfo{
		HeadBlockNum:             0,
		LastIrreversibleBlockNum: 0,
		HeadBlockTime:            time.Now(),
	}); err != nil {
		return err
	}

	return nil
}

func (c *Client) RegisterHandler(h Handler) {
	c.handlersLock.Lock()
	defer c.handlersLock.Unlock()

	c.handlers = append(c.handlers, h)
}

func (c *Client) RegisterHandlerFunc(f func(Message)) Handler {
	h := HandlerFunc(f)
	c.RegisterHandler(h)
	return h
}

func (c *Client) UnregisterHandler(h Handler) {
	c.handlersLock.Lock()
	defer c.handlersLock.Unlock()

	var newHandlers []Handler
	for _, handler := range c.handlers {
		if handler != h {
			newHandlers = append(newHandlers, handler)
		}
	}
	c.handlers = newHandlers
}

func (c *Client) registerInitHandler() {
	initHandler := HandlerFunc(func(processable Message) {
		msg, ok := processable.Envelope.P2PMessage.(*eos.HandshakeMessage)
		if !ok {
			return
		}

		// Keep track of last handshake received..
		c.LastHandshakeReceived = msg

		var hInfo HandshakeInfo
		if c.ShowEmptyChain {
			hInfo = HandshakeInfo{
				HeadBlockNum:             0,
				HeadBlockID:              make([]byte, 32, 32),
				HeadBlockTime:            msg.Time.Time,
				LastIrreversibleBlockNum: 0,
				LastIrreversibleBlockID:  make([]byte, 32, 32),
			}
		} else {
			hInfo = HandshakeInfo{
				HeadBlockNum:             msg.HeadNum,
				HeadBlockID:              msg.HeadID,
				HeadBlockTime:            msg.Time.Time,
				LastIrreversibleBlockNum: msg.LastIrreversibleBlockNum,
				LastIrreversibleBlockID:  msg.LastIrreversibleBlockID,
			}
		}
		if err := c.SendHandshake(&hInfo); err != nil {
			log.Println("Failed sending handshake:", err)
		}

	})
	c.RegisterHandler(initHandler)
}

type HandshakeInfo struct {
	HeadBlockNum             uint32
	HeadBlockID              eos.SHA256Bytes
	HeadBlockTime            time.Time
	LastIrreversibleBlockNum uint32
	LastIrreversibleBlockID  eos.SHA256Bytes
}

func (c *Client) SendHandshake(info *HandshakeInfo) (err error) {
	publicKey, err := ecc.NewPublicKey("EOS1111111111111111111111111111111114T1Anm")
	if err != nil {
		return
	}

	tstamp := eos.Tstamp{Time: info.HeadBlockTime}

	println("Time from fake: ", tstamp)
	//tData, err := eos.MarshalBinary(&tstamp)
	//if err != nil {
	//	return fmt.Errorf("marshalling tstamp, %s", err)
	//}
	//h := ripemd160.New()
	//_, err = h.Write(tData)
	//if err != nil {
	//	return fmt.Errorf("hashing tstamp data, %s", err)
	//}

	//time := fmt.Sprintf("%d", tstamp.Unix())
	//token := sha256.Sum256([]byte("1526431521355589"))

	//c.SigningKey.Curve = ecc.CurveR1
	// signature, err := c.SigningKey.Sign(token[:])
	// fmt.Println("signature: ", signature)
	// if err != nil {
	// 	return fmt.Errorf("signing token data, %s", err)
	// }
	signature := ecc.Signature{
		Curve:   ecc.CurveK1,
		Content: make([]byte, 65, 65),
	}

	handshake := &eos.HandshakeMessage{
		NetworkVersion:           c.NetworkVersion,
		ChainID:                  c.ChainID,
		NodeID:                   c.NodeID,
		Key:                      publicKey,
		Time:                     tstamp,
		Token:                    make([]byte, 32, 32), // token[:]
		Signature:                signature,
		P2PAddress:               c.p2pAddress,
		LastIrreversibleBlockNum: info.LastIrreversibleBlockNum,
		LastIrreversibleBlockID:  info.LastIrreversibleBlockID,
		HeadNum:                  info.HeadBlockNum,
		HeadID:                   info.HeadBlockID,
		OS:                       runtime.GOOS,
		Agent:                    c.AgentName,
		Generation:               int16(1),
	}

	err = c.sendMessage(handshake)
	return
}

func (c *Client) SendSyncRequest(startBlockNum uint32, endBlockNumber uint32) (err error) {
	println("SendSyncRequest start [%d] end [%d]\n", startBlockNum, endBlockNumber)
	syncRequest := &eos.SyncRequestMessage{
		StartBlock: startBlockNum,
		EndBlock:   endBlockNumber,
	}

	return c.sendMessage(syncRequest)
}

func (c *Client) sendMessage(message eos.P2PMessage) (err error) {
	n, _ := message.GetType().Name()
	println("Sending message [%s] to server\n", n)

	envelope := &eos.P2PMessageEnvelope{
		Type:       message.GetType(),
		P2PMessage: message,
	}

	encoder := eos.NewEncoder(c.Conn)
	err = encoder.Encode(envelope)

	return
}

func (c *Client) handleConnection(route *Route, ready chan bool) {

	r := bufio.NewReader(c.Conn)

	ready <- true
	for {

		envelope, err := eos.ReadP2PMessageData(r)
		if err != nil {
			log.Println("Error reading from p2p client:", err)
			// TODO: kill the socket, do something !
			return
		}

		pp := Message{
			Route:    route,
			Envelope: envelope,
		}

		c.handlersLock.Lock()
		for _, handle := range c.handlers {
			handle.Handle(pp)
		}
		c.handlersLock.Unlock()

	}
}
