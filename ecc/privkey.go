package ecc

import (
	cryptorand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
	"github.com/eoscanada/eos-go/btcsuite/btcutil"
)

const PrivateKeyPrefix = "PVT_"

func NewRandomPrivateKey() (*PrivateKey, error) {
	return newRandomPrivateKey(cryptorand.Reader)
}

func NewDeterministicPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	return newRandomPrivateKey(randSource)
}

func newRandomPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	rawPrivKey := make([]byte, 32)
	written, err := io.ReadFull(randSource, rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("error feeding crypto-rand numbers to seed ephemeral private key: %s", err)
	}
	if written != 32 {
		return nil, fmt.Errorf("couldn't write 32 bytes of randomness to seed ephemeral private key")
	}

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), rawPrivKey)

	inner := &innerK1PrivateKey{privKey: privKey}
	return &PrivateKey{Curve: CurveK1, inner: inner}, nil
}

func NewPrivateKey(wif string) (*PrivateKey, error) {
	// Strip potential prefix, and set curve
	var privKeyMaterial string
	if strings.HasPrefix(wif, PrivateKeyPrefix) { // "PVT_"
		privKeyMaterial = wif[len(PrivateKeyPrefix):]

		curvePrefix := privKeyMaterial[:3]
		privKeyMaterial = privKeyMaterial[3:] // remove "K1_"...

		switch curvePrefix {
		case "K1_":

			wifObj, err := btcutil.DecodeWIF(privKeyMaterial)
			if err != nil {
				return nil, err
			}
			inner := &innerK1PrivateKey{privKey: wifObj.PrivKey}
			return &PrivateKey{Curve: CurveK1, inner: inner}, nil
		case "R1_":

			inner := &innerR1PrivateKey{}
			return &PrivateKey{Curve: CurveR1, inner: inner}, nil

		default:
			return nil, fmt.Errorf("unsupported curve prefix %q", curvePrefix)
		}

	} else { // no-prefix, like before

		wifObj, err := btcutil.DecodeWIF(wif)
		if err != nil {
			return nil, err
		}
		inner := &innerK1PrivateKey{privKey: wifObj.PrivKey}
		return &PrivateKey{Curve: CurveK1, inner: inner}, nil
	}
}

type innerPrivateKey interface {
	publicKey() PublicKey
	sign(hash []byte) (out Signature, err error)
	string() string
}

type PrivateKey struct {
	Curve CurveID

	inner innerPrivateKey
}

func (p *PrivateKey) PublicKey() PublicKey {
	return p.inner.publicKey()
}

// Sign signs a 32 bytes SHA256 hash..
func (p *PrivateKey) Sign(hash []byte) (out Signature, err error) {
	return p.inner.sign(hash)
}

func (p *PrivateKey) String() string {
	return p.inner.string()
}

func (p *PrivateKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *PrivateKey) UnmarshalJSON(v []byte) (err error) {
	var s string
	if err = json.Unmarshal(v, &s); err != nil {
		return
	}

	newPrivKey, err := NewPrivateKey(s)
	if err != nil {
		return
	}

	*p = *newPrivKey

	return
}
