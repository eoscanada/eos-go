package main

import (
	"flag"
	"fmt"

	"github.com/eoscanada/eos-go/p2p"
)

var peer = flag.String("peer", "localhost:9876", "peer to connect to")
var chainID = flag.String("chain-id", "308cae83a690640be3726a725dde1fa72a845e28cfc63f28c3fa0a6ccdb6faf0", "chain id of the peer")

func main() {

	flag.Parse()
	fmt.Println("P2P Client", *peer)
	client := p2p.NewClient(
		p2p.NewOutgoingPeer(*peer, "eos-proxy"),
	)

	client.RegisterHandler(p2p.StringLoggerHandler)
	client.Start(*chainID)

}
