package eos

import (
	"encoding/hex"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestP2PMessage_UnmarshalBinaryRead(t *testing.T) {

	hexString := `09000000050100000019000000`
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		t.Error(err)
	}

	var s P2PMessageEnvelope

	assert.NoError(t, UnmarshalBinary(decoded, &s))
	assert.Equal(t, uint32(9), s.Length)
	assert.Equal(t, P2PMessageType(5), s.Type)
	assert.Equal(t, []byte{0x1, 0x0, 0x0, 0x0, 0x19, 0x0, 0x0, 0x0}, s.Payload)
}

/*
0f010000
06
0003338e14a48fd745ae3ff9d9dfb3d149c04507f9fde929c67311e1f004eab21849e444b7d9a2e6e5730eef8166a57ef1cc6e905cc0944fd8cee5d0c2ec6ec9a53076333442d7a57ce7820187f8889111443a1466666eac707ece1249d91a1a279ee04a856013e6b1fc7a14d74d8c237d47d8e2d951c422301f77342eb6f79f143958c90000000000ea30550000000000001f7ed3ef0a9a8bc7b695a390b20909852a7c43e3e54dce9d64aa107762a0b54b904b611084d16fb2bb4973374f380bd43548ac8e64e52c9911bcbb96128ccf5a11010000010100010000000000ea30550000000000ea305501001326b7530160d3c85db322fd2c550248fe22151d1e9423d2250a699595e5b1c1517d",
*/
func TestP2PMessage_DecodePayload(t *testing.T) {

	type Case struct {
		Name           string
		HexString      string
		ExpectedStruct P2PMessage
		ExpectedErr    bool
	}

	cases := []Case{
		{
			Name:           "SignedBlockSummaryMessage Good",
			HexString:      "0f010000060003338e14a48fd745ae3ff9d9dfb3d149c04507f9fde929c67311e1f004eab21849e444b7d9a2e6e5730eef8166a57ef1cc6e905cc0944fd8cee5d0c2ec6ec9a53076333442d7a57ce7820187f8889111443a1466666eac707ece1249d91a1a279ee04a856013e6b1fc7a14d74d8c237d47d8e2d951c422301f77342eb6f79f143958c90000000000ea30550000000000001f7ed3ef0a9a8bc7b695a390b20909852a7c43e3e54dce9d64aa107762a0b54b904b611084d16fb2bb4973374f380bd43548ac8e64e52c9911bcbb96128ccf5a11010000010100010000000000ea30550000000000ea305501001326b7530160d3c85db322fd2c550248fe22151d1e9423d2250a699595e5b1c1517d",
			ExpectedStruct: &SignedBlockSummaryMessage{},
			ExpectedErr:    false,
		},
		//{
		//	Name:           "SignedBlockMessage Good",
		//	HexString:      "100100000700022373115ab1ff9296887bd50c09a767ec55e9af84773cfd3981c7771d66c67c20e1449dca75777f692ada89a117e359c078e13ff5c237b2d11e35b2d84b01eb92c712547d5ceae4a7bbc4ba2e17d86f8df6b3f6f336e6dc884eea726e1f6c68d72dd7b6454c4f04d4610ff2852993d0292d6022e2e3bba07c25d7a369276e01f1c3ed0000000000ea30550000000000002074f4a360736226aa1e0a14e7a11df65197c847e1054835d4229d806fb109b2f065791a1b9bcbfc2eefddab7bccc3e284cdacf28a2550072fc14b57c0f0bc7bc4010000010100010000000000ea30550000000000ea3055010013261ddcaa56f9d00f2a55ebbd28d5607c977e8e513554483891a8d4c0532c4b662e00",
		//	ExpectedStruct: &SignedBlockMessage{},
		//	ExpectedErr:    false,
		//},
		//{
		//	Name:           "Time Good",
		//	HexString:      "210000000284226efdfa6e2815e80d7efdfa6e2815342b7efdfa6e28150000000000000000",
		//	ExpectedStruct: &TimeMessage{},
		//	ExpectedErr:    false,
		//},
	}

	for _, c := range cases {
		decoded, err := hex.DecodeString(c.HexString)
		assert.NoError(t, err)

		//s := string(decoded)
		//fmt.Println(s)

		var p2PMessageEnvelope P2PMessageEnvelope
		assert.NoError(t, UnmarshalBinary(decoded, &p2PMessageEnvelope), c.Name)

		assert.NoError(t, p2PMessageEnvelope.DecodePayload(c.ExpectedStruct), c.Name)
		fmt.Println("BILC: ", c.ExpectedStruct)
	}

	//todo : more assert

}
func TestP2PMessage_AsMessage(t *testing.T) {

	hexString := `2100000002000000000000000000000000000000004016e1d216df26150000000000000000`
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		t.Error(err)
	}
	var p2pMessage P2PMessageEnvelope
	assert.NoError(t, UnmarshalBinary(decoded, &p2pMessage))

	msg, err := p2pMessage.AsMessage()

	assert.NoError(t, err)
	assert.IsType(t, &TimeMessage{}, msg)

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
