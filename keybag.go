package eosapi

import (
	"fmt"

	"github.com/eosioca/eosapi/ecc"
)

// KeyBag holds private keys in memory, for signing transactions.
type KeyBag struct {
	Keys []*ecc.PrivateKey
}

func NewKeyBag() *KeyBag {
	return &KeyBag{
		Keys: make([]*ecc.PrivateKey, 0),
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

func (b *KeyBag) Sign(tx *Transaction, requiredKeys ...PublicKey) (*SignedTransaction, error) {
	s := &SignedTransaction{
		Transaction: tx,
	}

	keyMap := b.keyMap()

	hashdata := []byte("hashed-data-of-tx--we-need-to-serialize-the-whole-thing-thank-you")
	for _, key := range requiredKeys {
		privKey := keyMap[string(key)]
		if privKey == nil {
			return nil, fmt.Errorf("private key for %q not in keybag", key)
		}

		sig, err := privKey.Sign(hashdata)
		if err != nil {
			return nil, err
		}

		s.Signatures = append(s.Signatures, sig.String())
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
