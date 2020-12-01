package main

import (
	"flag"
	"fmt"

	"github.com/dfuse-io/logging"
	"github.com/eoscanada/eos-go/p2p"
)

var peer = flag.String("peer", "", "peer")
var listeningAddress = flag.String("listening-address", "", "address on with the relay will listen")
var showLog = flag.Bool("v", false, "show detail log")

func main() {
	flag.Parse()

	if *showLog {
		logging.Set(logging.MustCreateLogger(), "github.com/eoscanada/eos-go/p2p")
	}
	defer p2p.SyncLogger()

	relay := p2p.NewRelay(*listeningAddress, *peer)
	relay.RegisterHandler(p2p.StringLoggerHandler)
	fmt.Println(relay.Start())
}
