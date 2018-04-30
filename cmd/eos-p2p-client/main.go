package main

import (
	"log"

	"flag"

	"bytes"

	"net/url"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/p2p"
)

var apiAddr = flag.String("api-addr", "http://localhost:8888", "RPC endpoint of the nodeos instance")
var p2pAddr = flag.String("p2p-addr", "localhost:9876", "P2P socket connection")
var signingKey = flag.String("signing-key", "", "Key to sign transactions we're about to blast")
var chainID = flag.String("chain-id", "00000000000000000000000000000000", "Chain id")
var networkVersion = flag.Int("network-version", 25431, "Chain id")

func main() {

	flag.Parse()

	apiAddrURL, err := url.Parse(*apiAddr)
	if err != nil {
		log.Fatalln("could not parse --api-addr:", err)
	}

	done := make(chan bool)

	api := eos.New(apiAddrURL, bytes.Repeat([]byte{0}, 32))
	client := p2p.NewClient(*p2pAddr, api, *chainID, int16(*networkVersion))
	client.RegisterHandler(p2p.HandlerFunc(p2p.LoggerHandler))
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
