package ecc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/golangcrypto/ripemd160"
)

const PublicKeyPrefix = "EOS"

type PublicKey struct {
	pubKey *btcec.PublicKey
}

func NewPublicKey(pubKey string) (*PublicKey, error) {
	if len(pubKey) < 5 {
		return nil, fmt.Errorf("invalid format")
	}

	if !strings.HasPrefix(pubKey, PublicKeyPrefix) {
		return nil, fmt.Errorf("public key should start with %q", PublicKeyPrefix)
	}

	pubdecoded, version, err := checkDecode(pubKey[3:])
	if err != nil {
		return nil, fmt.Errorf("checkDecode: %s", err)
	}
	if version != 2 {
		return nil, fmt.Errorf("invalid version")
	}

	key, err := btcec.ParsePubKey(pubdecoded, btcec.S256())
	if err != nil {
		return nil, fmt.Errorf("parsePubKey: %s", err)
	}

	return &PublicKey{pubKey: key}, nil
}

// CheckDecode decodes a string that was encoded with CheckEncode and verifies the checksum.
func checkDecode(input string) (result []byte, version byte, err error) {
	decoded := base58.Decode(input)
	if len(decoded) < 5 {
		return nil, 0, fmt.Errorf("invalid format")
	}
	version = decoded[0]
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if bytes.Compare(ripemd160checksum(decoded[:len(decoded)-4]), cksum[:]) != 0 {
		return nil, 0, fmt.Errorf("invalid checksum")
	}
	// perhaps bitcoin has a leading net ID, but EOS doesn't
	payload := decoded[0 : len(decoded)-4]
	result = append(result, payload...)
	return
}

func ripemd160checksum(in []byte) []byte {
	h := ripemd160.New()
	_, _ = h.Write(in) // this implementation has no error path
	return h.Sum(nil)[:4]
}

func (p *PublicKey) String() string {
	rawkey := p.pubKey.SerializeCompressed()
	hash := ripemd160checksum(rawkey)
	rawkey = append(rawkey, hash[:4]...)
	return PublicKeyPrefix + base58.Encode(rawkey)
}

func (p *PublicKey) ToHex() string {
	return hex.EncodeToString(p.pubKey.SerializeCompressed())
}
