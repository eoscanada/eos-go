package main

import (
	"encoding/hex"
	"log"
	"net"

	"time"

	"fmt"

	"github.com/eoscanada/eos-go/p2p"
)

func peer(address string, chainID []byte) *p2p.Peer {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Dial %s id: %s", address, err)
	}

	originConnection := p2p.NewConnection(address, chainID, "eos-proxy", conn)
	return p2p.NewPeer(originConnection)

}
func readMessage(connection *p2p.Connection) {
	for {
		msg, err := connection.Read()
		if err != nil {
			log.Fatal(fmt.Errorf("read message from: %s", err))
		}
		fmt.Println("MSG:", msg)
	}
}

func main() {

	chainID, err := hex.DecodeString("9bf6c5d3610260507f3a37340c43ff186c1810c984e9ad0b99b6fb8d6a3c94a3")
	if err != nil {
		log.Fatal("Chain id:", err)
	}

	dummyHandshakeInfo := &p2p.HandshakeInfo{
		HeadBlockID:   make([]byte, 32),
		HeadBlockNum:  0,
		HeadBlockTime: time.Now(),
	}
	peer := peer("localhost:9876", chainID)
	go readMessage(peer.Connection)
	peer.Connection.SendHandshake(dummyHandshakeInfo)
	peer.Connection.SendSyncRequest(1, 10)
	fmt.Println("Handshake sent!")

	select {}
	fmt.Println("Done!")
}
