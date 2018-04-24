package websocket

import (
	"log"
	"net/url"

	"fmt"

	"github.com/eoscanada/eos-go/proxy"
	"github.com/gorilla/websocket"
)

var Plugin = WebSocketPlugin{
	webSocketChannel: make(proxy.PostProcessChannel),
	serverAddress:    "Charless-MacBook-Pro-2.local:8080",
}

type WebSocketPlugin struct {
	webSocketChannel proxy.PostProcessChannel
	serverAddress    string
	connection       *websocket.Conn
}

func (p *WebSocketPlugin) getConnection() (conn *websocket.Conn, err error) {

	if p.connection == nil {
		fmt.Printf("New connection for [%v] on server [%s] \n", p, p.serverAddress)
		u := url.URL{Scheme: "ws", Host: p.serverAddress, Path: "/ws"}

		urlStr := u.String()
		fmt.Printf("URL to server [%s]\n", urlStr)

		c, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
		if err != nil {
			return nil, err
		}
		p.connection = c
	}

	fmt.Println("Returning connect")
	return p.connection, nil

}

func (p *WebSocketPlugin) Start() {

	//go handleWebSocket(p.webSocketChannel)

	var t proxy.PostProcessorPlugin
	t = p

	fmt.Println(t)

}

func (p *WebSocketPlugin) Channel() proxy.PostProcessChannel {
	return p.webSocketChannel
}

func (p *WebSocketPlugin) Handler() proxy.Handler {
	return HandleWebSocket
}

func HandleWebSocket(postProcessChannel proxy.PostProcessChannel) {

	for processable := range postProcessChannel {

		//u := url.URL{Scheme: "ws", Host: "Charless-MacBook-Pro-2.local:8080", Path: "/ws"}
		//fmt.Printf("connecting to %s", u.String())
		//
		//conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		conn, err := Plugin.getConnection()
		if err != nil {
			log.Fatal("dial:", err)
		}

		err = conn.WriteJSON(processable)
		if err != nil {
			// repost on channel ???/
			log.Fatal("JSON convertion err: ", err)
		}
	}
}
