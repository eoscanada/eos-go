package ecc

import (
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
)

type innerWAPublicKey struct {
}

func newInnerWAPublicKey() innerPublicKey {
	return &innerWAPublicKey{}
}

func (p *innerWAPublicKey) key(content []byte) (*btcec.PublicKey, error) {
	return nil, fmt.Errorf("")
}

func (p *innerWAPublicKey) keyMaterialSize() *int {
	return nil
}

func (p *innerWAPublicKey) prefix() string {
	return PublicKeyWAPrefix
}
