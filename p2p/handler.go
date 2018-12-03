package p2p

import (
	"encoding/json"

	"go.uber.org/zap"
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
		logErr("Marshal err", err)
		return
	}

	p2pLog.Info("handler", zap.String("message", string(data)))
})

// StringLoggerHandler simply prints the messages as they go through the client.
var StringLoggerHandler = HandlerFunc(func(envelope *Envelope) {
	name, _ := envelope.Packet.Type.Name()
	p2pLog.Info(
		"handler Packet",
		zap.String("name", name),
		zap.String("sender", envelope.Sender.Address),
		zap.String("receiver", envelope.Receiver.Address),
		zap.Stringer("msg", envelope.Packet.P2PMessage), // this will use by String()
	)
})
