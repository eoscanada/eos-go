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

type Handler func(processable PostProcessable)

var LoggerHandler = func(processable PostProcessable) {

	data, error := json.Marshal(processable)
	if error != nil {
		fmt.Println("logger plugin err: ", error)
		return
	}

	fmt.Println("logger - message : ", string(data))
}

var StringLoggerHandler = func(processable PostProcessable) {

	fmt.Printf("Message -> from [%s] to [%s] [%s]\n", processable.Route.From, processable.Route.To, processable.P2PMessage)
}
