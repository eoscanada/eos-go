package ecc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const PublicKeyPrefix = "PUB_"
const PublicKeyK1Prefix = "PUB_K1_"
const PublicKeyR1Prefix = "PUB_R1_"
const PublicKeyWAPrefix = "PUB_WA_"
const PublicKeyPrefixCompat = "EOS"

var publicKeyDataSize = new(int)

func init() { *publicKeyDataSize = 33 }

type innerPublicKey interface {
	key(content []byte) (*btcec.PublicKey, error)
	prefix() string
	keyMaterialSize() *int
}

type PublicKey struct {
	Curve   CurveID
	Content []byte

	inner innerPublicKey
}

func NewPublicKeyFromData(data []byte) (out PublicKey, err error) {
	if len(data) <= 0 {
		return out, errors.New("data must have at least one byte, got 0")
	}

	out = PublicKey{
		Curve:   CurveID(data[0]), // 1 byte
		Content: data[1:],         // 33 bytes for K1 & R1 keys, variable size for WA
	}

	switch out.Curve {
	case CurveK1:
		out.inner = &innerK1PublicKey{}
	case CurveR1:
		out.inner = &innerR1PublicKey{}
	case CurveWA:
		out.inner = &innerWAPublicKey{}
	default:
		return out, fmt.Errorf("unsupported curve prefix %q", out.Curve)
	}

	return out, out.Validate()
}

func MustNewPublicKeyFromData(data []byte) PublicKey {
	key, err := NewPublicKeyFromData(data)
	if err != nil {
		panic(err.Error())
	}
	return key
}

type pubkeyReaderManifest struct {
	curveID CurveID
	inner   func() innerPublicKey
}

var pubKeyReaderManifest = map[string]pubkeyReaderManifest{
	PublicKeyPrefixCompat: pubkeyReaderManifest{CurveK1, newInnerK1PublicKey},
	PublicKeyK1Prefix:     pubkeyReaderManifest{CurveK1, newInnerK1PublicKey},
	PublicKeyR1Prefix:     pubkeyReaderManifest{CurveR1, newInnerR1PublicKey},
	PublicKeyWAPrefix:     pubkeyReaderManifest{CurveWA, newInnerWAPublicKey},
}

func NewPublicKey(pubKey string) (out PublicKey, err error) {
	if len(pubKey) < 8 {
		return out, fmt.Errorf("invalid format")
	}

	for prefix, manifest := range pubKeyReaderManifest {
		if !strings.HasPrefix(pubKey, prefix) {
			continue
		}

		pubKeyMaterial := pubKey[len(prefix):]
		decodedPubKey, err := decodeKeyMaterial(pubKeyMaterial, manifest.curveID)
		if err != nil {
			return out, fmt.Errorf("checkDecode: %s", err)
		}

		return PublicKey{Curve: manifest.curveID, Content: decodedPubKey, inner: manifest.inner()}, nil
	}

	return out, fmt.Errorf("public key should start with [%q | %q | %q] (or the old %q)", PublicKeyK1Prefix, PublicKeyR1Prefix, PublicKeyWAPrefix, PublicKeyPrefixCompat)
}

func MustNewPublicKey(pubKey string) PublicKey {
	key, err := NewPublicKey(pubKey)
	if err != nil {
		panic(err.Error())
	}
	return key
}

func (p PublicKey) Validate() error {
	if p.inner == nil {
		return fmt.Errorf("the inner public key structure must be present, was nil")
	}

	if p.Curve == CurveK1 || p.Curve == CurveR1 {
		size := p.inner.keyMaterialSize()
		if size == nil {
			return fmt.Errorf("R1 & K1 public keys must have a fixed key material size")
		}

		if len(p.Content) != *size {
			return fmt.Errorf("public key data must have a length of %d, got %d", *size, len(p.Content))
		}
	}

	return nil
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

var emptyKeyMaterial = make([]byte, 33)

func (p PublicKey) String() string {
	data := p.Content
	if len(data) == 0 {
		// Nothing really to do, just output some garbage
		return p.inner.prefix() + base58.Encode(emptyKeyMaterial)
	}

	hash := ripemd160checksum(data, p.Curve)
	size := p.KeyMaterialSize()

	rawKey := make([]byte, size+4)
	copy(rawKey, data[:size])
	copy(rawKey[size:], hash[:4])

	return p.inner.prefix() + base58.Encode(rawKey)
}

func (p PublicKey) KeyMaterialSize() int {
	size := p.inner.keyMaterialSize()
	if size != nil {
		return *size
	}

	return len(p.Content)
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
