package main

import (
	"net"

	"fmt"

	"bufio"

	"flag"

	"log"

	"reflect"

	"plugin"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/proxy"
	"github.com/eoscanada/eos-go/proxy/plugins/websocket"
)

var Routes = []*proxy.Route{
	{From: ":8900", To: "cbillett.eoscanada.com:9876"},
	{From: ":8901", To: "cbillett.eoscanada.com:9876"},
	{From: ":8902", To: "Charless-MacBook-Pro-2.local:19876"},
}

type ActionType int

const (
	AddRoute ActionType = iota
)

type RouteAction struct {
	ActionType ActionType
	Route      *proxy.Route
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
	Route                 *proxy.Route `json:"route"`
	DestinationConnection net.Conn
	P2PMessageEnvelope    *eos.P2PMessageEnvelope
}

type TransmissionChannel chan Communication

var transmissionChannel = make(TransmissionChannel)

func handleTransmission(channel TransmissionChannel) {

	for communication := range channel {

		encoder := eos.NewEncoder(communication.DestinationConnection)
		err := encoder.Encode(communication.P2PMessageEnvelope)
		if err != nil {
			fmt.Println("Sender error: ", err)
		} else {
			fmt.Printf("Message forwarded to [%s]\n", communication.DestinationConnection.RemoteAddr().String())
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

func handlePostProcess(postProcessChannel PostProcessorChannel, postProcessChannels []proxy.PostProcessChannel) {

	fmt.Println("Wait for comm on web socket")

	for communication := range postProcessChannel {

		pp := proxy.PostProcessable{
			Route:              communication.Route,
			P2PMessageEnvelope: communication.P2PMessageEnvelope,
		}

		//msg, err := communication.P2PMessageEnvelope.AsMessage()
		//if err != nil {
		//	fmt.Println("Post processing err: ", err)
		//	continue
		//} else {
		//	pp.P2PMessage = &msg
		//}

		for _, c := range postProcessChannels {
			fmt.Printf("Sending pp [%s] to channel [%s]\n", pp, reflect.TypeOf(c))
			c <- pp
		}
	}
}

func startForwarding(route *proxy.Route) {

	fmt.Printf("Starting forwarding [%s] -> [%s] \n", route.From, route.To)

	ln, err := net.Listen("tcp", route.From)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	for {
		fmt.Printf("Accepting connection on port [%s]\n", route.From)
		fromConn, err := ln.Accept()
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Printf("Connection on port [%s]\n", route.From)

		toConn, err := net.Dial("tcp", route.To)
		if err != nil {
			fmt.Println("error: ", err)
			fromConn.Close()
		} else {
			fmt.Println("Connected to: ", route.To)
			go handleConnection(fromConn, toConn, route)
			go handleConnection(toConn, fromConn, route)
		}
	}
}

func handleConnection(connection net.Conn, forwardConnection net.Conn, route *proxy.Route) (err error) {

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
			Route: route,
			DestinationConnection: forwardConnection,
			P2PMessageEnvelope:    &msg,
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

var pluginFiles pluginFlags
var plugins = make([]proxy.PostProcessorPlugin, 0)

func main() {

	done := make(chan bool)

	flag.Var(&pluginFiles, "plugin", "Plugin SO file path")
	flag.Parse()

	var postProcessChannels []proxy.PostProcessChannel

	for _, p := range pluginFiles {
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

		plugin.Start()
		plugins = append(plugins, plugin)
		postProcessChannels = append(postProcessChannels, plugin.Channel())
		h := plugin.Handler()
		go h(plugin.Channel())
	}

	//postProcessChannels = append(postProcessChannels, websocket.Plugin.Channel())
	go websocket.HandleWebSocket(websocket.Plugin.Channel())

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
