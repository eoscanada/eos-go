package ecc

import (
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

type innerK1PublicKey struct {
}

func (p *innerK1PublicKey) key(content []byte) (*btcec.PublicKey, error) {
	key, err := btcec.ParsePubKey(content, btcec.S256())
	if err != nil {
		return nil, fmt.Errorf("parsePubKey: %s", err)
	}

	return key, nil
}

func (p *innerK1PublicKey) string(content []byte, curveID CurveID) string {

	data := content
	if len(data) == 0 {
		data = make([]byte, 33)
	}

	hash := ripemd160checksum(data, curveID)

	rawKey := make([]byte, 37)
	copy(rawKey, data[:33])
	copy(rawKey[33:], hash[:4])

	return PublicKeyPrefixCompat + base58.Encode(rawKey)
}
