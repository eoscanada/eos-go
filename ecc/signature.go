package ecc

import (
	"encoding/json"
	"fmt"
)

const SignatureK1Prefix = "SIG_K1_"
const SignatureR1Prefix = "SIG_R1_"
const SignatureWAPrefix = "SIG_WA_"

var signatureDataSize = new(int)

func init() { *signatureDataSize = 65 }

type innerSignature interface {
	verify(content []byte, hash []byte, pubKey PublicKey) bool
	publicKey(content []byte, hash []byte) (out PublicKey, err error)
	string(content []byte) string
	signatureMaterialSize() *int
}

// Signature represents a signature for some hash
type Signature struct {
	Curve   CurveID
	Content []byte // the Compact signature as bytes

	inner innerSignature
}

func (s Signature) Verify(hash []byte, pubKey PublicKey) bool {
	return s.inner.verify(s.Content, hash, pubKey)
}

func (s Signature) PublicKey(hash []byte) (out PublicKey, err error) {
	return s.inner.publicKey(s.Content, hash)
}

func (s Signature) String() string {
	return s.inner.string(s.Content)
}

func (s Signature) Validate() error {
	if s.inner == nil {
		return fmt.Errorf("the inner signature structure must be present, was nil")
	}

	if s.Curve == CurveK1 || s.Curve == CurveR1 {
		size := s.inner.signatureMaterialSize()
		if size == nil {
			return fmt.Errorf("R1 & K1 signatures must have a fixed key material size")
		}

		if len(s.Content) != *size {
			return fmt.Errorf("signature data must have a length of %d, got %d", *size, len(s.Content))
		}
	}

	return nil
}

func NewSignatureFromData(data []byte) (Signature, error) {
	signature := Signature{
		Curve:   CurveID(data[0]), // 1 byte
		Content: data[1:],         // 65 bytes for K1 and R1, variable length for WA
	}

	switch signature.Curve {
	case CurveK1:
		signature.inner = &innerK1Signature{}
	case CurveR1:
		signature.inner = &innerR1Signature{}
	case CurveWA:
		signature.inner = &innerWASignature{}
	default:
		return Signature{}, fmt.Errorf("invalid curve  %q", signature.Curve)
	}

	return signature, signature.Validate()
}

func MustNewSignatureFromData(data []byte) Signature {
	sig, err := NewSignatureFromData(data)
	if err != nil {
		panic(err.Error())
	}

	return sig
}

func MustNewSignature(fromText string) Signature {
	signature, err := NewSignature(fromText)
	if err != nil {
		panic(fmt.Errorf("invalid signature string: %w", err))
	}

	return signature
}

func NewSignature(signature string) (out Signature, err error) {
	if len(signature) < 8 {
		return out, fmt.Errorf("invalid format")
	}

	// We had a for/loop using a map before, this a disavantadge. The ordering was
	// not constant so we were not optimizing for the fact that compat keys appears way more
	// often than all others.
	//
	// We now have an unrolled for/loop specially ordered so that the most occurring prefix
	// is checked first.

	prefix := signature[0:7]
	if prefix == SignatureK1Prefix {
		return newSignature(CurveK1, signature[7:], newInnerK1Signature)
	}

	if prefix == SignatureR1Prefix {
		return newSignature(CurveR1, signature[7:], newInnerR1Signature)
	}

	if prefix == SignatureWAPrefix {
		return newSignature(CurveWA, signature[7:], newInnerWASignature)
	}

	return out, fmt.Errorf("unknown prefix %q", prefix)
}

func newSignature(curveID CurveID, in string, innerFactory func() innerSignature) (out Signature, err error) {
	payload, err := decodeSignatureMaterial(in, curveID)
	if err != nil {
		return out, err
	}

	return Signature{Curve: curveID, Content: payload, inner: innerFactory()}, nil
}

func (s Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Signature) UnmarshalJSON(data []byte) (err error) {
	var sig string
	err = json.Unmarshal(data, &sig)
	if err != nil {
		return
	}

	*s, err = NewSignature(sig)

	return
}
