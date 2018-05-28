package main

import (
	"log"

	"flag"

	"bytes"

	"github.com/eoscanada/eos-go/p2p"
)

//var p2pAddr = flag.String("p2p-addr", "eosio.multibp.eoscanada.com:9876", "P2P socket connection")
var p2pAddr = flag.String("p2p-addr", "Charless-MacBook-Pro-2.local:9876", "P2P socket connection")
var signingKey = flag.String("signing-key", "5JRFptLJfq16Qk9fZummqyKkhuaDjK5R1PAp1uGZ1E29SXVUfbJ", "Key to sign transactions we're about to blast")
var chainID = flag.String("chain-id", "0000000000000000000000000000000000000000000000000000000000000000", "Chain id")
var networkVersion = flag.Int("network-version", 1, "Network version")

func main() {

	flag.Parse()

	done := make(chan bool)
	//eos.Debug = true
	client := p2p.NewClient(*p2pAddr, bytes.Repeat([]byte{0}, 32), uint16(*networkVersion))
	client.RegisterHandler(p2p.HandlerFunc(p2p.LoggerHandler))
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
