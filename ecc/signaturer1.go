package ecc

import (
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
)

type innerR1Signature struct {
}

func (s innerR1Signature) verify(content []byte, hash []byte, pubKey PublicKey) bool {
	return false
}

func (s *innerR1Signature) publicKey(content []byte, hash []byte) (out PublicKey, err error) {
	return out, fmt.Errorf("R1 not supported")
}

func (s innerR1Signature) string(content []byte) string {
	return "SIG_R1_" + base58.Encode(content)
}
