package eos

import (
	"encoding/hex"
	"testing"

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
			HexString:      "0f01000006000330907937f721e855bb1ff0b6143af8bd82cba550d1c20438b29dbb30b2a74170e244d756399dd08666bbbe2e44c0f0967b9dfaa02b33234a2a31e52319961cf2de529eb4bc72f5b7b28b86848e837ce30805ae57e81023189eabd77d516a12ac02a7e9f40f2a2d66d2199a8bc9983e30066f61cbf24230653bb74b9c91e8d7b972ab0000000000ea3055000000000000205c87c866b524d313866c8c0c036518ee4e898b133d570b41403e31bb072f909d1b8e40775e261e24ae787aeeb127e5f25725ab684f382f7f92a68ee4672e1fac010000010100010000000000ea30550000000000ea305501001326fd25e5ad35924d41b9611304b682dcfca0bf1154008f97268d43bfaba999b104",
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
		//	HexString:      "2100000002000000000000000000000000000000004016e1d216df26150000000000000000",
		//	ExpectedStruct: &TimeMessage{},
		//	ExpectedErr:    false,
		//},
	}

	for _, c := range cases {
		decoded, err := hex.DecodeString(c.HexString)
		assert.NoError(t, err)

		//s := string(decoded)
		//fmt.Println(s)

		var p2pMessage P2PMessageEnvelope
		assert.NoError(t, UnmarshalBinary(decoded, &p2pMessage), c.Name)

		assert.NoError(t, p2pMessage.DecodePayload(c.ExpectedStruct), c.Name)
		//fmt.Println(c.ExpectedStruct)
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
