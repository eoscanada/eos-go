package main

import (
	"fmt"

	"flag"

	"github.com/eoscanada/eos-go/p2p"
)

var peer = flag.String("peer", "", "peer")
var listeningAddress = flag.String("listening-address", "", "address on with the relay will listen")

func main() {

	flag.Parse()

	relay := p2p.NewRelay(*listeningAddress, *peer)

	relay.RegisterHandler(p2p.StringLoggerHandler)

	fmt.Println(relay.Start())
}
