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

	cID, err := hex.DecodeString("308cae83a690640be3726a725dde1fa72a845e28cfc63f28c3fa0a6ccdb6faf0")
	if err != nil {

		log.Fatal(err)
	}

	fmt.Println("P2P Client", *peer)
	client := p2p.NewClient(
		p2p.NewOutgoingPeer("localhost:9876", "eos-proxy", &p2p.HandshakeInfo{
			ChainID:      cID,
			HeadBlockNum: 1,
		}),
		true,
	)

	client.RegisterHandler(p2p.StringLoggerHandler)
	client.Start()
}
