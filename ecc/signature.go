package ecc

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
)

// Signature represents a signature for some hash
type Signature []byte

// Verify checks the signature against the pubKey. `hash` is a sha256
// hash of the payload to verify.
func (s Signature) Verify(payload []byte, pubKey *PublicKey) bool {
	hash := sha256.New()
	hash.Write(payload)

	recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), s, hash.Sum(nil))
	if err != nil {
		return false
	}
	if recoveredKey.IsEqual(pubKey.pubKey) {
		return true
	}
	return false
}

// PublicKey retrieves the public key, but requires the
// payload.. that's the way to validate the signature. Use Verify() if
// you only want to validate.
func (s Signature) PublicKey(payload []byte) (*PublicKey, error) {
	hash := sha256.New()
	hash.Write(payload)

	recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), s, hash.Sum(nil))
	if err != nil {
		return nil, err
	}

	return &PublicKey{recoveredKey}, err
}

func (s Signature) String() string {
	checksum := ripemd160checksum(s)
	buf := append(s[:], checksum[:4]...)
	return "EOS" + base58.Encode(buf)
}

func NewSignature(fromText string) (Signature, error) {
	sigbytes := base58.Decode(fromText[3:]) // simply remove the `EOS` in front..

	content := sigbytes[:len(sigbytes)-4]
	checksum := sigbytes[len(sigbytes)-4:]
	verifyChecksum := ripemd160checksum(content)
	if !bytes.Equal(verifyChecksum, checksum) {
		return nil, fmt.Errorf("signature checksum failed, found %x expected %x", verifyChecksum, checksum)
	}

	// TODO: validate the checksum..
	return Signature(content), nil
}

func (a Signature) MarshalBinary() ([]byte, error) {
	return append(bytes.Repeat([]byte{0}, 66-len(a)), a...), nil
}

func (a *Signature) UnmarshalBinary(data []byte) error {
	*a = Signature(data)
	return nil
}

func (a *Signature) UnmarshalBinarySize() int { return 66 }

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
