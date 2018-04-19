package eos

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {

	hexString := `09000000050100000019000000`
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		t.Error(err)
	}

	var s P2PMessage

	assert.NoError(t, UnmarshalBinary(decoded, &s))
	assert.Equal(t, uint32(9), s.Length)
	assert.Equal(t, P2PMessageType(5), s.Type)
	assert.Equal(t, []byte{0x1, 0x0, 0x0, 0x0, 0x19, 0x0, 0x0, 0x0}, s.Payload)
}

func TestMessageType_Name(t *testing.T) {

	type Case struct {
		Type         P2PMessageType
		ExpectedName interface{}
		OK           bool
	}

	cases := []Case{
		{Type: HandshakeMessageType, ExpectedName: "Handshake", OK: true},
		{Type: GoAwayMessageType, ExpectedName: "GoAway", OK: true},
		{Type: TimeMessageType, ExpectedName: "Time", OK: true},
		{Type: NoticeMessageType, ExpectedName: "Notice", OK: true},
		{Type: RequestMessageType, ExpectedName: "Request", OK: true},
		{Type: SyncRequestMessageType, ExpectedName: "SyncRequest", OK: true},
		{Type: SignedBlockSummaryMessageType, ExpectedName: "SignedBlockSummary", OK: true},
		{Type: SignedBlockMessageType, ExpectedName: "SignedBlock", OK: true},
		{Type: SignedTransactionMessageType, ExpectedName: "SignedTransaction", OK: true},
		{Type: PackedTransactionMessageType, ExpectedName: "PackedTransaction", OK: true},
		{Type: P2PMessageType(100), ExpectedName: "Unknown", OK: false},
	}

	for _, c := range cases {

		name, ok := c.Type.Name()
		assert.Equal(t, c.OK, ok)
		assert.Equal(t, c.ExpectedName, name)
	}
}
