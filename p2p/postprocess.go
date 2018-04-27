package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eos-go"
)

type PostProcessable struct {
	Route              *Route                  `json:"route"`
	P2PMessageEnvelope *eos.P2PMessageEnvelope `json:"p2p_message_envelope"`
	P2PMessage         eos.P2PMessage          `json:"p2p_message"`
}

type Handler interface {
	Handle(msg PostProcessable)
}

type HandlerFunc func(msg PostProcessable)

func (f HandlerFunc) Handle(msg PostProcessable) {
	f(msg)
}

// LoggerHandler logs the messages back and forth.
var LoggerHandler = HandlerFunc(func(msg PostProcessable) {
	data, error := json.Marshal(msg)
	if error != nil {
		fmt.Println("logger plugin err: ", error)
		return
	}

	fmt.Println("logger - message : ", string(data))
})

// StringLoggerHandler simply prints the messages as they go through the client.
var StringLoggerHandler = HandlerFunc(func(msg PostProcessable) {
	fmt.Printf("Message -> from [%s] to [%s] [%s]\n", msg.Route.From, msg.Route.To, msg.P2PMessage)
})
