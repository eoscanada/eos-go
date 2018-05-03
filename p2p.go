package eos

import (
	"bytes"
	"encoding/binary"
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

type MessageAttributes struct {
	Name        string
	ReflectType reflect.Type
}

var messageAttributes = []MessageAttributes{
	{Name: "Handshake", ReflectType: reflect.TypeOf(HandshakeMessage{})},
	{Name: "GoAway", ReflectType: reflect.TypeOf(GoAwayMessage{})},
	{Name: "Time", ReflectType: reflect.TypeOf(TimeMessage{})},
	{Name: "Notice", ReflectType: reflect.TypeOf(NoticeMessage{})},
	{Name: "Request", ReflectType: reflect.TypeOf(RequestMessage{})},
	{Name: "SyncRequest", ReflectType: reflect.TypeOf(SyncRequestMessage{})},
	{Name: "SignedBlockSummary", ReflectType: reflect.TypeOf(SignedBlockSummaryMessage{})},
	{Name: "SignedBlock", ReflectType: reflect.TypeOf(SignedBlockMessage{})},
	{Name: "SignedTransaction", ReflectType: reflect.TypeOf(SignedTransactionMessage{})},
	{Name: "PackedTransaction", ReflectType: reflect.TypeOf(PackedTransactionMessage{})},
}

var UnknownMessageTypeError = errors.New("unknown type")

func NewMessageType(aType byte) (t P2PMessageType, err error) {

	t = P2PMessageType(aType)
	if !t.isValid() {
		err = UnknownMessageTypeError
		return
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

func (t P2PMessageType) Attributes() (MessageAttributes, bool) {
	index := byte(t)

	if !t.isValid() {
		return MessageAttributes{}, false
	}

	attr := messageAttributes[index]
	return attr, true
}

type P2PMessageEnvelope struct {
	Length     uint32         `json:"length"`
	Type       P2PMessageType `json:"type"`
	Payload    []byte         `json:"-"`
	P2PMessage *P2PMessage    `json:"message" eos:"-"`
}

func (p2pMsg P2PMessageEnvelope) AsMessage() (P2PMessage, error) {
	attr, ok := p2pMsg.Type.Attributes()
	if !ok {
		return nil, UnknownMessageTypeError
	}

	if attr.ReflectType == nil {
		return nil, errors.New("Missing reflect type ")
	}

	msg := reflect.New(attr.ReflectType)

	err := p2pMsg.DecodePayload(msg.Interface())
	if err != nil {
		return nil, err
	}

	return msg.Interface().(P2PMessage), err
}

func (p2pMsg P2PMessageEnvelope) DecodePayload(message interface{}) error {

	attr, ok := p2pMsg.Type.Attributes()
	if !ok {
		return UnknownMessageTypeError
	}

	if attr.ReflectType == nil {
		return errors.New("missing reflect type")
	}

	messageType := reflect.TypeOf(message).Elem()
	if messageType != attr.ReflectType {
		return fmt.Errorf("given message type [%s] to not match payload type [%s]", messageType.Name(), attr.ReflectType.Name())
	}

	r := bytes.NewReader(p2pMsg.Payload)
	//lr := &LoggerReader{
	//	Reader: r,
	//}

	return NewOldDecoder(r).Decode(message)

}

func (p2pMsg P2PMessageEnvelope) MarshalBinary() ([]byte, error) {
	l := len(p2pMsg.Payload) + 1
	data := make([]byte, l+4, l+4)
	binary.LittleEndian.PutUint32(data[0:4], uint32(l))
	data[4] = byte(p2pMsg.Type)
	copy(data[5:], p2pMsg.Payload)

	return data, nil
}

func (p2pMsg *P2PMessageEnvelope) UnmarshalBinaryRead(r io.Reader) (err error) {

	lengthBytes := make([]byte, 4, 4)
	_, err = r.Read(lengthBytes)
	if err != nil {
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

	//fmt.Printf("Length: [%s] Payload: [%s]\n", hex.EncodeToString(lengthBytes), hex.EncodeToString(payloadBytes[:int(math.Min(float64(1000), float64(len(payloadBytes))))]))
	messageType, err := NewMessageType(payloadBytes[0])
	if err != nil {
		return
	}

	//fmt.Println("Payload type: ", messageType)

	*p2pMsg = P2PMessageEnvelope{
		Length:  size,
		Type:    messageType,
		Payload: payloadBytes[1:],
	}

	return nil
}
