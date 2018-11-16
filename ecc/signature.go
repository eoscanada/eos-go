package ecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

type innerSignature interface {
	verify(content []byte, hash []byte, pubKey PublicKey) bool
	publicKey(content []byte, hash []byte) (out PublicKey, err error)
	string(content []byte) string
}

// Signature represents a signature for some hash
type Signature struct {
	Curve   CurveID
	Content []byte // the Compact signature as bytes

	innerSignature innerSignature
}

func (s Signature) Verify(hash []byte, pubKey PublicKey) bool {
	return s.innerSignature.verify(s.Content, hash, pubKey)
}

func (s Signature) PublicKey(hash []byte) (out PublicKey, err error) {
	return s.innerSignature.publicKey(s.Content, hash)
}

func (s Signature) String() string {
	return s.innerSignature.string(s.Content)
}

func NewSignatureFromData(data []byte) (Signature, error) {
	if len(data) != 66 {
		return Signature{}, fmt.Errorf("data length of a signature should be 66, reveived %d", len(data))
	}

	signature := Signature{
		Curve:   CurveID(data[0]), // 1 byte
		Content: data[1:],         // 65 bytes
	}

	switch signature.Curve {
	case CurveK1:
		signature.innerSignature = &innerK1Signature{}
	case CurveR1:
		signature.innerSignature = &innerR1Signature{}
	default:
		return Signature{}, fmt.Errorf("invalid curve  %q", signature.Curve)
	}
	return signature, nil
}

func MustNewSignatureFromData(data []byte) Signature {
	sig, err := NewSignatureFromData(data)
	if err != nil {
		panic(err.Error())
	}
	return sig
}

func NewSignature(fromText string) (Signature, error) {
	if !strings.HasPrefix(fromText, "SIG_") {
		return Signature{}, fmt.Errorf("signature should start with SIG_")
	}
	if len(fromText) < 8 {
		return Signature{}, fmt.Errorf("invalid signature length")
	}

	fromText = fromText[4:] // remove the `SIG_` prefix

	var curvePrefix = fromText[:3]
	switch curvePrefix {
	case "K1_":

		fromText = fromText[3:] // strip curve ID

		sigbytes := base58.Decode(fromText)

		content := sigbytes[:len(sigbytes)-4]
		checksum := sigbytes[len(sigbytes)-4:]
		verifyChecksum := Ripemd160checksumHashCurve(content, CurveK1)
		if !bytes.Equal(verifyChecksum, checksum) {
			return Signature{}, fmt.Errorf("signature checksum failed, found %x expected %x", verifyChecksum, checksum)
		}

		return Signature{Curve: CurveK1, Content: content, innerSignature: &innerK1Signature{}}, nil

	case "R1_":

		fromText = fromText[3:] // strip R1_
		content := base58.Decode(fromText)
		//todo: stuff here

		return Signature{Curve: CurveR1, Content: content, innerSignature: &innerR1Signature{}}, nil

	default:
		return Signature{}, fmt.Errorf("invalid curve prefix %q", curvePrefix)
	}
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
