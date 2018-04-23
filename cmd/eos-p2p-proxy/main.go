package main

import (
	"net"

	"fmt"

	"bufio"

	"flag"

	"plugin"

	"log"

	"reflect"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/proxy"
)

var Routes = []Route{
	{From: ":8900", To: "cbillett.eoscanada.com:9876"},
	{From: ":8901", To: "cbillett.eoscanada.com:9876"},
	{From: ":8902", To: "Charless-MacBook-Pro-2.local:19876"},
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
	P2PMessage            eos.P2PMessageEnvelope
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

type Router chan Communication

var router = make(Router)

func handleRouting(routingChannel Router) {

	for communication := range routingChannel {
		for _, channel := range routingChannels {

			channel <- communication
		}
	}
}

type PostProcessorChannel chan Communication

var postProcessorChannel = make(PostProcessorChannel)

func handlePostProcess(postProcessChannel PostProcessorChannel, postProcessChannels []proxy.P2PMessageChannel) {

	fmt.Println("Wait for comm on web socket")

	for communication := range postProcessChannel {

		msg, err := communication.P2PMessage.AsMessage()
		if err != nil {
			fmt.Println("Post processing err: ", err)
			continue
		}
		for _, c := range postProcessChannels {
			fmt.Printf("Sending message [%s] to channel [%s]\n", msg, reflect.TypeOf(c))
			c <- msg
		}
	}
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
		var msg eos.P2PMessageEnvelope

		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Connection error: ", err)
			forwardConnection.Close()
			return
		}

		typeName, _ := msg.Type.Name()
		fmt.Printf("Message received from [%s] with length: [%d] type: [%d - %s]\n", connection.RemoteAddr().String(), msg.Length, msg.Type, typeName)

		router <- Communication{
			DestinationConnection: forwardConnection,
			P2PMessage:            msg,
		}
	}
}

var routingChannels []chan Communication

type pluginFlags []string

func (p *pluginFlags) String() string {
	return "TODO"
}

func (p *pluginFlags) Set(value string) error {
	*p = append(*p, value)
	return nil
}

var plugins pluginFlags

func main() {

	done := make(chan bool)

	flag.Var(&plugins, "plugin", "Plugin so file path")
	flag.Parse()

	var postProcessChannels []proxy.P2PMessageChannel
	for _, p := range plugins {
		fmt.Println("Loading plugin: ", p)
		plug, err := plugin.Open(p)
		if err != nil {
			log.Fatal("Failed to load plugin: ", err)
		}
		pluginSymbol, err := plug.Lookup("Plugin")
		if err != nil {
			log.Fatal("Failed to load plugin: ", err)
		}

		var plugin proxy.PostProcessorPlugin
		plugin, ok := pluginSymbol.(proxy.PostProcessorPlugin)
		if !ok {
			log.Fatal("unexpected type from module symbol: ", reflect.TypeOf(pluginSymbol).String())
		}
		go plugin.Start()
		postProcessChannels = append(postProcessChannels, plugin.Channel())
	}

	routingChannels = []chan Communication{transmissionChannel, postProcessorChannel}

	go handleRouteAction(routeActionChannel)

	go handleRouting(router)
	go handleTransmission(transmissionChannel)
	go handlePostProcess(postProcessorChannel, postProcessChannels)

	for _, route := range Routes {

		routeActionChannel <- RouteAction{ActionType: AddRoute, Route: route}
	}

	<-done
}
