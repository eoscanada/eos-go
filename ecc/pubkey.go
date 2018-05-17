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
const PublicKeyPrefixCompat = "EOS"

type PublicKey struct {
	Curve   CurveID
	Content []byte
}

func NewPublicKey(pubKey string) (out PublicKey, err error) {
	if len(pubKey) < 8 {
		return out, fmt.Errorf("invalid format")
	}

	var pubKeyMaterial string
	var curveID CurveID
	if strings.HasPrefix(pubKey, PublicKeyPrefix) {
		pubKeyMaterial = pubKey[len(PublicKeyPrefix):] // strip "PUB_"

		curvePrefix := pubKeyMaterial[:3]
		switch curvePrefix {
		case "K1_":
			curveID = CurveK1
		case "R1_":
			curveID = CurveR1
		default:
			return out, fmt.Errorf("unsupported curve prefix %q", curvePrefix)
		}
		pubKeyMaterial = pubKeyMaterial[3:] // strip "K1_"

	} else if strings.HasPrefix(pubKey, PublicKeyPrefixCompat) { // "EOS"
		pubKeyMaterial = pubKey[len(PublicKeyPrefixCompat):] // strip "EOS"
		curveID = CurveK1

	} else {
		return out, fmt.Errorf("public key should start with %q (or the old %q)", PublicKeyPrefix, PublicKeyPrefixCompat)
	}

	pubDecoded, err := checkDecode(pubKeyMaterial, curveID)
	if err != nil {
		return out, fmt.Errorf("checkDecode: %s", err)
	}

	return PublicKey{Curve: curveID, Content: pubDecoded}, nil
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

	// if curve != CurveK1 {
	// 	_, _ = h.Write([]byte(curve.String())) // conditionally ?
	// }
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
	// TODO: implement the curve switch according to `p.Curve`
	key, err := btcec.ParsePubKey(p.Content, btcec.S256())
	if err != nil {
		return nil, fmt.Errorf("parsePubKey: %s", err)
	}

	return key, nil
}

func (p PublicKey) String() string {
	//hash := ripemd160checksum(append([]byte{byte(p.Curve)}, p.Content...))  does the checksum include the curve ID?!
	hash := ripemd160checksum(p.Content, p.Curve)
	rawkey := append(p.Content, hash[:4]...)
	return PublicKeyPrefixCompat + base58.Encode(rawkey)
	// FIXME: when we decide to go ahead with the new representation.
	//return PublicKeyPrefix + p.Curve.StringPrefix() + base58.Encode(rawkey)
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
