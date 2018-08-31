package main

import (
	"flag"
	"fmt"

	"github.com/eoscanada/eos-go/p2p"
)

var peer1 = flag.String("peer1", "localhost:9876", "peer 1")
var peer2 = flag.String("peer2", "localhost:2222", "peer 2")
var chainID = flag.String("chain-id", "9bf6c5d3610260507f3a37340c43ff186c1810c984e9ad0b99b6fb8d6a3c94a3", "peer 1")

func main() {

	fmt.Println("P2P Proxy")

	flag.Parse()
	//chainID, err := hex.DecodeString("aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906")

	proxy := p2p.NewProxy(
		p2p.NewOutgoingPeer(*peer1, "eos-proxy"),
		p2p.NewOutgoingPeer(*peer2, "eos-proxy"),
	)

	//proxy := p2p.NewProxy(
	//	p2p.NewOutgoingPeer("localhost:9876", chainID, "eos-proxy", false),
	//	p2p.NewIncommingPeer("localhost:1111", chainID, "eos-proxy"),
	//)

	//proxy := p2p.NewProxy(
	//	p2p.NewIncommingPeer("localhost:2222", chainID, "eos-proxy"),
	//	p2p.NewIncommingPeer("localhost:1111", chainID, "eos-proxy"),
	//)

	proxy.RegisterHandler(p2p.StringLoggerHandler)
	fmt.Println(proxy.ConnectAndStart(*chainID))

}
