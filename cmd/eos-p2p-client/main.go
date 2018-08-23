package main

import (
	"github.com/eoscanada/eos-go/p2p"
)

func main() {

	client := p2p.NewClient(
		p2p.NewOutgoingPeer("localhost:9876", "eos-proxy"),
	)

	client.RegisterHandler(p2p.StringLoggerHandler)
	client.Start("9bf6c5d3610260507f3a37340c43ff186c1810c984e9ad0b99b6fb8d6a3c94a3")

}
