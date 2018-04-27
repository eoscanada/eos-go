package main

import (
	"log"

	"github.com/eoscanada/eos-go/p2p"
)

func main() {

	done := make(chan bool)

	client := p2p.Client{
		Handlers: []p2p.Handler{
			p2p.LoggerHandler,
		},
	}

	err := client.Dial("testnet1.eos.io:9876", "testnet1.eos.io:8888")
	//err := client.Dial(":8902", ":8888")
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
