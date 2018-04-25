package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eos-go"
)

type PostProcessable struct {
	Route              *Route                  `json:"route"`
	P2PMessageEnvelope *eos.P2PMessageEnvelope `json:"p2p_message_envelope"`
	P2PMessage         *eos.P2PMessage         `json:"p2p_message"`
}

type PostProcessor interface {
	Handle(postProcessable PostProcessable)
}

type LoggerPostProcessor struct {
}

func (p *LoggerPostProcessor) Handle(postProcessable PostProcessable) {

	data, error := json.Marshal(postProcessable)
	if error != nil {
		fmt.Println("logger plugin err: ", error)
		return
	}

	fmt.Println("logger - message : ", string(data))

}
