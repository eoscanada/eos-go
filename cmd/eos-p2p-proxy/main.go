package main

import (
	"flag"
	"fmt"

	"github.com/eoscanada/eos-go/p2p"
)

var peer1 = flag.String("peer1", "localhost:9876", "peer 1")
var peer2 = flag.String("peer2", "localhost:2222", "peer 2")
var chainID = flag.String("chain-id", "308cae83a690640be3726a725dde1fa72a845e28cfc63f28c3fa0a6ccdb6faf0", "peer 1")

func main() {

	fmt.Println("P2P Proxy")

	flag.Parse()
	//chainID, err := hex.DecodeString("aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906")

	proxy := p2p.NewProxy(
		p2p.NewOutgoingPeer(*peer1, "eos-proxy", nil),
		p2p.NewOutgoingPeer(*peer2, "eos-proxy", &p2p.HandshakeInfo{}),
	)

	//proxy := p2p.NewProxy(
	//	p2p.NewOutgoingPeer("localhost:9876", chainID, "eos-proxy", false),
	//	p2p.NewIncommingPeer("localhost:1111", chainID, "eos-proxy"),
	//)

	//proxy := p2p.NewProxy(
	//	p2p.NewIncommingPeer("localhost:2222", "eos-proxy"),
	//	p2p.NewIncommingPeer("localhost:1111", "eos-proxy"),
	//)

	proxy.RegisterHandler(p2p.StringLoggerHandler)
	fmt.Println(proxy.ConnectAndStart(*chainID))

}
