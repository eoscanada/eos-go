package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {

	type Case struct {
		Type         MessageType
		ExpectedName interface{}
		OK           bool
	}

	cases := []Case{
		{Type: HandshakeMessage, ExpectedName: "Handshake", OK: true},
		{Type: GoAwayMessage, ExpectedName: "GoAway", OK: true},
		{Type: TimeMessage, ExpectedName: "Time", OK: true},
		{Type: NoticeMessage, ExpectedName: "Notice", OK: true},
		{Type: RequestMessage, ExpectedName: "Request", OK: true},
		{Type: SyncRequestMessage, ExpectedName: "SyncRequest", OK: true},
		{Type: SignedBlockSummaryMessage, ExpectedName: "SignedBlockSummary", OK: true},
		{Type: SignedBlockMessage, ExpectedName: "SignedBlock", OK: true},
		{Type: SignedTransactionMessage, ExpectedName: "SignedTransaction", OK: true},
		{Type: PackedTransactionMessage, ExpectedName: "PackedTransaction", OK: true},
		{Type: MessageType(100), ExpectedName: nil, OK: false},
	}

	for _, c := range cases {

		name, ok := c.Type.Name()
		assert.Equal(t, c.OK, ok)
		if ok {
			assert.Equal(t, c.ExpectedName, name)
		}
	}
}
