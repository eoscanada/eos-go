package p2p

import (
	"bufio"
	"fmt"
	"net"

	"encoding/hex"
	"log"

	"bytes"
	"net/url"

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

type Client struct {
	Handlers     []Handler
	flowHandlers map[eos.P2PMessageType]Handler
	Api          *eos.API
}

func decodeHex(hexString string) (data []byte) {

	data, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println("decodeHex error: ", err)
	}
	return data
}
func (c *Client) Dial(p2pAddress string, webserviceAddress string) (err error) {

	c.Api = eos.New(&url.URL{Scheme: "http", Host: webserviceAddress}, bytes.Repeat([]byte{0}, 32))
	c.flowHandlers = make(map[eos.P2PMessageType]Handler)

	handshakeInfo, err := c.getHandshakeInfo()
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", p2pAddress)
	if err != nil {
		return err
	}

	c.setupFlow(conn)

	fmt.Println("Connected to: ", p2pAddress)
	ready := make(chan bool)
	go c.handleConnection(conn, &Route{From: p2pAddress}, ready)
	<-ready

	c.SendHandshake(handshakeInfo, conn)

	return
}

func (c *Client) setupFlow(conn net.Conn) {
	var hInfo handshakeInfo
	c.flowHandlers[eos.HandshakeMessageType] = func(processable PostProcessable) {

		handshakeMessage := processable.P2PMessage.(*eos.HandshakeMessage)

		c.SendSyncRequest(handshakeMessage.LastIrreversibleBlockNum, handshakeMessage.HeadNum, conn)

		fmt.Println("Handshake time from node : ", handshakeMessage.Time)

		hInfo = handshakeInfo{
			HeadBlockNum:             handshakeMessage.HeadNum,
			HeadBlockID:              handshakeMessage.HeadID,
			HeadBlockTime:            handshakeMessage.Time.Time,
			LastIrreversibleBlockNum: handshakeMessage.LastIrreversibleBlockNum,
			LastIrreversibleBlockID:  handshakeMessage.LastIrreversibleBlockID,
		}
	}

	count := 0
	c.flowHandlers[eos.SignedBlockMessageType] = func(processable PostProcessable) {

		count = count + 1

		if count >= 2 {
			c.SendHandshake(hInfo, conn)
		}
	}
}

func (c *Client) getHandshakeInfo() (info handshakeInfo, err error) {

	peerInfo, err := c.Api.GetInfo()
	if err != nil {
		return
	}

	fmt.Println("Peer info: ", peerInfo)

	blockInfo, err := c.Api.GetBlockByNum(uint64(peerInfo.LastIrreversibleBlockNum))
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

func (c *Client) SendHandshake(info handshakeInfo, toConnection net.Conn) (err error) {

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
		NodeID:                   decodeHex("b79243d6facfb19de89dd50405dd7958cf17afebedb10203b86442348b14c7a5"),
		Key:                      pulbicKey,
		Time:                     tstamp,
		Token:                    decodeHex("0000000000000000000000000000000000000000000000000000000000000000"),
		Signature:                signature,
		P2PAddress:               "Charless-MacBook-Pro-2.local:29876 - 1234567",
		LastIrreversibleBlockNum: info.LastIrreversibleBlockNum,
		LastIrreversibleBlockID:  info.LastIrreversibleBlockID,
		HeadNum:                  info.HeadBlockNum,
		HeadID:                   info.HeadBlockID,
		OS:                       "linux",
		Agent:                    "Charles Billette Agent",
		Generation:               int16(1),
	}

	err = c.sendMessage(handshake, toConnection)
	return
}

func (c *Client) SendSyncRequest(startBlockNum uint32, endBlockNumber uint32, toConnection net.Conn) (err error) {

	syncRequest := &eos.SyncRequestMessage{
		StartBlock: startBlockNum,
		EndBlock:   endBlockNumber,
	}

	c.sendMessage(syncRequest, toConnection)

	return
}

func (c *Client) sendMessage(message eos.P2PMessage, conn net.Conn) (err error) {

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
	_, err = conn.Write(data)
	return
}

func (c *Client) handleConnection(connection net.Conn, route *Route, ready chan bool) {

	decoder := eos.NewDecoder(bufio.NewReader(connection))

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

		c.handleEnvelop(&envelope, route)

	}
}

func (c *Client) handleEnvelop(envelope *eos.P2PMessageEnvelope, route *Route) error {

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

	for _, handle := range c.Handlers {
		handle(pp)
	}

	if handle, ok := c.flowHandlers[msg.GetType()]; ok {
		handle(pp)
	}

	return nil
}
