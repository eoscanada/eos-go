package main

import (
	"flag"
	"log"
	"time"

	"encoding/hex"

	"github.com/eoscanada/eos-go/p2p"
)

var p2pAddr = flag.String("p2p-addr", "peering.mainnet.eoscanada.com:9876", "P2P socket connection")
var chainID = flag.String("chain-id", "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906", "Chain id")
var networkVersion = flag.Int("network-version", 1206, "Network version")

func main() {
	flag.Parse()

	chainID, err := hex.DecodeString("9bf6c5d3610260507f3a37340c43ff186c1810c984e9ad0b99b6fb8d6a3c94a3")
	if err != nil {
		log.Fatal("Chain id:", err)
	}

	client := p2p.NewClient(
		p2p.NewOutgoingPeer("localhost:9876", chainID, "eos-proxy", &p2p.HandshakeInfo{HeadBlockNum: 0, LastIrreversibleBlockNum: 0, HeadBlockTime: time.Now()}),
	)

	client.RegisterHandler(p2p.StringLoggerHandler)
	client.Start()
}
