package eosapi

import (
	"crypto/sha256"
	"fmt"

	"github.com/eosioca/eosapi/ecc"
)

// KeyBag holds private keys in memory, for signing transactions.
type KeyBag struct {
	ChainID []byte
	Keys    []*ecc.PrivateKey
}

func NewKeyBag(chainID []byte) *KeyBag {
	return &KeyBag{
		ChainID: chainID,
		Keys:    make([]*ecc.PrivateKey, 0),
	}
}

func (b *KeyBag) Add(wifKey string) error {
	privKey, err := ecc.NewPrivateKey(wifKey)
	if err != nil {
		return err
	}
	b.Keys = append(b.Keys, privKey)
	return nil
}

func (b *KeyBag) AvailableKeys() (out []string) {
	for _, k := range b.Keys {
		out = append(out, k.PublicKey().String())
	}
	return
}

func (b *KeyBag) Sign(tx *Transaction, chainID []byte, requiredKeys ...PublicKey) (*SignedTransaction, error) {
	s := &SignedTransaction{
		Transaction: tx,
	}

	txdata, err := MarshalBinary(tx)
	if err != nil {
		return nil, err
	}

	// chainID, err := hex.DecodeString(hexChainID)
	// if err != nil {
	// 	return nil, err
	// }

	hash := sha256.New()
	_, _ = hash.Write(chainID)
	_, _ = hash.Write(txdata)
	hashdata := hash.Sum(nil)

	keyMap := b.keyMap()
	for _, key := range requiredKeys {
		privKey := keyMap[string(key)]
		if privKey == nil {
			return nil, fmt.Errorf("private key for %q not in keybag", key)
		}

		sig, err := privKey.Sign(hashdata)
		if err != nil {
			return nil, err
		}

		s.Signatures = append(s.Signatures, sig)
	}

	return s, nil
}

func (b *KeyBag) keyMap() map[string]*ecc.PrivateKey {
	out := map[string]*ecc.PrivateKey{}
	for _, key := range b.Keys {
		out[key.PublicKey().String()] = key
	}
	return out
}
