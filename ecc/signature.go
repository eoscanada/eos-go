package ecc

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
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
	buf := s.sig.Serialize()
	checksum := ripemd160checksum(buf)
	buf = append(buf, checksum[:4]...)
	return "EOS" + base58.Encode(buf)
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
