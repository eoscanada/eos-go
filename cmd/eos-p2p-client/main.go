package main

import (
	"log"

	"flag"

	"bytes"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/p2p"
)

var apiAddr = flag.String("api-addr", "http://localhost:8888", "RPC endpoint of the nodeos instance")
var p2pAddr = flag.String("p2p-addr", "localhost:9876", "P2P socket connection")
var signingKey = flag.String("signing-key", "5J5EE2cBDM4d3vWpKGcJsgiagsLVZkgWjJpxacz9mXodemXex6K", "Key to sign transactions we're about to blast")
var chainID = flag.String("chain-id", "0000000000000000000000000000000000000000000000000000000000000000", "Chain id")
var networkVersion = flag.Int("network-version", 25431, "Chain id")

func main() {

	flag.Parse()

	done := make(chan bool)

	api := eos.New(*apiAddr, bytes.Repeat([]byte{0}, 32))
	client := p2p.NewClient(*p2pAddr, api, p2p.DecodeHex(*chainID), int16(*networkVersion), *p2pAddr)
	client.RegisterHandler(p2p.HandlerFunc(p2p.LoggerHandler))
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
