package main

import (
	"fmt"

	"github.com/eoscanada/eos-go/p2p"
)

func main() {

	relay := p2p.NewRelay("0.0.0.0:6789", "localhost:9876")
	//relay.RegisterHandler(p2p.StringLoggerHandler)

	fmt.Println(relay.Start())
}
