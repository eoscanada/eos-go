package main

import (
	"log"

	"github.com/eoscanada/eos-go/p2p"
)

func main() {

	done := make(chan bool)

	client := p2p.Client{
		PostProcessors: []p2p.PostProcessor{
			&p2p.LoggerPostProcessor{},
		},
	}

	err := client.Dial(":9876", ":8888")
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
