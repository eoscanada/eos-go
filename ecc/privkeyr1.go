package ecc

import (
	"fmt"
)

type innerR1PrivateKey struct {
}

func (k *innerR1PrivateKey) publicKey() PublicKey {
	return PublicKey{Curve: CurveK1, Content: nil, inner: &innerK1PublicKey{}}
}

func (p *innerR1PrivateKey) sign(hash []byte) (out Signature, err error) {
	if len(hash) != 32 {
		return out, fmt.Errorf("hash should be 32 bytes")
	}

	return Signature{Curve: CurveR1, Content: nil, innerSignature: &innerR1Signature{}}, nil
}

func (k *innerR1PrivateKey) string() string {

	return "a string here"
}
