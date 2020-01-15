package ecc

import (
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

type innerWASignature struct {
}

func newInnerWASignature() innerSignature {
	return &innerWASignature{}
}

func (s innerWASignature) verify(content []byte, hash []byte, pubKey PublicKey) bool {
	// It seems from my understanding that WA uses standard ECDSA P256 algorithm, so we
	// should be able to verify signature of message against PublicKey.
	//
	// See https://thanethomson.com/2018/11/30/validating-ecdsa-signatures-golang/
	return false
}

func (s *innerWASignature) publicKey(content []byte, hash []byte) (out PublicKey, err error) {
	return out, fmt.Errorf("WA not supported")
}

func (s innerWASignature) string(content []byte) string {
	checksum := Ripemd160checksumHashCurve(content, CurveWA)
	buf := append(content[:], checksum...)
	return "SIG_WA_" + base58.Encode(buf)
}

func (s innerWASignature) signatureMaterialSize() *int {
	return nil
}
