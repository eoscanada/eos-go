package main

import (
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eos-go/proxy"
)

var Plugin = WebSocketPlugin{
	webSocketChannel: make(proxy.P2PMessageChannel),
}

type WebSocketPlugin struct {
	webSocketChannel proxy.P2PMessageChannel
}

func (p *WebSocketPlugin) Start() {
	handleWebSocket(p.webSocketChannel)
}

func (p *WebSocketPlugin) Channel() proxy.P2PMessageChannel {
	return p.webSocketChannel
}

func handleWebSocket(messageChannel proxy.P2PMessageChannel) {

	for message := range messageChannel {

		fmt.Println("websocket received message: ", message)
		b, err := json.Marshal(message)
		if err != nil {
			fmt.Println("WebSocket err : ", err)
		}
		fmt.Println("WebSocket data ------> ", string(b))
	}
}
