package main

import (
	"net"

	"fmt"

	"bufio"

	"github.com/eoscanada/eos-go"
)

var Routes = []Route{
	{From: ":8900", To: "cbillett.eoscanada.com:9876"},
	{From: ":8901", To: "cbillett.eoscanada.com:9876"},
	{From: ":8902", To: "localhost:19876"},
}

type ActionType int

const (
	AddRoute ActionType = iota
)

type Route struct {
	From string
	To   string
}

type RouteAction struct {
	ActionType ActionType
	Route
}

type RouteActionChannel chan RouteAction

var routeActionChannel = make(RouteActionChannel)

func handleRouteAction(channel RouteActionChannel) {

	for routeAction := range channel {
		//todo : handle action type
		go startForwarding(routeAction.Route)
	}
}

type Communication struct {
	Source                string
	Destination           string
	DestinationConnection net.Conn
	P2PMessage            eos.P2PMessage
}

type TransmissionChannel chan Communication

var transmissionChannel = make(TransmissionChannel)

func handleTransmission(channel TransmissionChannel) {

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

type CommunicationRouter chan Communication

var communicationRouter = make(CommunicationRouter)

func handleRouting(routingChannel CommunicationRouter) {

	for communication := range routingChannel {
		for _, channel := range routingChannels {
			channel <- communication
		}
	}
}

var routingChannels []chan Communication

func main() {

	done := make(chan bool)

	go handleTransmission(transmissionChannel)
	go handleRouteAction(routeActionChannel)
	go handleRouting(communicationRouter)

	routingChannels = []chan Communication{transmissionChannel}

	for _, route := range Routes {

		routeActionChannel <- RouteAction{ActionType: AddRoute, Route: route}
	}

	<-done

}

func startForwarding(setting Route) {

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

		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Connection error: ", err)
			forwardConnection.Close()
			return
		}

		typeName, _ := msg.Type.Name()
		fmt.Printf("Message received from [%s] with length: [%d] type: [%d - %s]\n", connection.RemoteAddr().String(), msg.Length, msg.Type, typeName)

		transmissionChannel <- Communication{
			DestinationConnection: forwardConnection,
			P2PMessage:            msg,
		}
	}
}
