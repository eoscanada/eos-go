package main

import (
	"fmt"

	"github.com/eoscanada/eos-go/p2p"
)

func main() {

	//chainID, err := hex.DecodeString("aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906")

	proxy := p2p.NewProxy(
		p2p.NewOutgoingPeer("localhost:9876", "eos-proxy"),
		p2p.NewOutgoingPeer("localhost:2222", "eos-proxy"),
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
	fmt.Println(proxy.ConnectAndStart("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))

}
