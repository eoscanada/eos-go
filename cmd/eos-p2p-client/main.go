package main

import (
	"log"

	"flag"

	"encoding/hex"

	"github.com/eoscanada/eos-go/p2p"
)

//var p2pAddr = flag.String("p2p-addr", "eosio.multibp.eoscanada.com:9876", "P2P socket connection")
var p2pAddr = flag.String("p2p-addr", "Charless-MacBook-Pro-2.local:9876", "P2P socket connection")
var signingKey = flag.String("signing-key", "5JRFptLJfq16Qk9fZummqyKkhuaDjK5R1PAp1uGZ1E29SXVUfbJ", "Key to sign transactions we're about to blast")
var chainID = flag.String("chain-id", "bf4dd1a5e59c6dc90bddb4f678ec24a2d2d4678b1513f5b483f134e586fc4643", "Chain id")
var networkVersion = flag.Int("network-version", 1206, "Network version")

func main() {

	flag.Parse()

	done := make(chan bool)
	cID, err := hex.DecodeString(*chainID)
	if err != nil {
		log.Fatal(err)
	}

	client := p2p.NewClient(*p2pAddr, cID, uint16(*networkVersion))
	if err != nil {
		log.Fatal(err)
	}
	client.RegisterHandler(p2p.HandlerFunc(p2p.LoggerHandler))
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
