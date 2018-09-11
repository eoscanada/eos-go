package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/p2p"
)

func read(peer *p2p.Peer, errChannel chan error, chainID eos.SHA256Bytes) {
	for {
		packet, err := peer.Read()
		if err != nil {
			errChannel <- fmt.Errorf("read message from %s: %s", peer.Address, err)
		}

		name, _ := packet.P2PMessage.GetType().Name()
		//fmt.Println("Read:", name)
		fmt.Printf("--> Received Msg: %s :: %s\n", name, packet.P2PMessage)
		//switch m := packet.P2PMessage.(type) {
		//case *eos.GoAwayMessage:
		//	log.Fatalf("handling message: go away: reason [%d]", m.Reason)
		//
		//case *eos.HandshakeMessage:
		//	//handshakeInfo := &p2p.HandshakeInfo{
		//	//	ChainID:                  chainID,
		//	//	LastIrreversibleBlockID:  m.LastIrreversibleBlockID,
		//	//	LastIrreversibleBlockNum: m.LastIrreversibleBlockNum,
		//	//	HeadBlockTime:            m.Time.Time,
		//	//	HeadBlockNum:             m.HeadNum,
		//	//	HeadBlockID:              m.HeadID,
		//	//}
		//	//if err := sendHanshake(handshakeInfo, peer); err != nil {
		//	//	panic(err)
		//	//}
		//	triggerHandshake(peer, chainID)
		//	peer.SendNotice(m.HeadNum-10, m.LastIrreversibleBlockNum-10)
		//	peer.SendTime()
		//
		//case *eos.NoticeMessage:
		//	//peer.SendRequest(0, 0)
		//	//peer.SendRequest(0, 0)
		//default:
		//	//name, _ := m.GetType().Name()
		//	//fmt.Println("Read:", name)
		//}

	}
}

func Start(peer *p2p.Peer, chainID string) error {

	fmt.Println("Getting info from server")
	api = eos.New("http://localhost:8888")
	info, err := api.GetInfo()
	if err != nil {
		return fmt.Errorf("connect and start: api get info: %s", err)
	}

	cID, err := hex.DecodeString(chainID)
	if err != nil {
		return fmt.Errorf("connect and start: parsing chain id: %s", err)
	}

	dummyHandshakeInfo := &p2p.HandshakeInfo{
		ChainID:                  cID,
		HeadBlockNum:             info.HeadBlockNum,
		HeadBlockID:              info.HeadBlockID,
		LastIrreversibleBlockNum: info.LastIrreversibleBlockNum,
		LastIrreversibleBlockID:  info.LastIrreversibleBlockID,
		HeadBlockTime:            info.HeadBlockTime.Time,
	}

	errorChannel := make(chan error)

	readyChannel := peer.Connect(errorChannel)

	for {

		select {
		case <-readyChannel:
			go read(peer, errorChannel, cID)
			time.Sleep(500 * time.Millisecond)
			sendHanshake(dummyHandshakeInfo, peer)
			peer.SendNotice(info.HeadBlockNum, info.LastIrreversibleBlockNum)
		case err := <-errorChannel:
			return err
		}
	}
}

func triggerUpToDateHandshake(peer *p2p.Peer, chainID eos.SHA256Bytes) error {
	info, err := api.GetInfo()
	if err != nil {
		return err
	}

	head, err := api.GetBlockByNum(info.LastIrreversibleBlockNum)
	if err != nil {
		return err
	}
	lib, err := api.GetBlockByNum(info.LastIrreversibleBlockNum - 5)
	if err != nil {
		return err
	}
	dummyHandshakeInfo := &p2p.HandshakeInfo{
		ChainID:                  chainID,
		HeadBlockNum:             head.BlockNum,
		HeadBlockID:              head.ID,
		LastIrreversibleBlockNum: lib.BlockNumber(),
		LastIrreversibleBlockID:  lib.ID,
		HeadBlockTime:            info.HeadBlockTime.Time,
	}
	return sendHanshake(dummyHandshakeInfo, peer)
}

func triggerHandshake(peer *p2p.Peer, chainID eos.SHA256Bytes) error {
	dummyHandshakeInfo := &p2p.HandshakeInfo{
		ChainID:                  chainID,
		HeadBlockNum:             0,
		LastIrreversibleBlockNum: 0,
	}
	return sendHanshake(dummyHandshakeInfo, peer)
}

func sendHanshake(handshakeInfo *p2p.HandshakeInfo, toPeer *p2p.Peer) error {
	fmt.Printf("<-- Sending handshake head %d lib %d\n", handshakeInfo.HeadBlockNum, handshakeInfo.LastIrreversibleBlockNum)
	if err := toPeer.SendHandshake(handshakeInfo); err != nil {
		return err
	}
	return nil
}

var api *eos.API

func main() {

	//peer := p2p.NewOutgoingPeer("35.203.0.168:9876", "eos-proxy")
	peer := p2p.NewOutgoingPeer("localhost:9876", "eos-proxy")

	//Start(peer, "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906") //mainnet
	Start(peer, "cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f")
	select {}

}
