package ecc

import (
	"bytes"
	"fmt"
)

type innerR1PrivateKey struct {
}

func (k *innerR1PrivateKey) publicKey() PublicKey {
	var pubKeyData []byte
	pubKeyData = append(pubKeyData, byte(1))
	pubKeyData = append(pubKeyData, bytes.Repeat([]byte{0}, 33)...)
	return PublicKey{Curve: CurveK1, Content: pubKeyData, inner: &innerK1PublicKey{}}
}

func (k *innerR1PrivateKey) sign(hash []byte) (out Signature, err error) {
	return out, fmt.Errorf("R1 not supported")
}

func (k *innerR1PrivateKey) string() string {
	return "PVT_R1_PLACE_HOLDER"
}
