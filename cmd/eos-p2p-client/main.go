package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/eoscanada/eos-go/p2p"
)

var peer = flag.String("peer", "localhost:9876", "peer to connect to")

func main() {
	flag.Parse()

	flag.Parse()
	fmt.Println("P2P Client", *peer)
	client := p2p.NewClient(
		p2p.NewOutgoingPeer("localhost:9876", "eos-proxy", &p2p.HandshakeInfo{HeadBlockNum: 0, LastIrreversibleBlockNum: 0, HeadBlockTime: time.Now()}),
	)

	client.RegisterHandler(p2p.StringLoggerHandler)
	client.Start()
}
