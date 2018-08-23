package main

import (
	"fmt"

	"github.com/eoscanada/eos-go/p2p"
)

func main() {

	//chainID, err := hex.DecodeString("aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906")

	proxy := p2p.NewProxy(
		p2p.NewOutgoingPeer("localhost:9876", "eos-proxy"),
		p2p.NewOutgoingPeer("localhost:9875", "eos-proxy"),
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
	//fmt.Println(proxy.ConnectAndStart("aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906"))
	fmt.Println(proxy.ConnectAndStart("9bf6c5d3610260507f3a37340c43ff186c1810c984e9ad0b99b6fb8d6a3c94a3"))

}
