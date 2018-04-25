package eos

import (
	"encoding/hex"
	"testing"

	"encoding/json"
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
func TestP2PMessage_DecodePayload(t *testing.T) {

	type Case struct {
		Name           string
		HexString      string
		ExpectedStruct P2PMessage
		ExpectedErr    bool
	}

	cases := []Case{
		{
			Name:           "HandshakeMessage Good",
			HexString:      "660100000057630000000000000000000000000000000000000000000000000000000000000000dc5b6788499eddf38fac9e70e7572a2cd1db1572dce8172ccd311f2a5de7682d000000000000000000000000000000000000000000000000000000000000000000002c0a9a97b6a8281500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002c436861726c6573732d4d6163426f6f6b2d50726f2d322e6c6f63616c3a3239383736202d2064633562363738ec330300000333ecb5c24000b4a478edc5557936fd1887cf12b55ab1a49301e8cdd318f1ed330300000333ed4bd4c21253ce599b038242574e9fec8ac7eff2dbb75cc469618f789c056c696e75781922454f532043616e616461206b38732d64332d62696f7333220100",
			ExpectedStruct: &HandshakeMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "GoAwayMessage Good",
			HexString:      "",
			ExpectedStruct: &GoAwayMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "Time Good",
			HexString:      "210000000284226efdfa6e2815e80d7efdfa6e2815342b7efdfa6e28150000000000000000",
			ExpectedStruct: &TimeMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "NoticeMessage Good",
			HexString:      "330000000300000000000000000001000000123403000100033412ead7e5d4ebc2f2a47edbbbbde260f2f13a29668dc0366aa72a8e665d",
			ExpectedStruct: &NoticeMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "NoticeMessage Good",
			HexString:      "130000000302000000bd2400000002000000be24000000",
			ExpectedStruct: &NoticeMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "RequestMessage Good",
			HexString:      "1300000004000000000000000000010000000000000000",
			ExpectedStruct: &RequestMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "SyncRequest Good",
			HexString:      "0900000005ed33030011340300",
			ExpectedStruct: &SyncRequestMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "SignedBlockSummaryMessage Good",
			HexString:      "0f010000060003338e14a48fd745ae3ff9d9dfb3d149c04507f9fde929c67311e1f004eab21849e444b7d9a2e6e5730eef8166a57ef1cc6e905cc0944fd8cee5d0c2ec6ec9a53076333442d7a57ce7820187f8889111443a1466666eac707ece1249d91a1a279ee04a856013e6b1fc7a14d74d8c237d47d8e2d951c422301f77342eb6f79f143958c90000000000ea30550000000000001f7ed3ef0a9a8bc7b695a390b20909852a7c43e3e54dce9d64aa107762a0b54b904b611084d16fb2bb4973374f380bd43548ac8e64e52c9911bcbb96128ccf5a11010000010100010000000000ea30550000000000ea305501001326b7530160d3c85db322fd2c550248fe22151d1e9423d2250a699595e5b1c1517d",
			ExpectedStruct: &SignedBlockSummaryMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "SignedBlockMessage Good",
			HexString:      "100100000700022373115ab1ff9296887bd50c09a767ec55e9af84773cfd3981c7771d66c67c20e1449dca75777f692ada89a117e359c078e13ff5c237b2d11e35b2d84b01eb92c712547d5ceae4a7bbc4ba2e17d86f8df6b3f6f336e6dc884eea726e1f6c68d72dd7b6454c4f04d4610ff2852993d0292d6022e2e3bba07c25d7a369276e01f1c3ed0000000000ea30550000000000002074f4a360736226aa1e0a14e7a11df65197c847e1054835d4229d806fb109b2f065791a1b9bcbfc2eefddab7bccc3e284cdacf28a2550072fc14b57c0f0bc7bc4010000010100010000000000ea30550000000000ea3055010013261ddcaa56f9d00f2a55ebbd28d5607c977e8e513554483891a8d4c0532c4b662e00",
			ExpectedStruct: &SignedBlockMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "SignedBlockMessage Good",
			HexString:      "0001000007000002b6fe7206df847da1e7e7691bd4ad25d0998f4fe226a2fbbd901221657d8988e64451b3b73face644d29603d2dcebf00af121c4af5e4ed86f5e80844fb4832dbeb03394ae8dde1ee9b0df4edbcd4cd7478768e97ca72801699ac07e9614b44707172ed846471b8f9d46914be0822d4d43dc2eebe08f759af1b6261af05107d23b880000000000ea305500000000000020783d2ff33f718c6cef10d5d0875ba99177259e184ae7ac92022f09302693b8925ed5429d7ec474cea6dec0fe4f5ab2e3da55188515c6a914e6e97a30d82ee00001000001010000010002268149cec70c1ed3fafc6f3ebe90e4c7af423d9b1edd627a4fe1a8f54d9c90f81e00",
			ExpectedStruct: &SignedBlockMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "SignedBlockMessage Good",
			HexString:      "0704000007000000c79e6ceeb00129996a74111e208a900265b728e9460c0d5ca119442847f1e0e64478156541313bfbd8572fb813dc29ee817cc3ae8d2d3fc5db3f32030cee969a0d760e71f0fe0a74ec82def3b746e7fa4f80e9e5a9b85eba9357420230e8f781ecefd023d49c218dea37902c82cb48fe8cb6d5ff776f2d4ebf04a2e9372dc53d2d0000000000ea30550000000000002040bc8a2c87c4af36b2166dc295babbb76bc55322ca0bac6735e968d1f7d714010fc77907c9240a64210c30476ef775426bcccf4f908fcaf3684a07a1b09e9138010000020100000100022619295edbc20ec236289bde12b51b2b2a99e1d662f355844acdc3390b4af410630100010000000000ea30550040cbda00ea3055010067674a094cb542b4314c425f5d5b8948ab1b1b329c09f7d8a1b4e645cdfe2d3d467d0101001f6dd53f77a5120d6c6dce756ca52c1c1c999e9895322e5172625a8149102d7b546fde27b6146bafefb86cf2e2ed81bc1fc9b1bcf3aaf4c236f35d4f49c149e4810000890516b4e05a0000c7000129996a00000000040000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea30550000735802ea305501000000010003fb47adddc6f5c4d7ad3c864771957aa9494860259666e61d6aeac397ca7430f401000001000000010003fb47adddc6f5c4d7ad3c864771957aa9494860259666e61d6aeac397ca7430f40100000100000000010000000000ea305500000000a8ed323201000000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500a6823403ea305501000000010003fb47adddc6f5c4d7ad3c864771957aa9494860259666e61d6aeac397ca7430f401000001000000010003fb47adddc6f5c4d7ad3c864771957aa9494860259666e61d6aeac397ca7430f40100000100000000010000000000ea305500000000a8ed323201000000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea30557015d289de8aa23b01000000010003fb47adddc6f5c4d7ad3c864771957aa9494860259666e61d6aeac397ca7430f401000001000000010003fb47adddc6f5c4d7ad3c864771957aa9494860259666e61d6aeac397ca7430f40100000100000000010000000000ea305500000000a8ed323201000000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000039ab18dd4101000000010003fb47adddc6f5c4d7ad3c864771957aa9494860259666e61d6aeac397ca7430f401000001000000010003fb47adddc6f5c4d7ad3c864771957aa9494860259666e61d6aeac397ca7430f40100000100000000010000000000ea305500000000a8ed32320100",
			ExpectedStruct: &SignedBlockMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "SignedTransactionMessage Good",
			HexString:      "",
			ExpectedStruct: &SignedTransactionMessage{},
			ExpectedErr:    false,
		},
		{
			Name:           "PackedTransactionMessage Good",
			HexString:      "9b000000090100204ccbd87f27e40f790c97297189657f323ea552664735571068f36140aebc548526072f7e1d543aa18260d8aef4956119165c7416da6d455325e8b7b29f592d87000054d561e05a0000743596958d05000000000100a6823403ea3055000000572d3ccdcd010000000000ea305500000000a8ed3232210000000000ea305500000039ab18dd41a08601000000000004454f530000000000",
			ExpectedStruct: &PackedTransactionMessage{},
			ExpectedErr:    false,
		},
	}

	for _, c := range cases {
		decoded, err := hex.DecodeString(c.HexString)
		assert.NoError(t, err)

		var p2PMessageEnvelope P2PMessageEnvelope
		assert.NoError(t, UnmarshalBinary(decoded, &p2PMessageEnvelope), c.Name)
		fmt.Println("Payload length: ", p2PMessageEnvelope.Length)
		assert.NoError(t, p2PMessageEnvelope.DecodePayload(c.ExpectedStruct), c.Name)

		data, err := json.Marshal(c.ExpectedStruct)
		assert.NoError(t, err)
		fmt.Println("JSON : ", string(data))
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
