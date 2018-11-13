package ecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

// Signature represents a signature for some hash
type Signature struct {
	Curve   CurveID
	Content []byte // the Compact signature as bytes
}

// Verify checks the signature against the pubKey. `hash` is a sha256
// hash of the payload to verify.
func (s Signature) Verify(hash []byte, pubKey PublicKey) bool {

	// TODO: choose the S256 curve, based on s.Curve
	recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), s.Content, hash)
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

// PublicKey retrieves the public key, but requires the
// payload.. that's the way to validate the signature. Use Verify() if
// you only want to validate.
func (s Signature) PublicKey(hash []byte) (out PublicKey, err error) {

	var recoveredKey *btcec.PublicKey
	switch s.Curve {
	case CurveK1:
		recoveredKey, _, err = btcec.RecoverCompact(btcec.S256(), s.Content, hash)
	case CurveR1:
		curve := btcec.S256R1()

		recoveredKey, _, err = btcec.RecoverCompact(curve, s.Content, hash)
	default:
		return PublicKey{}, fmt.Errorf("invalid curve: %s", s.Curve)
	}

	if err != nil {
		return out, err
	}

	return PublicKey{
		Curve:   s.Curve,
		Content: recoveredKey.SerializeCompressed(),
	}, nil
}

func (s Signature) String() string {
	checksum := Ripemd160checksumHashCurve(s.Content, s.Curve)
	buf := append(s.Content[:], checksum...)
	return "SIG_" + s.Curve.StringPrefix() + base58.Encode(buf)
	//return "SIG_" + base58.Encode(buf)
	//return base58.Encode(buf)
}

func NewSignature(fromText string) (Signature, error) {
	if !strings.HasPrefix(fromText, "SIG_") {
		return Signature{}, fmt.Errorf("signature should start with SIG_")
	}
	if len(fromText) < 8 {
		return Signature{}, fmt.Errorf("invalid signature length")
	}

	fromText = fromText[4:] // remove the `SIG_` prefix

	var curveID CurveID
	var curvePrefix = fromText[:3]
	switch curvePrefix {
	case "K1_":
		curveID = CurveK1
	case "R1_":
		curveID = CurveR1
	default:
		return Signature{}, fmt.Errorf("invalid curve prefix %q", curvePrefix)
	}
	fromText = fromText[3:] // strip curve ID

	sigbytes := base58.Decode(fromText)

	content := sigbytes[:len(sigbytes)-4]
	checksum := sigbytes[len(sigbytes)-4:]
	verifyChecksum := Ripemd160checksumHashCurve(content, curveID)
	if !bytes.Equal(verifyChecksum, checksum) {
		return Signature{}, fmt.Errorf("signature checksum failed, found %x expected %x", verifyChecksum, checksum)
	}

	return Signature{Curve: curveID, Content: content}, nil
}

func (a Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Signature) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*a, err = NewSignature(s)

	return
}
