package eos

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Work-in-progress p2p comms implementation
//
// See /home/abourget/build/eos3/plugins/net_plugin/include/eosio/net_plugin/protocol.hpp:219
//

type P2PMessageType byte

const (
	HandshakeMessageType P2PMessageType = iota // 0
	ChainSizeType
	GoAwayMessageType // 2
	TimeMessageType
	NoticeMessageType // 4
	RequestMessageType
	SyncRequestMessageType // 6
	SignedBlockType
	PackedTransactionMessageType // 8
)

type MessageReflectTypes struct {
	Name        string
	ReflectType reflect.Type
}

var messageAttributes = []MessageReflectTypes{
	{Name: "Handshake", ReflectType: reflect.TypeOf(HandshakeMessage{})},
	{Name: "ChainSize", ReflectType: reflect.TypeOf(ChainSizeMessage{})},
	{Name: "GoAway", ReflectType: reflect.TypeOf(GoAwayMessage{})},
	{Name: "Time", ReflectType: reflect.TypeOf(TimeMessage{})},
	{Name: "Notice", ReflectType: reflect.TypeOf(NoticeMessage{})},
	{Name: "Request", ReflectType: reflect.TypeOf(RequestMessage{})},
	{Name: "SyncRequest", ReflectType: reflect.TypeOf(SyncRequestMessage{})},
	{Name: "SignedBlock", ReflectType: reflect.TypeOf(SignedBlock{})},
	{Name: "PackedTransaction", ReflectType: reflect.TypeOf(PackedTransactionMessage{})},
}

var ErrUnknownMessageType = errors.New("unknown type")

func NewMessageType(aType byte) (t P2PMessageType, err error) {
	t = P2PMessageType(aType)
	if !t.isValid() {
		return t, ErrUnknownMessageType
	}

	return
}

func (t P2PMessageType) isValid() bool {
	index := byte(t)
	return int(index) < len(messageAttributes) && index >= 0
}

func (t P2PMessageType) Name() (string, bool) {
	index := byte(t)

	if !t.isValid() {
		return "Unknown", false
	}

	attr := messageAttributes[index]
	return attr.Name, true
}

func (t P2PMessageType) reflectTypes() (MessageReflectTypes, bool) {
	index := byte(t)

	if !t.isValid() {
		return MessageReflectTypes{}, false
	}

	attr := messageAttributes[index]
	return attr, true
}

type P2PMessageEnvelope struct {
	Length     uint32         `json:"length"`
	Type       P2PMessageType `json:"type"`
	Payload    []byte         `json:"-"`
	P2PMessage P2PMessage     `json:"message" eos:"-"`
}

func ReadP2PMessageData(r io.Reader) (envelope *P2PMessageEnvelope, err error) {
	data := make([]byte, 0)

	lengthBytes := make([]byte, 4, 4)
	_, err = r.Read(lengthBytes)
	if err != nil {
		return
	}

	data = append(data, lengthBytes...)

	size := binary.LittleEndian.Uint32(lengthBytes)

	payloadBytes := make([]byte, size, size)
	count, err := io.ReadFull(r, payloadBytes)

	if count != int(size) {
		err = fmt.Errorf("readfull not full read[%d] expected[%d]", count, size)
		return
	}

	if err != nil {
		fmt.Println("Connection error: ", err)
		return
	}

	data = append(data, payloadBytes...)

	envelope = &P2PMessageEnvelope{}
	decoder := NewDecoder(data)
	decoder.DecodeActions(false)
	err = decoder.Decode(envelope)
	if err != nil {
		fmt.Println("Failing data: ", hex.EncodeToString(data))
	}
	return
}
