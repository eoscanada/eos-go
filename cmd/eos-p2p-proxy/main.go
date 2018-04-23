package main

import (
	"net"
	"fmt"
	"bufio"
	"encoding/json"
	"github.com/eoscanada/eos-go"
	"google.golang.org/api/monitoring/v3"
	"golang.org/x/net/context"
	"time"
	"math/rand"
	"os"
	"log"
	"golang.org/x/oauth2/google"
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
	SourceConnection      net.Conn
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

type WebSocketChannel chan Communication

var webSocketChannel = make(WebSocketChannel)

func handleWebSocket(webSocketChannel WebSocketChannel) {

	for communication := range webSocketChannel {

		msg, err := communication.P2PMessage.AsMessage()
		if err != nil {
			fmt.Println("websocket err: ", err)
			continue
		}

		b, err := json.Marshal(msg)
		fmt.Println("WebSocket data ------> ", string(b))
	}
}

type MonitoringChannel chan Communication

var monitoringChannel = make(MonitoringChannel)

func handleMonitoring(monitoringChannel MonitoringChannel) {

	for communication := range monitoringChannel {

		msg, err := communication.P2PMessage.AsMessage()
		if err != nil {
			continue
		}
		typeName, _ := communication.P2PMessage.Type.Name()
		fmt.Printf("Message received from [%s] with length: [%d] type: [%d - %s]\n", communication.SourceConnection.RemoteAddr().String(), communication.P2PMessage.Length, communication.P2PMessage.Type, typeName)

		b, err := json.Marshal(msg)
		fmt.Println("Monitoring data ------> ", string(b))
		///////////////////////////////////////////////////////////

		if len(os.Args) < 2 {
			fmt.Println("Usage: custommetric <project_id>")
			return
		}

		ctx := context.Background()
		s, err := createService(ctx)
		if err != nil {
			log.Fatal(err)
		}

		//projectID := os.Args[1]
		projectID := "eoscanada-sandbox-patrick-test"
		metricType := "custom.googleapis.com/patrick/test/monitoring3"

		if err := writeTimeSeriesValue(s, projectID, metricType, typeName, communication.SourceConnection.RemoteAddr().String()); err != nil {
			log.Fatal(err)
		}

		///////////////////////////////////////////////////////////
	}
}

func writeTimeSeriesValue(s *monitoring.Service, projectID, metricType string, typeName string, msgFrom string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	rand.Seed(time.Now().UTC().UnixNano())
	randVal := rand.Int63n(10)
	timeseries := monitoring.TimeSeries{
		Metric: &monitoring.Metric{
			Type: metricType,
			Labels: map[string]string{
				//"environment":  "STAGING",
				"msgTypeLabel": typeName,
				"msgFrom": msgFrom,
				//"msgTypeCode": communication.P2PMessage.Type,
			},
		},
		Resource: &monitoring.MonitoredResource{
			Labels: map[string]string{
				"project_id": projectID,
				//"instance_id": "test-instance",
				//"zone":        "us-central1-f",
			},
			Type: "global",
			//Type: "gce_instance",
		},
		Points: []*monitoring.Point{
			{
				Interval: &monitoring.TimeInterval{
					//StartTime: now,
					EndTime:   now,
				},
				Value: &monitoring.TypedValue{
					Int64Value: &randVal,
					//Int64Value: int64(&randVal),
				},
			},
		},
	}

	createTimeseriesRequest := monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{&timeseries},
	}

	//log.Printf("writeTimeseriesRequest: %s\n", formatResource(createTimeseriesRequest))
	_, err := s.Projects.TimeSeries.Create(projectResource(projectID), &createTimeseriesRequest).Do()
	if err != nil {
		return fmt.Errorf("Could not write time series value, %v ", err)
	}
	return nil
}

func createService(ctx context.Context) (*monitoring.Service, error) {
	hc, err := google.DefaultClient(ctx, monitoring.MonitoringScope)
	if err != nil {
		return nil, err
	}
	s, err := monitoring.New(hc)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func projectResource(projectID string) string {
	return "projects/" + projectID
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
			//TODO: verifier si la connection doit etre fermee des deux bords
			return
		}

		communicationRouter <- Communication{
			SourceConnection:      connection,
			DestinationConnection: forwardConnection,
			P2PMessage:            msg,
		}
	}
}

var routingChannels []chan Communication

func main() {

	done := make(chan bool)

	routingChannels = []chan Communication{transmissionChannel, webSocketChannel, monitoringChannel}

	go handleRouteAction(routeActionChannel)

	go handleRouting(communicationRouter)
	go handleTransmission(transmissionChannel)
	go handleWebSocket(webSocketChannel)
	go handleMonitoring(monitoringChannel)

	for _, route := range Routes {

		routeActionChannel <- RouteAction{ActionType: AddRoute, Route: route}
	}

	<-done
}
