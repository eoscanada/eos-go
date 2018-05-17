package main

import (
	"log"

	"flag"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/p2p"
)

var apiAddr = flag.String("api-addr", "http://stage5.testnets.eoscanada.com", "RPC endpoint of the nodeos instance")
var p2pAddr = flag.String("p2p-addr", "stage5.testnets.eoscanada.com:9876", "P2P socket connection")
var signingKey = flag.String("signing-key", "5JRFptLJfq16Qk9fZummqyKkhuaDjK5R1PAp1uGZ1E29SXVUfbJ", "Key to sign transactions we're about to blast")
var chainID = flag.String("chain-id", "0000000000000000000000000000000000000000000000000000000000000000", "Chain id")
var networkVersion = flag.Int("network-version", 25431, "Chain id")

func main() {

	flag.Parse()

	done := make(chan bool)

	privKey, err := ecc.NewPrivateKey(*signingKey)
	if err != nil {
		log.Fatal(err)
	}
	client := p2p.NewClient(*p2pAddr, p2p.DecodeHex(*chainID), int16(*networkVersion), privKey)
	client.RegisterHandler(p2p.HandlerFunc(p2p.LoggerHandler))
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
