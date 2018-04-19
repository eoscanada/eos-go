package eos

import (
	"encoding/binary"
	"errors"
	"io"

	"fmt"

	"github.com/eoscanada/eos-go/ecc"
)

// Work-in-progress p2p comms implementation
//
// See /home/abourget/build/eos3/plugins/net_plugin/include/eosio/net_plugin/protocol.hpp:219
//

type P2PMessageType byte

const (
	HandshakeMessageType P2PMessageType = iota
	GoAwayMessageType
	TimeMessageType
	NoticeMessageType
	RequestMessageType
	SyncRequestMessageType
	SignedBlockSummaryMessageType
	SignedBlockMessageType
	SignedTransactionMessageType
	PackedTransactionMessageType
)

var messageNames = []string{
	"Handshake",
	"GoAway",
	"Time",
	"Notice",
	"Request",
	"SyncRequest",
	"SignedBlockSummary",
	"SignedBlock",
	"SignedTransaction",
	"PackedTransaction",
}

func NewMessageType(aType byte) (t P2PMessageType, err error) {

	t = P2PMessageType(aType)
	if !t.isValid() {
		err = errors.New(fmt.Sprintf("unknown type [%d]", aType))
		return
	}

	return
}

func (t P2PMessageType) isValid() bool {

	index := byte(t)
	return int(index) < len(messageNames) && index >= 0

}

func (t P2PMessageType) Name() (string, bool) {

	index := byte(t)

	if !t.isValid() {
		return "Unknown", false
	}

	name := messageNames[index]
	return name, true
}

type P2PMessage struct {
	Length  uint32
	Type    P2PMessageType
	Payload []byte
}

func (a P2PMessage) MarshalBinary() ([]byte, error) {

	data := make([]byte, a.Length+4, a.Length+4)
	binary.LittleEndian.PutUint32(data[0:4], a.Length)
	data[4] = byte(a.Type)
	copy(data[5:], a.Payload)

	return data, nil
}

func (a *P2PMessage) UnmarshalBinaryRead(r io.Reader) (err error) {

	lengthBytes := make([]byte, 4, 4)
	_, err = r.Read(lengthBytes)
	if err != nil {
		fmt.Errorf("error: [%s]\n", err)
		return
	}

	size := binary.LittleEndian.Uint32(lengthBytes)

	payloadBytes := make([]byte, size, size)

	_, err = io.ReadFull(r, payloadBytes)

	if err != nil {
		return
	}
	//fmt.Printf("--> Payload length [%d] read count [%d]\n", size, count)

	if size < 1 {
		return errors.New("empty message")
	}

	//headerBytes := append(lengthBytes, payloadBytes[:int(math.Min(float64(10), float64(len(payloadBytes))))]...)

	//fmt.Printf("Length: [%s] Payload: [%s]\n", hex.EncodeToString(lengthBytes), hex.EncodeToString(payloadBytes[:int(math.Min(float64(7000), float64(len(payloadBytes))))]))

	messageType, err := NewMessageType(payloadBytes[0])
	if err != nil {
		return
	}

	*a = P2PMessage{
		Length:  size,
		Type:    messageType,
		Payload: payloadBytes[1:],
	}

	return nil
}

type HandshakeMessage struct {
	// net_plugin/protocol.hpp handshake_message
	NetworkVersion           int16         `json:"network_version"`
	ChainID                  HexBytes      `json:"chain_id"`
	NodeID                   HexBytes      `json:"node_id"` // sha256
	Key                      ecc.PublicKey `json:"key"`     // can be empty, producer key, or peer key
	Time                     Tstamp        `json:"time"`    // time?!
	Token                    HexBytes      `json:"token"`   // digest of time to prove we own the private `key`
	Signature                ecc.Signature `json:"sig"`     // can be empty if no key, signature of the digest above
	P2PAddress               string        `json:"p2p_address"`
	LastIrreversibleBlockNum uint32        `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  HexBytes      `json:"last_irreversible_block_id"`
	HeadNum                  uint32        `json:"head_num"`
	HeadID                   HexBytes      `json:"head_id"`
	OS                       string        `json:"os"`
	Agent                    string        `json:"agent"`
	Generation               int16         `json:"generation"`
}

type GoAwayReason uint8

const (
	GoAwayNoReason = uint8(iota)
	GoAwaySelfConnect
	GoAwayDuplicate
	GoAwayWrongChain
	GoAwayWrongVersion
	GoAwayForked
	GoAwayUnlinkable
	GoAwayBadTransaction
	GoAwayValidation
	GoAwayAuthentication
	GoAwayFatalOther
	GoAwayBenignOther
	GoAwayCrazy
)

type GoAwayMessage struct {
	GoAwayReason
}
