package ecc

import (
	"bytes"
	"fmt"
)

type innerWAPrivateKey struct {
}

func (k *innerWAPrivateKey) publicKey() PublicKey {
	//TODO: Finish WA support here - for now we return bogus key
	var pubKeyData []byte
	pubKeyData = append(pubKeyData, 0x03)                           // ySign (either 0x02 or 0x03)
	pubKeyData = append(pubKeyData, bytes.Repeat([]byte{0}, 32)...) // X
	pubKeyData = append(pubKeyData, 0x02)                           // Flags Presence (2 Verified, 1 Present, 0 None)
	pubKeyData = append(pubKeyData, 0x00)                           // Relay Party ID (String encoded, i.e. varuint32 length + characters)

	return PublicKey{Curve: CurveWA, Content: pubKeyData, inner: &innerWAPublicKey{}}
}

func (k *innerWAPrivateKey) sign(hash []byte) (out Signature, err error) {
	return out, fmt.Errorf("WA not supported")
}

func (k *innerWAPrivateKey) string() string {
	return "PVT_WA_PLACE_HOLDER"
}
