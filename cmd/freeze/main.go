package main

import (
	"encoding/hex"
	"log"
	"net"

	"time"

	"fmt"

	"github.com/eoscanada/eos-go"
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
func readMessage(connection *p2p.Connection, msgChannel chan *eos.P2PMessageEnvelope, errChannel chan error) {
	for {
		msg, err := connection.Read()
		if err != nil {
			errChannel <- fmt.Errorf("read message from %s: %s", connection., err)
		}
		msgChannel <- msg
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
	peer.Connection.SendHandshake(dummyHandshakeInfo)
	fmt.Println("Handshake sent!")

	select {}
	fmt.Println("Done!")
}
