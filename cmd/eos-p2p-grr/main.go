package main

import (
	"encoding/hex"
	"fmt"
	"log"

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
		switch m := packet.P2PMessage.(type) {
		case *eos.GoAwayMessage:
			log.Fatalf("handling message: go away: reason [%d]", m.Reason)

		case *eos.HandshakeMessage:
			//peer.SendHandshake(&p2p.HandshakeInfo{
			//	ChainID:                  chainID,
			//	LastIrreversibleBlockID:  m.LastIrreversibleBlockID,
			//	LastIrreversibleBlockNum: m.LastIrreversibleBlockNum,
			//	HeadBlockTime:            m.Time.Time,
			//	HeadBlockNum:             m.HeadNum,
			//	//HeadBlockID:              m.HeadID,
			//})
			//peer.SendNotice(m.HeadNum-1, m.LastIrreversibleBlockNum-1)
			//peer.SendSyncRequest(m.HeadNum-1, m.LastIrreversibleBlockNum-1)
		case *eos.NoticeMessage:
			//peer.SendRequest(0, 0)
			//peer.SendRequest(0, 0)
		default:
			//name, _ := m.GetType().Name()
			//fmt.Println("Read:", name)
		}

	}
}

func Start(peer *p2p.Peer, chainID string) error {

	errorChannel := make(chan error)

	readyChannel := peer.Connect(errorChannel)
	cID, err := hex.DecodeString(chainID)
	if err != nil {
		return fmt.Errorf("connect and start: parsing chain id: %s", err)
	}

	for {

		select {
		case <-readyChannel:
			go read(peer, errorChannel, cID)
			if chainID != "" {
				err = triggerHandshake(peer, cID)
				//peer.SendNotice(0, 100)
				//peer.SendSyncRequest(0, 100)
				//if err != nil {
				//	return fmt.Errorf("connect and start: trigger handshake: %s", err)
				//}
			}
		case err := <-errorChannel:
			return err
		}
	}
}

func triggerHandshake(peer *p2p.Peer, chainID eos.SHA256Bytes) error {
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
	fmt.Printf("<-- Sending handshake head %d lib %d\n", head.BlockNum, lib.BlockNum)
	err = peer.SendHandshake(dummyHandshakeInfo)
	if err != nil {
		return err
	}

	return peer.SendNotice(head.BlockNum, lib.BlockNum)
}

var api *eos.API

func main() {

	api = eos.New("http://localhost:8888")
	peer := p2p.NewOutgoingPeer("localhost:9876", "eos-proxy")

	Start(peer, "cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f")
	select {}

}
