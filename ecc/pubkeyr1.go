package ecc

import (
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
)

type innerR1PublicKey struct {
}

func (p *innerR1PublicKey) key(content []byte) (*btcec.PublicKey, error) {
	key, err := btcec.ParsePubKey(content, btcec.S256())
	if err != nil {
		return nil, fmt.Errorf("parsePubKey: %s", err)
	}

	return key, nil
}

func (p *innerR1PublicKey) prefix() string {
	return PublicKeyR1Prefix
}
