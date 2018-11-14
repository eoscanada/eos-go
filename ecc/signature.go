package ecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

type InnerSignature interface {
	Verify(content []byte, hash []byte, pubKey PublicKey) bool
	PublicKey(content []byte, hash []byte) (out PublicKey, err error)
}

// Signature represents a signature for some hash
type Signature struct {
	Curve   CurveID
	Content []byte // the Compact signature as bytes

	innerSignature InnerSignature
}

func (s Signature) Verify(hash []byte, pubKey PublicKey) bool {
	return s.innerSignature.Verify(s.Content, hash, pubKey)
}

func (s Signature) PublicKey(hash []byte) (out PublicKey, err error) {
	return s.innerSignature.PublicKey(s.Content, hash)
}

func (s Signature) String() string {
	checksum := Ripemd160checksumHashCurve(s.Content, s.Curve)
	buf := append(s.Content[:], checksum...)
	return "SIG_" + s.Curve.StringPrefix() + base58.Encode(buf)
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

		return Signature{Curve: CurveK1, Content: content, innerSignature: &InnerK1Signature{}}, nil

	case "R1_":

		// ICI!

		return Signature{Curve: CurveK1, Content: nil, innerSignature: &InnerK1Signature{}}, fmt.Errorf("invalid curve prefix %q", curvePrefix)
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
