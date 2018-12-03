package ecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const PublicKeyPrefix = "PUB_"
const PublicKeyK1Prefix = "PUB_K1_"
const PublicKeyR1Prefix = "PUB_R1_"
const PublicKeyPrefixCompat = "EOS"

type innerPublicKey interface {
	key(content []byte) (*btcec.PublicKey, error)
	prefix() string
}

type PublicKey struct {
	Curve   CurveID
	Content []byte

	inner innerPublicKey
}

func NewPublicKeyFromData(data []byte) (out PublicKey, err error) {
	if len(data) != 34 {
		return out, fmt.Errorf("public key data must have a length of 33 ")
	}

	out = PublicKey{
		Curve:   CurveID(data[0]), // 1 byte
		Content: data[1:],         // 33 bytes
	}

	switch out.Curve {
	case CurveK1:
		out.inner = &innerK1PublicKey{}
	case CurveR1:
		out.inner = &innerR1PublicKey{}
	default:
		return out, fmt.Errorf("unsupported curve prefix %q", out.Curve)
	}

	return out, nil
}

func MustNewPublicKeyFromData(data []byte) PublicKey {
	key, err := NewPublicKeyFromData(data)
	if err != nil {
		panic(err.Error())
	}
	return key
}

func NewPublicKey(pubKey string) (out PublicKey, err error) {
	if len(pubKey) < 8 {
		return out, fmt.Errorf("invalid format")
	}

	var decodedPubKey []byte
	var curveID CurveID
	var inner innerPublicKey

	if strings.HasPrefix(pubKey, PublicKeyR1Prefix) {
		pubKeyMaterial := pubKey[len(PublicKeyR1Prefix):] // strip "PUB_R1_"
		curveID = CurveR1
		decodedPubKey, err = checkDecode(pubKeyMaterial, curveID)
		if err != nil {
			return out, fmt.Errorf("checkDecode: %s", err)
		}
		inner = &innerR1PublicKey{}
	} else if strings.HasPrefix(pubKey, PublicKeyK1Prefix) {
		pubKeyMaterial := pubKey[len(PublicKeyK1Prefix):] // strip "PUB_K1_"
		curveID = CurveK1
		decodedPubKey, err = checkDecode(pubKeyMaterial, curveID)
		if err != nil {
			return out, fmt.Errorf("checkDecode: %s", err)
		}
		inner = &innerK1PublicKey{}
	} else if strings.HasPrefix(pubKey, PublicKeyPrefixCompat) { // "EOS"
		pubKeyMaterial := pubKey[len(PublicKeyPrefixCompat):] // strip "EOS"
		curveID = CurveK1
		decodedPubKey, err = checkDecode(pubKeyMaterial, curveID)
		if err != nil {
			return out, fmt.Errorf("checkDecode: %s", err)
		}
		inner = &innerK1PublicKey{}
	} else {
		return out, fmt.Errorf("public key should start with [%q | %q] (or the old %q)", PublicKeyK1Prefix, PublicKeyR1Prefix, PublicKeyPrefixCompat)
	}

	return PublicKey{Curve: curveID, Content: decodedPubKey, inner: inner}, nil
}

func MustNewPublicKey(pubKey string) PublicKey {
	key, err := NewPublicKey(pubKey)
	if err != nil {
		panic(err.Error())
	}
	return key
}

// CheckDecode decodes a string that was encoded with CheckEncode and verifies the checksum.
func checkDecode(input string, curve CurveID) (result []byte, err error) {
	decoded := base58.Decode(input)
	if len(decoded) < 5 {
		return nil, fmt.Errorf("invalid format")
	}
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	///// WARN: ok the ripemd160checksum should include the prefix in CERTAIN situations,
	// like when we imported the PubKey without a prefix ?! tied to the string representation
	// or something ? weird.. checksum shouldn't change based on the string reprsentation.
	if bytes.Compare(ripemd160checksum(decoded[:len(decoded)-4], curve), cksum[:]) != 0 {
		return nil, fmt.Errorf("invalid checksum")
	}
	// perhaps bitcoin has a leading net ID / version, but EOS doesn't
	payload := decoded[:len(decoded)-4]
	result = append(result, payload...)
	return
}

func ripemd160checksum(in []byte, curve CurveID) []byte {
	h := ripemd160.New()
	_, _ = h.Write(in) // this implementation has no error path

	if curve != CurveK1 {
		_, _ = h.Write([]byte(curve.String()))
	}

	sum := h.Sum(nil)
	return sum[:4]
}

func Ripemd160checksumHashCurve(in []byte, curve CurveID) []byte {
	h := ripemd160.New()
	_, _ = h.Write(in) // this implementation has no error path

	// FIXME: this seems to be only rolled out to the `SIG_` things..
	// proper support for importing `EOS` keys isn't rolled out into `dawn4`.
	_, _ = h.Write([]byte(curve.String())) // conditionally ?
	sum := h.Sum(nil)
	return sum[:4]
}

func (p PublicKey) Key() (*btcec.PublicKey, error) {
	return p.inner.key(p.Content)
}

func (p PublicKey) String() string {
	data := p.Content
	if len(data) == 0 {
		data = make([]byte, 33)
	}

	hash := ripemd160checksum(data, p.Curve)

	rawKey := make([]byte, 37)
	copy(rawKey, data[:33])
	copy(rawKey[33:], hash[:4])

	return p.inner.prefix() + base58.Encode(rawKey)
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	s := p.String()
	return json.Marshal(s)
}

func (p *PublicKey) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	newKey, err := NewPublicKey(s)
	if err != nil {
		return err
	}

	*p = newKey

	return nil
}
