package main

import (
	"flag"
	"fmt"

	"github.com/eoscanada/eos-go/p2p"
)

var peer = flag.String("peer", "localhost:9876", "peer to connect to")
var chainID = flag.String("chain-id", "", "chain id of the peer")

func main() {
	fmt.Println("P2P Client", *peer)

	flag.Parse()
	client := p2p.NewClient(
		p2p.NewOutgoingPeer(*peer, "eos-proxy"),
	)

	client.RegisterHandler(p2p.StringLoggerHandler)
	client.Start(*chainID)

}
