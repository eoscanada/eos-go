package main

import (
	"log"
	"time"

	"flag"

	"encoding/hex"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/p2p"
)

var p2pAddr = flag.String("p2p-addr", "peering.mainnet.eoscanada.com:9876", "P2P socket connection")
var chainID = flag.String("chain-id", "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906", "Chain id")
var networkVersion = flag.Int("network-version", 1206, "Network version")

func main() {

	flag.Parse()

	done := make(chan bool)
	cID, err := hex.DecodeString(*chainID)
	if err != nil {
		log.Fatal(err)
	}

	api := eos.New("http://mainnet.eoscanada.com")
	info, err := api.GetInfo()
	if err != nil {
		log.Fatal("Error getting info: ", err)
	}

	client := p2p.NewClient(*p2pAddr, cID, uint16(*networkVersion))
	if err != nil {
		log.Fatal(err)
	}
	client.RegisterHandler(p2p.HandlerFunc(p2p.LoggerHandler))
	time.Sleep(120 * time.Second)

	err = client.ConnectAndSync(info.HeadBlockNum, info.HeadBlockID, info.HeadBlockTime.Time, 0, make([]byte, 32))
	//err = client.ConnectRecent()
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
