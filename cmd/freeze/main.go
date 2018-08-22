package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"

	"time"

	eos "github.com/eoscanada/eos-go"
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
func readMessageAndForward(fromPeer *p2p.Peer, toPeer *p2p.Peer) {
	for {
		envelope, err := fromPeer.Connection.Read()
		if err != nil {
			log.Fatal(fmt.Errorf("read message from: %s", err))
		}

		switch m := envelope.P2PMessage.(type) {
		case *eos.GoAwayMessage:
			log.Fatalf("handling message: go away: reason [%d]", m.Reason)
		case *eos.HandshakeMessage:
			nodeID, err := hex.DecodeString("9bf6c5d3610260507f3a37340c43ff186c1810c984e9ad0b99b6fb8d6a3c94a3")
			if err != nil {
				log.Fatal(fmt.Errorf("nodeID: %s", err))
			}
			fmt.Println("Handshake resent!")
			m.P2PAddress = "localhost:9875"
			m.NodeID = nodeID
			err = send(m, toPeer)
			if err != nil {
				log.Fatal(fmt.Errorf("HandshakeMessage: %s", err))
			}
		default:
			name, _ := envelope.Type.Name()
			fmt.Println("Drop:", name)
		}
	}
}

func send(message eos.P2PMessage, peer *p2p.Peer) error {

	err := peer.Connection.Write(message)
	if err != nil {
		return fmt.Errorf("writing message to %s: %s", peer.Connection.Address, err)
	}
	return nil
}

func main() {

	chainID, err := hex.DecodeString("9bf6c5d3610260507f3a37340c43ff186c1810c984e9ad0b99b6fb8d6a3c94a3")
	if err != nil {
		log.Fatal("Chain id:", err)
	}

	peer1 := peer("localhost:9876", chainID)
	//peer2 := peer("localhost:9875", chainID)
	go readMessageAndForward(peer1, peer1)
	//go readMessage(peer2.Connection, peer1.Connection)

	dummyHandshakeInfo := &p2p.HandshakeInfo{
		HeadBlockID:   make([]byte, 32),
		HeadBlockNum:  0,
		HeadBlockTime: time.Now(),
	}

	peer1.Connection.SendHandshake(dummyHandshakeInfo)
	//peer2.Connection.SendHandshake(dummyHandshakeInfo)
	fmt.Println("Handshake sent!")

	select {}
	fmt.Println("Done!")
}
