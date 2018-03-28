package eosapi

import (
	"crypto/sha256"
	"fmt"

	"github.com/eosioca/eosapi/ecc"
)

type Signer interface {
	AvailableKeys() (out []ecc.PublicKey, err error)

	// Sign signs a `tx` transaction. It gets passed a
	// SignedTransaction because it is possible that it holds a few
	// signatures and requests this wallet only to add one or more
	// signatures it requires.
	Sign(tx *SignedTransaction, chainID []byte, requiredKeys ...ecc.PublicKey) (*SignedTransaction, error)

	ImportPrivateKey(wifPrivKey string) error
}

// `eosiowd` wallet-based signer
type WalletSigner struct {
	api        *EOSAPI
	walletName string
}

// NewWalletSigner takes an `api`, because often the wallet will be a
// second endpoint, and not the server node with whom you're pushing
// transactions to.
func NewWalletSigner(api *EOSAPI, walletName string) *WalletSigner {
	return &WalletSigner{api, walletName}
}

func (s *WalletSigner) ImportPrivateKey(wifKey string) (err error) {
	return s.api.WalletImportKey(s.walletName, wifKey)
}

func (s *WalletSigner) AvailableKeys() (out []ecc.PublicKey, err error) {
	return s.api.WalletPublicKeys()
}

func (s *WalletSigner) Sign(tx *SignedTransaction, chainID []byte, requiredKeys ...ecc.PublicKey) (*SignedTransaction, error) {
	// Fetch the available keys over there... and ask this wallet
	// provider to sign with the keys he has..

	// TODO: If there's not a full overlap between the required keys
	// and the available keys, return something about
	// `SignatureIncomplete`.

	resp, err := s.api.WalletSignTransaction(tx, requiredKeys...)
	if err != nil {
		return nil, err
	}

	tx.Signatures = resp.Signatures

	return tx, nil
}

// KeyBag, local signing - NOT COMPLETE

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

func (b *KeyBag) AvailableKeys() (out []ecc.PublicKey, err error) {
	for _, k := range b.Keys {
		out = append(out, k.PublicKey())
	}
	return
}

func (b *KeyBag) ImportPrivateKey(wifPrivKey string) (err error) {
	return b.Add(wifPrivKey)
}

func (b *KeyBag) Sign(tx *SignedTransaction, chainID []byte, requiredKeys ...ecc.PublicKey) (*SignedTransaction, error) {
	txdata, err := MarshalBinary(tx.Transaction)
	if err != nil {
		return nil, err
	}

	keyMap := b.keyMap()
	for _, key := range requiredKeys {
		privKey := keyMap[key.String()]
		if privKey == nil {
			return nil, fmt.Errorf("private key for %q not in keybag", key)
		}

		// TODO: handle ContextFreeData later.. will be added to
		// signature if it exists in tx.ContextFreeData .. and there
		// can be many []byte in there.. so the serialization isn't
		// clear to me yet.  Shouldn't be very complex though.
		sig, err := privKey.Sign(SigDigest(chainID, txdata, nil))
		if err != nil {
			return nil, err
		}

		tx.Signatures = append(tx.Signatures, sig.String())
	}

	return tx, nil
}

func (b *KeyBag) keyMap() map[string]*ecc.PrivateKey {
	out := map[string]*ecc.PrivateKey{}
	for _, key := range b.Keys {
		out[key.PublicKey().String()] = key
	}
	return out
}

func SigDigest(chainID, payload, contextFreeData []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(chainID)
	_, _ = h.Write(payload)
	if len(contextFreeData) > 0 {
		_, _ = h.Write(contextFreeData)
	}
	return h.Sum(nil)
}
