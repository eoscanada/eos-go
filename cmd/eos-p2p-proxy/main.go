package main

import (
	"github.com/eoscanada/eos-go/p2p"
)

func main() {

	proxy := p2p.Proxy{
		Routes: []*p2p.Route{
			{From: ":8902", To: ":9876"},
		},
		PostProcessors: []p2p.PostProcessor{
			&p2p.LoggerPostProcessor{},
		},
	}

	proxy.Start()

}
