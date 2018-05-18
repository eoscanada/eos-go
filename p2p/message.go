package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eos-go"
)

type Message struct {
	Route    *Route                  `json:"route"`
	Envelope *eos.P2PMessageEnvelope `json:"envelope"`
}

type Handler interface {
	Handle(msg Message)
}

type HandlerFunc func(msg Message)

func (f HandlerFunc) Handle(msg Message) {
	f(msg)
}

// LoggerHandler logs the messages back and forth.
var LoggerHandler = HandlerFunc(func(msg Message) {
	data, error := json.Marshal(msg)
	if error != nil {
		fmt.Println("logger plugin err: ", error)
		return
	}

	fmt.Println("logger - message : ", string(data))
})

// StringLoggerHandler simply prints the messages as they go through the client.
var StringLoggerHandler = HandlerFunc(func(msg Message) {
	fmt.Printf("Received message %T\n", msg.Envelope.P2PMessage)
})
