package p2p

import (
	"encoding/json"
	"fmt"
)

type Handler interface {
	Handle(envelope *Envelope)
}

type HandlerFunc func(envelope *Envelope)

func (f HandlerFunc) Handle(envelope *Envelope) {
	f(envelope)
}

// LoggerHandler logs the messages back and forth.
var LoggerHandler = HandlerFunc(func(envelope *Envelope) {
	data, err := json.Marshal(envelope)
	if err != nil {
		fmt.Println("logger plugin err: ", err)
		return
	}

	fmt.Println("logger - message : ", string(data))
})

// StringLoggerHandler simply prints the messages as they go through the client.
var StringLoggerHandler = HandlerFunc(func(envelope *Envelope) {
	name, _ := envelope.Packet.Type.Name()
	fmt.Printf(
		"type %s from %s to %s: %s\n",
		name,
		envelope.Sender.Address,
		envelope.Receiver.Address,
		envelope.Packet.P2PMessage,
	)
})
