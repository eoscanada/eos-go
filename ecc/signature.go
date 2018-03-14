package ecc

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
)

type Signature struct {
	sig *btcec.Signature
}

func (s *Signature) Verify(hash []byte, pubKey *PublicKey) bool {
	return s.sig.Verify(hash, pubKey.pubKey)
}

func (s *Signature) IsEqual(otherSig *Signature) bool {
	return s.sig.IsEqual(otherSig.sig)
}

func (s *Signature) String() string {
	return hex.EncodeToString(s.sig.Serialize())
}

func (s *Signature) Serialize() []byte {
	return s.sig.Serialize()
}

func ParseSignature(sigStr []byte) (*Signature, error) {
	sig, err := btcec.ParseSignature(sigStr, btcec.S256())
	if err != nil {
		return nil, err
	}

	return &Signature{sig: sig}, nil
}

func ParseDERSignature(sigStr []byte) (*Signature, error) {
	sig, err := btcec.ParseDERSignature(sigStr, btcec.S256())
	if err != nil {
		return nil, err
	}

	return &Signature{sig: sig}, nil
}
