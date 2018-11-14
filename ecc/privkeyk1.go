package ecc

import (
	"fmt"

	"github.com/eoscanada/eos-go/btcsuite/btcd/btcec"
	"github.com/eoscanada/eos-go/btcsuite/btcutil"
)

type InnerK1PrivateKey struct {
	privKey *btcec.PrivateKey
}

func (K *InnerK1PrivateKey) PublicKey() PublicKey {
	return PublicKey{Curve: CurveK1, Content: K.privKey.PubKey().SerializeCompressed(), inner: &InnerK1PublicKey{}}
}

func (k *InnerK1PrivateKey) Sign(hash []byte) (out Signature, err error) {
	if len(hash) != 32 {
		return out, fmt.Errorf("hash should be 32 bytes")
	}

	compactSig, err := k.privKey.SignCanonical(btcec.S256(), hash)

	if err != nil {
		return out, fmt.Errorf("canonical, %s", err)
	}

	return Signature{Curve: CurveK1, Content: compactSig, innerSignature: &InnerK1Signature{}}, nil
}

func (k *InnerK1PrivateKey) String() string {
	wif, _ := btcutil.NewWIF(k.privKey, '\x80', false) // no error possible
	return wif.String()
}
