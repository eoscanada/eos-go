package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/eoscanada/eos-go/p2p"
	"github.com/streamingfast/logging"
)

var peer = flag.String("peer", "localhost:9876", "peer to connect to")
var chainID = flag.String("chain-id", "cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f", "net chainID to connect to")
var showLog = flag.Bool("v", false, "show detail log")

func main() {
	flag.Parse()

	if *showLog {
		logging.Set(logging.MustCreateLogger(), "github.com/eoscanada/eos-go/p2p")
	}
	defer p2p.SyncLogger()

	cID, err := hex.DecodeString(*chainID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("P2P Client ", *peer, " With Chain ID :", *chainID)
	client := p2p.NewClient(
		p2p.NewOutgoingPeer(*peer, "eos-proxy", &p2p.HandshakeInfo{
			ChainID:      cID,
			HeadBlockNum: 1,
		}),
		false,
	)

	client.RegisterHandler(p2p.StringLoggerHandler)
	client.Start()
}
