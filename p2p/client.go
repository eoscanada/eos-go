package p2p

import (
	"bufio"
	"fmt"
	"net"

	"encoding/hex"
	"log"

	"github.com/eoscanada/eos-go"
)

type Client struct {
	PostProcessors []PostProcessor
}

type clientCommunication struct {
	connection         *net.Conn
	p2pMessageEnvelope *eos.P2PMessageEnvelope
}

func (c *Client) Dial(address string) error {

	conn, err := net.Dial("tcp", address)

	if err != nil {
		return err
	}

	fmt.Println("Connected to: ", address)
	go c.handleConnection(conn, &Route{From: address})

	return nil
}

func (c *Client) handleConnection(connection net.Conn, route *Route) {

	decoder := eos.NewDecoder(bufio.NewReader(connection))

	for {

		var envelope eos.P2PMessageEnvelope
		fmt.Println("Waiting for payload")
		err := decoder.Decode(&envelope)
		if err != nil {
			fmt.Println("Connection error: ", err)
			continue
		}

		typeName, _ := envelope.Type.Name()
		fmt.Printf("Message received from [%s] with length: [%d] type: [%d - %s]\n", connection.RemoteAddr().String(), envelope.Length, envelope.Type, typeName)

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

	pp.P2PMessage = &msg

	for _, p := range c.PostProcessors {
		p.Handle(pp)
	}

	return nil
}
