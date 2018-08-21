package main

import (
	"log"
	"net"

	"encoding/hex"

	"github.com/eoscanada/eos-go/p2p"
)

func peer(address string, chainID []byte) *p2p.Peer {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Dial %s id: %s", address, err)
	}

	originConnection := p2p.NewConnection(address, chainID, "eos-proxy", conn)
	return p2p.NewPeer(originConnection)

}

func main() {

	chainID, err := hex.DecodeString("9bf6c5d3610260507f3a37340c43ff186c1810c984e9ad0b99b6fb8d6a3c94a3")
	//chainID, err := hex.DecodeString("aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906")
	if err != nil {
		log.Fatal("Chain id:", err)
	}

	route := &p2p.Route{
		Origin: peer("localhost:9876", chainID),
		//Origin:      peer("35.203.0.168:9876", chainID),
		Destination: peer("localhost:9875", chainID),
	}

	proxy := p2p.NewProxy(route)

	proxy.Start()

}
