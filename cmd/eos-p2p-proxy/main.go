package main

import (
	"github.com/eoscanada/eos-go/p2p"
)

func main() {

	proxy := p2p.Proxy{
		Routes: []*p2p.Route{
			{From: ":8900", To: "cbillett.eoscanada.com:9876"},
			{From: ":8901", To: "cbillett.eoscanada.com:9876"},
			{From: ":8902", To: "Charless-MacBook-Pro-2.local:19876"},
		},
		PostProcessors: []p2p.PostProcessor{
			&p2p.LoggerPostProcessor{},
		},
	}

	proxy.Start()

}
