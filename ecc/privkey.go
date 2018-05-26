package ecc

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
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

	h := sha256.New()
	h.Write(rawPrivKey)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), h.Sum(nil))

	return &PrivateKey{Curve: CurveK1, privKey: privKey}, nil
}

func NewPrivateKey(wif string) (*PrivateKey, error) {
	// Strip potential prefix, and set curve
	var privKeyMaterial string
	var curveID CurveID
	if strings.HasPrefix(wif, PrivateKeyPrefix) { // "PVT_"
		privKeyMaterial = wif[len(PrivateKeyPrefix):]

		// check the subcurve
		curvePrefix := privKeyMaterial[:3]
		switch curvePrefix {
		case "K1_":
			curveID = CurveK1
		case "R1_":
			curveID = CurveR1
		default:
			return nil, fmt.Errorf("unsupported curve prefix %q", curvePrefix)
		}

		privKeyMaterial = privKeyMaterial[3:] // remove "K1_"...

	} else { // no-prefix, like before
		privKeyMaterial = wif
		curveID = CurveK1
	}

	wifObj, err := btcutil.DecodeWIF(privKeyMaterial)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{Curve: curveID, privKey: wifObj.PrivKey}, nil
}

type PrivateKey struct {
	Curve   CurveID
	privKey *btcec.PrivateKey
}

func (p *PrivateKey) PublicKey() PublicKey {
	return PublicKey{Curve: p.Curve, Content: p.privKey.PubKey().SerializeCompressed()}
}

// Sign signs a 32 bytes SHA256 hash..
func (p *PrivateKey) Sign(hash []byte) (out Signature, err error) {
	if len(hash) != 32 {
		return out, fmt.Errorf("hash should be 32 bytes")
	}

	if p.Curve != CurveK1 {
		return out, fmt.Errorf("curve R1 not supported for signature")
	}

	// TODO: implement the R1 curve..
	compactSig, err := p.privKey.SignCanonical(btcec.S256(), hash)
	if err != nil {
		return out, fmt.Errorf("canonical, %s", err)
	}

	return Signature{Curve: p.Curve, Content: compactSig}, nil
}

func (p *PrivateKey) String() string {
	wif, _ := btcutil.NewWIF(p.privKey, '\x80', false) // no error possible
	return wif.String()
	// FIXME: when we decide to go ahead with the new representation.
	//return PrivateKeyPrefix + p.Curve.StringPrefix() + wif.String()
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
