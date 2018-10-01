package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/eoscanada/eos-go/p2p"
)

var peer = flag.String("peer", "localhost:9876", "peer to connect to")

func main() {
	flag.Parse()

	cID, err := hex.DecodeString("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f")
	if err != nil {

		log.Fatal(err)
	}

	fmt.Println("P2P Client", *peer)
	client := p2p.NewClient(
		p2p.NewOutgoingPeer("localhost:9876", "eos-proxy", &p2p.HandshakeInfo{
			ChainID:      cID,
			HeadBlockNum: 1,
		}),
		false,
	)

	client.RegisterHandler(p2p.StringLoggerHandler)
	client.Start()
}
