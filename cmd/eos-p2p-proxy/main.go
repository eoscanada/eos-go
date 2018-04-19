package main

import (
	"net"

	"fmt"

	"bufio"

	"github.com/eoscanada/eos-go"
)

var Settings = []ForwardSetting{
	{From: ":8900", To: "cbillett.eoscanada.com:9876"},
	{From: ":8901", To: "cbillett.eoscanada.com:9876"},
	{From: ":8902", To: "localhost:19876"},
}

type ForwardSetting struct {
	From string
	To   string
}

type ForwardingSettingChannel chan ForwardSetting

var forwardingSettingChannel = make(ForwardingSettingChannel)

func handleForwardingSettings(channel ForwardingSettingChannel) {

	for setting := range channel {
		go startForwarding(setting)
	}
}

type Message struct {
	Source                string
	Destination           string
	DestinationConnection net.Conn
	P2PMessage            eos.P2PMessage
}

type SenderChannel chan Message

var senderChannel = make(SenderChannel)

func handleSend(channel SenderChannel) {

	for forward := range channel {

		encoder := eos.NewEncoder(forward.DestinationConnection)
		err := encoder.Encode(forward.P2PMessage)
		if err != nil {
			fmt.Println("Sender error: ", err)
		} else {
			fmt.Printf("Message forwarded to [%s]\n", forward.DestinationConnection.RemoteAddr().String())
		}
	}

}

func main() {

	done := make(chan bool)

	go handleSend(senderChannel)
	go handleForwardingSettings(forwardingSettingChannel)

	for _, forwardSetting := range Settings {

		forwardingSettingChannel <- forwardSetting
	}

	<-done

}

func startForwarding(setting ForwardSetting) {

	fmt.Printf("Starting forwarding [%s] -> [%s] \n", setting.From, setting.To)

	ln, err := net.Listen("tcp", setting.From)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	for {
		fmt.Printf("Accepting connection on port [%s]\n", setting.From)
		fromConn, err := ln.Accept()
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Printf("Connection on port [%s]\n", setting.From)

		toConn, err := net.Dial("tcp", setting.To)
		if err != nil {
			fmt.Println("error: ", err)
			fromConn.Close()
		} else {
			fmt.Println("Connected to: ", setting.To)
			go handleConnection(fromConn, toConn)
			go handleConnection(toConn, fromConn)
		}
	}
}

func handleConnection(connection net.Conn, forwardConnection net.Conn) (err error) {

	decoder := eos.NewDecoder(bufio.NewReader(connection))

	for {
		var msg eos.P2PMessage

		//fmt.Printf("Waiting for message from [%s]\n", connection.RemoteAddr().String())
		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Connection error: ", err)
			forwardConnection.Close()
			return
		}

		typeName, _ := msg.Type.Name()
		fmt.Printf("Message received from [%s] with length: [%d] type: [%d - %s]\n", connection.RemoteAddr().String(), msg.Length, msg.Type, typeName)

		senderChannel <- Message{
			DestinationConnection: forwardConnection,
			P2PMessage:            msg,
		}
	}
}
