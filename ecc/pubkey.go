package ecc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eoscanada/eos-go/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/golangcrypto/ripemd160"
)

const PublicKeyPrefix = "EOS"

type PublicKey []byte

func NewPublicKey(pubKey string) (PublicKey, error) {
	if len(pubKey) < 5 {
		return nil, fmt.Errorf("invalid format")
	}

	if !strings.HasPrefix(pubKey, PublicKeyPrefix) {
		return nil, fmt.Errorf("public key should start with %q", PublicKeyPrefix)
	}

	pubDecoded, err := checkDecode(pubKey[3:])
	if err != nil {
		return nil, fmt.Errorf("checkDecode: %s", err)
	}

	// key, err := btcec.ParsePubKey(pubDecoded, btcec.S256())
	// if err != nil {
	// 	return nil, fmt.Errorf("parsePubKey: %s", err)
	// }

	return PublicKey(pubDecoded), nil
}

func MustNewPublicKey(pubKey string) PublicKey {
	key, err := NewPublicKey(pubKey)
	if err != nil {
		panic(err.Error())
	}
	return key
}

// CheckDecode decodes a string that was encoded with CheckEncode and verifies the checksum.
func checkDecode(input string) (result []byte, err error) {
	decoded := base58.Decode(input)
	if len(decoded) < 5 {
		return nil, fmt.Errorf("invalid format")
	}
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if bytes.Compare(ripemd160checksum(decoded[:len(decoded)-4]), cksum[:]) != 0 {
		return nil, fmt.Errorf("invalid checksum")
	}
	// perhaps bitcoin has a leading net ID / version, but EOS doesn't
	payload := decoded[:len(decoded)-4]
	result = append(result, payload...)
	return
}

func ripemd160checksum(in []byte) []byte {
	h := ripemd160.New()
	_, _ = h.Write(in) // this implementation has no error path
	sum := h.Sum(nil)
	return sum[:4]
}

func (p PublicKey) Key() (*btcec.PublicKey, error) {
	key, err := btcec.ParsePubKey(p, btcec.S256())
	if err != nil {
		return nil, fmt.Errorf("parsePubKey: %s", err)
	}

	return key, nil
}

func (p PublicKey) String() string {
	//rawkey := p.pubKey.SerializeCompressed()
	hash := ripemd160checksum(p)
	rawkey := append(p, hash[:4]...)
	return PublicKeyPrefix + base58.Encode(rawkey)
}

func (p PublicKey) ToHex() string {
	// pubKey.SerializeCompressed()
	return hex.EncodeToString(p)
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

func (p PublicKey) MarshalBinary() ([]byte, error) {
	// str := a.String()
	// raw := base58.Decode(str[3:])
	// raw = raw[:33]
	// ..34-len(raw), raw...
	return append(bytes.Repeat([]byte{0}, 34-len(p)), p...), nil
}

func (p *PublicKey) UnmarshalBinary(data []byte) (err error) {
	if len(data) != 34 {
		return fmt.Errorf("public key should be 34 bytes: \x00 + 33 bytes of key material, had %d bytes", len(data))
	}

	*p = PublicKey(data[1:])
	// newKey, err := NewPublicKey("EOS" + base58.Encode(data))
	// if err != nil {
	// 	return err
	// }

	// a.pubKey = newKey.pubKey

	return nil
}

func (a PublicKey) UnmarshalBinarySize() int { return 34 }
