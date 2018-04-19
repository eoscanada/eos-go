package p2p

import (
	"fmt"

	"github.com/pkg/errors"
)

type MessageType byte

const (
	HandshakeMessage MessageType = iota
	GoAwayMessage
	TimeMessage
	NoticeMessage
	RequestMessage
	SyncRequestMessage
	SignedBlockSummaryMessage
	SignedBlockMessage
	SignedTransactionMessage
	PackedTransactionMessage
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

func NewMessageType(aType byte) (t MessageType, err error) {

	t = MessageType(aType)
	if !t.isValid() {
		err = errors.New(fmt.Sprintf("unknown type [%d]", aType))
		return
	}

	return
}

func (t MessageType) isValid() bool {

	index := byte(t)
	return int(index) < len(messageNames) && index >= 0

}

func (t MessageType) Name() (string, bool) {

	index := byte(t)

	if !t.isValid() {
		return "unknown", false
	}

	name := messageNames[index]
	return name, true
}
