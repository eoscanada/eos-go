package p2p

import (
	"encoding/json"
	"fmt"
)

type Handler interface {
	Handle(msg *Packet)
}

type HandlerFunc func(packet *Packet)

func (f HandlerFunc) Handle(packet *Packet) {
	f(packet)
}

// LoggerHandler logs the messages back and forth.
var LoggerHandler = HandlerFunc(func(packet *Packet) {
	data, err := json.Marshal(packet)
	if err != nil {
		fmt.Println("logger plugin err: ", err)
		return
	}

	fmt.Println("logger - message : ", string(data))
})

// StringLoggerHandler simply prints the messages as they go through the client.
var StringLoggerHandler = HandlerFunc(func(packet *Packet) {
	name, _ := packet.Envelope.Type.Name()
	fmt.Printf(
		"type %s from %s to %s: %s\n",
		name,
		packet.Sender.Address,
		packet.Receiver.Address,
		packet.Envelope.P2PMessage,
	)
})
