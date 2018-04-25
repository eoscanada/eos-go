package p2p

import (
	"net"

	"bufio"
	"fmt"

	"encoding/hex"
	"log"

	"github.com/eoscanada/eos-go"
)

type actionType int

const (
	addRoute actionType = iota
)

type routeAction struct {
	actionType actionType
	route      *Route
}

type routeCommunication struct {
	Route                 *Route `json:"route"`
	DestinationConnection net.Conn
	P2PMessageEnvelope    *eos.P2PMessageEnvelope
}

var routerActionChannel = make(chan routeAction)

func (p *Proxy) handleRouteAction(channel chan routeAction) {

	for routeAction := range channel {
		//todo : handle action type
		go p.startForwarding(routeAction.route)
	}
}

var transmissionChannel = make(chan routeCommunication)

func (p *Proxy) handleTransmission(channel chan routeCommunication) {

	for communication := range channel {

		encoder := eos.NewEncoder(communication.DestinationConnection)
		err := encoder.Encode(communication.P2PMessageEnvelope)
		if err != nil {
			fmt.Println("Sender error: ", err)
		}
	}
}

var router = make(chan routeCommunication)
var routingChannels []chan routeCommunication

func (p *Proxy) handleRouting(routingChannel chan routeCommunication) {

	for communication := range routingChannel {
		for _, channel := range routingChannels {

			channel <- communication
		}
	}
}

var postProcessChannel = make(chan routeCommunication)

func (p *Proxy) handlePostProcess(postProcessChannel chan routeCommunication, postProcessorChannels []chan PostProcessable) {

	for communication := range postProcessChannel {

		pp := PostProcessable{
			Route:              communication.Route,
			P2PMessageEnvelope: communication.P2PMessageEnvelope,
		}

		msg, err := communication.P2PMessageEnvelope.AsMessage()
		if err != nil {

			env := communication.P2PMessageEnvelope

			msgData, err := eos.MarshalBinary(env)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Post process failed for message type [%d] len[%d] with data [%s]\n", env.Type, env.Length, hex.EncodeToString(msgData))
			log.Fatal("Post processing err: ", err)
			continue
		}

		pp.P2PMessage = &msg

		for _, c := range postProcessorChannels {
			c <- pp
		}
	}
}

func (p *Proxy) handlePluginPostProcess(postProcessor PostProcessor, channel chan PostProcessable) {

	for postProcessable := range channel {
		postProcessor.Handle(postProcessable)
	}
}

func (p *Proxy) startForwarding(route *Route) {

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
			go p.handleConnection(fromConn, toConn, route)
			go p.handleConnection(toConn, fromConn, route)
		}
	}
}

func (p *Proxy) handleConnection(connection net.Conn, forwardConnection net.Conn, route *Route) (err error) {

	decoder := eos.NewDecoder(bufio.NewReader(connection))

	for {
		var msg eos.P2PMessageEnvelope

		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Connection error: ", err)
			forwardConnection.Close()
			return // handleConnection will be restarted in startForwarding
		}

		router <- routeCommunication{
			Route:                 route,
			DestinationConnection: forwardConnection,
			P2PMessageEnvelope:    &msg,
		}
	}
}

type Proxy struct {
	Routes         []*Route
	PostProcessors []PostProcessor
}

func (p *Proxy) Start() {

	done := make(chan bool)

	var postProcessorChannels []chan PostProcessable

	for _, plugin := range p.PostProcessors {

		pc := make(chan PostProcessable)
		postProcessorChannels = append(postProcessorChannels, pc)
		go p.handlePluginPostProcess(plugin, pc)

	}

	routingChannels = []chan routeCommunication{transmissionChannel, postProcessChannel}

	go p.handleRouteAction(routerActionChannel)
	go p.handleRouting(router)
	go p.handleTransmission(transmissionChannel)
	go p.handlePostProcess(postProcessChannel, postProcessorChannels)

	for _, route := range p.Routes {

		routerActionChannel <- routeAction{actionType: addRoute, route: route}
	}

	fmt.Println("Proxy started")
	<-done
	fmt.Println("Proxy will stop")
}
