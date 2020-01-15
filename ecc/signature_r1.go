package ecc

import (
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

type innerR1Signature struct {
}

func newInnerR1Signature() innerSignature {
	return &innerR1Signature{}
}

func (s innerR1Signature) verify(content []byte, hash []byte, pubKey PublicKey) bool {
	return false
}

func (s *innerR1Signature) publicKey(content []byte, hash []byte) (out PublicKey, err error) {
	return out, fmt.Errorf("R1 not supported")
}

func (s innerR1Signature) string(content []byte) string {
	checksum := Ripemd160checksumHashCurve(content, CurveR1)
	buf := append(content[:], checksum...)
	return "SIG_R1_" + base58.Encode(buf)
}

func (s innerR1Signature) signatureMaterialSize() *int {
	return signatureDataSize
}
