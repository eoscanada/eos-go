package ecc

import (
	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
)

type InnerK1Signature struct {
}

// Verify checks the signature against the pubKey. `hash` is a sha256
// hash of the payload to verify.
func (s *InnerK1Signature) Verify(content []byte, hash []byte, pubKey PublicKey) bool {
	recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), content, hash)
	if err != nil {
		return false
	}
	key, err := pubKey.Key()
	if err != nil {
		return false
	}
	if recoveredKey.IsEqual(key) {
		return true
	}
	return false
}

func (s *InnerK1Signature) PublicKey(content []byte, hash []byte) (out PublicKey, err error) {

	recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), content, hash)

	if err != nil {
		return out, err
	}

	return PublicKey{
		Curve:   CurveK1,
		Content: recoveredKey.SerializeCompressed(),
	}, nil
}
