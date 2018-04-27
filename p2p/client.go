package p2p

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"encoding/hex"
	"log"

	time "time"

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

func NewClient(p2pAddr string, eosAPI *eos.API, advertiseAddress string) *Client {
	c := &Client{
		p2pAddress:          p2pAddr,
		AdvertiseP2PAddress: advertiseAddress,
		API:                 eosAPI,
	}
	copy(c.NodeID[:], []byte(advertiseAddress))
	return c
}

type Client struct {
	handlers            []Handler
	handlersLock        sync.Mutex
	p2pAddress          string
	API                 *eos.API
	AdvertiseP2PAddress string
	Conn                net.Conn
	NodeID              [32]byte
}

func (c *Client) Connect() (err error) {
	handshakeInfo, err := c.getHandshakeInfo()
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", c.p2pAddress)
	if err != nil {
		return err
	}

	c.Conn = conn

	if err := c.setupFlow(); err != nil {
		return err
	}

	fmt.Println("Connected to: ", c.p2pAddress)
	ready := make(chan bool)
	go c.handleConnection(&Route{From: c.p2pAddress}, ready)
	<-ready

	if err := c.SendHandshake(handshakeInfo); err != nil {
		return err
	}

	return nil
}

func (c *Client) RegisterHandler(h Handler) {
	c.handlersLock.Lock()
	defer c.handlersLock.Unlock()

	c.handlers = append(c.handlers, h)
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

func (c *Client) setupFlow() error {
	var hInfo handshakeInfo

	hInfo, err := c.getHandshakeInfo()
	if err != nil {
		return err
	}

	initHandler := HandlerFunc(func(processable PostProcessable) {
		msg, ok := processable.P2PMessage.(*eos.HandshakeMessage)
		if !ok {
			return
		}

		// c.SendSyncRequest(msg.LastIrreversibleBlockNum, msg.HeadNum)

		fmt.Println("Handshake time from node : ", msg.Time)

		hInfo = handshakeInfo{
			HeadBlockNum:             msg.HeadNum,
			HeadBlockID:              msg.HeadID,
			HeadBlockTime:            msg.Time.Time,
			LastIrreversibleBlockNum: msg.LastIrreversibleBlockNum,
			LastIrreversibleBlockID:  msg.LastIrreversibleBlockID,
		}
		if err := c.SendHandshake(hInfo); err != nil {
			log.Println("Failed sending handshake:", err)
		}
	})
	c.RegisterHandler(initHandler)

	return nil
}

func (c *Client) getHandshakeInfo() (info handshakeInfo, err error) {

	peerInfo, err := c.API.GetInfo()
	if err != nil {
		return
	}

	fmt.Println("Peer info: ", peerInfo)

	blockInfo, err := c.API.GetBlockByNum(uint64(peerInfo.LastIrreversibleBlockNum))
	if err != nil {
		return
	}

	info = handshakeInfo{
		HeadBlockNum:             peerInfo.HeadBlockNum,
		HeadBlockID:              decodeHex(peerInfo.HeadBlockID),
		HeadBlockTime:            peerInfo.HeadBlockTime.Time,
		LastIrreversibleBlockNum: uint32(blockInfo.BlockNum),
		LastIrreversibleBlockID:  decodeHex(blockInfo.ID),
	}

	return

}

type handshakeInfo struct {
	HeadBlockNum             uint32
	HeadBlockID              eos.SHA256Bytes
	HeadBlockTime            time.Time
	LastIrreversibleBlockNum uint32
	LastIrreversibleBlockID  eos.SHA256Bytes
}

func (c *Client) SendHandshake(info handshakeInfo) (err error) {
	pulbicKey, err := ecc.NewPublicKey("EOS1111111111111111111111111111111114T1Anm")
	if err != nil {
		return
	}
	signature, err := ecc.NewSignature("EOS111111111111111111111111111111111111111111111111111111111111111111LHpNx")
	if err != nil {
		return
	}

	tstamp := eos.Tstamp{Time: info.HeadBlockTime}

	fmt.Println("Time from fake: ", tstamp)

	handshake := &eos.HandshakeMessage{
		NetworkVersion:           int16(25431),
		ChainID:                  decodeHex("0000000000000000000000000000000000000000000000000000000000000000"),
		NodeID:                   c.NodeID[:],
		Key:                      pulbicKey,
		Time:                     tstamp,
		Token:                    decodeHex("0000000000000000000000000000000000000000000000000000000000000000"),
		Signature:                signature,
		P2PAddress:               c.AdvertiseP2PAddress,
		LastIrreversibleBlockNum: info.LastIrreversibleBlockNum,
		LastIrreversibleBlockID:  info.LastIrreversibleBlockID,
		HeadNum:                  info.HeadBlockNum,
		HeadID:                   info.HeadBlockID,
		OS:                       "linux",
		Agent:                    "Charles Billette Agent",
		Generation:               int16(1),
	}

	err = c.sendMessage(handshake)
	return
}

func (c *Client) SendSyncRequest(startBlockNum uint32, endBlockNumber uint32) (err error) {

	syncRequest := &eos.SyncRequestMessage{
		StartBlock: startBlockNum,
		EndBlock:   endBlockNumber,
	}

	c.sendMessage(syncRequest)

	return
}

func (c *Client) sendMessage(message eos.P2PMessage) (err error) {

	payload, err := eos.MarshalBinary(message)
	if err != nil {
		return
	}

	envelope := eos.P2PMessageEnvelope{
		Type:    message.GetType(),
		Payload: payload,
	}

	data, err := eos.MarshalBinary(envelope)

	var ev eos.P2PMessageEnvelope
	err = eos.UnmarshalBinary(data, &ev)
	if err != nil {
		return
	}

	n, _ := message.GetType().Name()
	fmt.Printf("Sending message [%s] to server\n", n)
	_, err = c.Conn.Write(data)
	return
}

func (c *Client) handleConnection(route *Route, ready chan bool) {

	decoder := eos.NewDecoder(bufio.NewReader(c.Conn))

	ready <- true
	for {

		var envelope eos.P2PMessageEnvelope
		err := decoder.Decode(&envelope)
		if err != nil {
			fmt.Println("Connection error: ", err)
			return
		}

		//typeName, _ := envelope.Type.Name()
		//fmt.Printf("Message received from [%s] with length: [%d] type: [%d - %s]\n", connection.RemoteAddr().String(), envelope.Length, envelope.Type, typeName)

		c.handleEnvelope(&envelope, route)

	}
}

func (c *Client) handleEnvelope(envelope *eos.P2PMessageEnvelope, route *Route) error {

	pp := PostProcessable{
		Route:              route,
		P2PMessageEnvelope: envelope,
	}

	msg, err := envelope.AsMessage()
	if err != nil {

		msgData, err := eos.MarshalBinary(envelope)
		if err != nil {
			log.Fatal(err)
		}

		return fmt.Errorf("failed for message type [%d] len[%d] with data [%s]\n", envelope.Type, envelope.Length, hex.EncodeToString(msgData))
	}

	pp.P2PMessage = msg

	c.handlersLock.Lock()
	for _, handle := range c.handlers {
		handle.Handle(pp)
	}
	c.handlersLock.Unlock()

	return nil
}
