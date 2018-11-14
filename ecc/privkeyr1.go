package ecc

import (
	"fmt"
)

type InnerR1PrivateKey struct {
}

func (k *InnerR1PrivateKey) PublicKey() PublicKey {
	return PublicKey{Curve: CurveK1, Content: nil, inner: &InnerK1PublicKey{}}
}

func (p *InnerR1PrivateKey) Sign(hash []byte) (out Signature, err error) {
	if len(hash) != 32 {
		return out, fmt.Errorf("hash should be 32 bytes")
	}

	return Signature{Curve: CurveR1, Content: nil, innerSignature: &InnerR1Signature{}}, nil
}

func (k *InnerR1PrivateKey) String() string {

	return "a string here"
}
