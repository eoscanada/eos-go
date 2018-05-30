package p2p

import (
	"net"

	"bufio"
	"fmt"

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
	Envelope              *eos.P2PMessageEnvelope
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

		//_, err := communication.DestinationConnection.Write(communication.Envelope.Payload)
		//buf := new(bytes.Buffer)
		//encoder := eos.NewEncoder(buf)
		//err := encoder.Encode(communication.Envelope)
		//if err != nil {
		//	fmt.Println("Sender encode error: ", err)
		//}

		//fmt.Println("Data to send: ", hex.EncodeToString(buf.Bytes()))

		_, err := communication.DestinationConnection.Write(communication.Envelope.Raw)
		if err != nil {
			fmt.Println("Sender comm error: ", err)
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

func (p *Proxy) handlePostProcess(postProcessChannel chan routeCommunication, postProcessorChannels []chan Message) {

	for communication := range postProcessChannel {

		pp := Message{
			Route:    communication.Route,
			Envelope: communication.Envelope,
		}

		for _, c := range postProcessorChannels {
			c <- pp
		}
	}
}

func (p *Proxy) handlePluginPostProcess(handle Handler, channel chan Message) {

	for postProcessable := range channel {
		handle.Handle(postProcessable)
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
			go p.handleConnection(toConn, fromConn, &Route{From: route.To, To: route.From})
		}
	}
}

func (p *Proxy) handleConnection(connection net.Conn, forwardConnection net.Conn, route *Route) (err error) {

	r := bufio.NewReader(connection)

	for {

		envelope, err := eos.ReadP2PMessageData(r)
		if err != nil {
			fmt.Printf("WARNING: Connection from [%s] to [%s] : %s\n ", route.From, route.To, err)
			forwardConnection.Close()
			//log.Fatal("Handle connection, ", err)
			return err
		}

		router <- routeCommunication{
			Route:                 route,
			DestinationConnection: forwardConnection,
			Envelope:              envelope,
		}

	}
}

type Proxy struct {
	Routes   []*Route
	Handlers []Handler
}

func (p *Proxy) Start() {

	done := make(chan bool)

	var postProcessorChannels []chan Message

	for _, plugin := range p.Handlers {

		pc := make(chan Message)
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
