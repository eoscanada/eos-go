package ecc

import (
	"bytes"
	"fmt"
)

type innerWAPrivateKey struct {
}

func (k *innerWAPrivateKey) publicKey() PublicKey {
	//TODO: Finish WA support here - for now we  return bogus key
	var pubKeyData []byte
	pubKeyData = append(pubKeyData, byte(1))
	pubKeyData = append(pubKeyData, bytes.Repeat([]byte{0}, 33)...)
	return PublicKey{Curve: CurveWA, Content: pubKeyData, inner: &innerWAPublicKey{}}
}

func (k *innerWAPrivateKey) sign(hash []byte) (out Signature, err error) {
	return out, fmt.Errorf("WA not supported")
}

func (k *innerWAPrivateKey) string() string {
	return "PVT_WA_PLACE_HOLDER"
}
