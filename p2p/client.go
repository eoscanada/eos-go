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

	"math"

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
	c.api = eos.New("http://mainnet.eoscanada.com")
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

	LastHandshakeReceived *eos.HandshakeMessage
	api                   *eos.API
}

func (c *Client) ConnectRecent() error {
	return c.connect(false, 0, make([]byte, 32), time.Now(), 0, make([]byte, 32))
}

func (c *Client) ConnectAndSync(headBlock uint32, headBlockID eos.SHA256Bytes, headBlockTime time.Time, lib uint32, libID eos.SHA256Bytes) error {
	return c.connect(true, headBlock, headBlockID, headBlockTime, lib, libID)
}

func (c *Client) connect(sync bool, headBlock uint32, headBlockID eos.SHA256Bytes, headBlockTime time.Time, lib uint32, libID eos.SHA256Bytes) (err error) {

	c.registerInitHandler(sync, headBlock, headBlockID, headBlockTime, lib, libID)

	conn, err := net.Dial("tcp", c.p2pAddress)
	if err != nil {
		return err
	}

	c.Conn = conn

	println("Connecting to: ", c.p2pAddress)
	ready := make(chan bool)
	errChannel := make(chan error)
	go c.handleConnection(&Route{From: c.p2pAddress}, ready, errChannel)
	<-ready

	println("Connected")

	if err := c.SendHandshake(&HandshakeInfo{
		HeadBlockNum:             headBlock,
		LastIrreversibleBlockNum: lib,
		HeadBlockTime:            headBlockTime,
	}); err != nil {
		return err
	}

	return <-errChannel
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

var peerHeadBlock = uint32(0)
var syncHeadBlock = uint32(0)
var requestedBlock = uint32(0)
var syncing = false

func (c *Client) registerInitHandler(sync bool, headBlock uint32, headBlockID eos.SHA256Bytes, headBlockTime time.Time, lib uint32, libID eos.SHA256Bytes) {

	initHandler := HandlerFunc(func(processable Message) {

		switch msg := processable.Envelope.P2PMessage.(type) {
		case *eos.HandshakeMessage:
			c.LastHandshakeReceived = msg

			hInfo := &HandshakeInfo{
				HeadBlockNum:             msg.HeadNum,
				HeadBlockID:              msg.HeadID,
				HeadBlockTime:            msg.Time.Time,
				LastIrreversibleBlockNum: msg.LastIrreversibleBlockNum,
				LastIrreversibleBlockID:  msg.LastIrreversibleBlockID,
			}

			if sync {

				if msg.HeadNum > headBlock {
					syncHeadBlock = headBlock + 1
					peerHeadBlock = msg.HeadNum

					delta := peerHeadBlock - syncHeadBlock
					fmt.Printf("Out of sync by %d blocks \n", delta)
					requestedBlock = syncHeadBlock + uint32(math.Min(float64(delta), 250))
					fmt.Printf("Requestion block from %d to %d\n", syncHeadBlock, requestedBlock)
					syncing = true
					c.SendSyncRequest(syncHeadBlock, requestedBlock)
					return

				} else {
					fmt.Println("In sync ... Sending handshake!!!")
					hInfo = &HandshakeInfo{
						HeadBlockNum:             headBlock,
						HeadBlockID:              headBlockID,
						HeadBlockTime:            headBlockTime,
						LastIrreversibleBlockNum: lib,
						LastIrreversibleBlockID:  libID,
					}
				}
			}

			if err := c.SendHandshake(hInfo); err != nil {
				log.Println("Failed sending handshake:", err)
			}

		case *eos.SignedBlock:

			syncHeadBlock = msg.BlockNumber()

			if syncHeadBlock == requestedBlock {

				delta := peerHeadBlock - syncHeadBlock
				if delta == 0 {

					syncing = false
					sync = false
					fmt.Println("Sync completed ... Sending handshake")
					id, err := msg.BlockID()
					if err != nil {
						log.Println("blockID: ", err)
						return
					}
					hInfo := &HandshakeInfo{
						HeadBlockNum:             msg.BlockNumber(),
						HeadBlockID:              id,
						HeadBlockTime:            msg.Timestamp.Time,
						LastIrreversibleBlockNum: 0,
						LastIrreversibleBlockID:  make([]byte, 32, 32),
					}
					if err := c.SendHandshake(hInfo); err != nil {
						log.Println("Failed sending handshake:", err)
						return
					}

					fmt.Println("Send handshake: ", hInfo)

					return
				}

				requestedBlock = syncHeadBlock + uint32(math.Min(float64(delta), 250))
				syncHeadBlock++
				fmt.Println("************************************")
				fmt.Printf("Requestion more block from %d to %d\n", syncHeadBlock, requestedBlock)
				fmt.Println("************************************")
				c.SendSyncRequest(syncHeadBlock, requestedBlock)
			}
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
		fmt.Println("publicKey : ", err)
		return
	}

	tstamp := eos.Tstamp{Time: info.HeadBlockTime}

	fmt.Println("Time from fake: ", tstamp)
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
	if err != nil {
		fmt.Println("send HandshakeMessage, ", err)
	}
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

	envelope := &eos.P2PMessageEnvelope{
		Type:       message.GetType(),
		P2PMessage: message,
	}

	encoder := eos.NewEncoder(c.Conn)
	err = encoder.Encode(envelope)

	return
}

func (c *Client) handleConnection(route *Route, ready chan bool, errChannel chan error) {

	r := bufio.NewReader(c.Conn)

	ready <- true
	for {

		envelope, err := eos.ReadP2PMessageData(r)
		if err != nil {
			log.Println("Error reading from p2p client:", err)
			errChannel <- err
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
