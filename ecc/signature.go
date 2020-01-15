package ecc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
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

type sigReaderManifest struct {
	curveID CurveID
	inner   func() innerSignature
}

var sigReaderManifests = map[string]sigReaderManifest{
	SignatureK1Prefix: sigReaderManifest{CurveK1, newInnerK1Signature},
	SignatureR1Prefix: sigReaderManifest{CurveR1, newInnerR1Signature},
	SignatureWAPrefix: sigReaderManifest{CurveWA, newInnerWASignature},
}

func NewSignature(signature string) (out Signature, err error) {
	if len(signature) < 8 {
		return out, fmt.Errorf("invalid format")
	}

	prefix := signature[0:7]
	manifest, found := sigReaderManifests[prefix]
	if !found {
		return out, fmt.Errorf("unknown prefix %q", prefix)
	}

	fromText := signature[7:]
	decoder := keyMaterialDecoders[manifest.curveID]
	if decoder == nil {
		decoder = keyMaterialDecoderFunc(base58.Decode)
	}

	decoded := decoder.Decode(fromText)

	content := decoded[:len(decoded)-4]
	checksum := decoded[len(decoded)-4:]
	verifyChecksum := Ripemd160checksumHashCurve(content, manifest.curveID)
	if !bytes.Equal(verifyChecksum, checksum) {
		return Signature{}, fmt.Errorf("signature checksum failed, found %x expected %x", verifyChecksum, checksum)
	}

	return Signature{Curve: manifest.curveID, Content: content, inner: manifest.inner()}, nil
}

func MustNewSignature(fromText string) Signature {
	signature, err := NewSignature(fromText)
	if err != nil {
		panic(fmt.Errorf("invalid signature string: %s", err))
	}

	return signature
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
