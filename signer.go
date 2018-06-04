package eos

import (
	"crypto/sha256"
	"fmt"

	"os"

	"bufio"

	"strings"

	"github.com/eoscanada/eos-go/ecc"
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
	api        *API
	walletName string
}

// NewWalletSigner takes an `api`, because often the wallet will be a
// second endpoint, and not the server node with whom you're pushing
// transactions to.
func NewWalletSigner(api *API, walletName string) *WalletSigner {
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

	resp, err := s.api.WalletSignTransaction(tx, chainID, requiredKeys...)
	if err != nil {
		return nil, err
	}

	tx.Signatures = resp.Signatures

	return tx, nil
}

// KeyBag, local signing - NOT COMPLETE

// KeyBag holds private keys in memory, for signing transactions.
type KeyBag struct {
	Keys []*ecc.PrivateKey `json:"keys"`
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

func (b *KeyBag) ImportFromFile(path string) error {
	inFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("import keys from file [%s], %s", path, err)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		key := strings.TrimSpace(strings.Split(scanner.Text(), " ")[0])

		if strings.Contains(key, "/") || strings.Contains(key, "#") || strings.Contains(key, ";") {
			return fmt.Errorf("lines should consist of a private key on each line, with an optional whitespace and comment")
		}

		if err := b.Add(key); err != nil {
			return err
		}
	}
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

func (b *KeyBag) SignDigest(digest []byte, requiredKey ecc.PublicKey) (ecc.Signature, error) {

	privateKey := b.keyMap()[requiredKey.String()]
	if privateKey == nil {
		return ecc.Signature{}, fmt.Errorf("private key not found for public key [%s]", requiredKey.String())
	}

	return privateKey.Sign(digest)
}

func (b *KeyBag) Sign(tx *SignedTransaction, chainID []byte, requiredKeys ...ecc.PublicKey) (*SignedTransaction, error) {
	// TODO: probably want to use `tx.packed` and hash the ContextFreeData also.
	txdata, err := MarshalBinary(tx.Transaction)
	if err != nil {
		return nil, err
	}

	cfd := []byte{}
	if len(tx.ContextFreeData) > 0 {
		cfd, err = MarshalBinary(tx.ContextFreeData)
		if err != nil {
			return nil, err
		}
	}

	keyMap := b.keyMap()
	for _, key := range requiredKeys {
		privKey := keyMap[key.String()]
		if privKey == nil {
			return nil, fmt.Errorf("private key for %q not in keybag", key)
		}

		sigDigest := SigDigest(chainID, txdata, cfd)
		// fmt.Println("Signing with", key.String(), privKey.String())
		// fmt.Println("SIGNING THIS DIGEST:", hex.EncodeToString(sigDigest))
		// fmt.Println("SIGNING THIS payload:", hex.EncodeToString(txdata))
		// fmt.Println("SIGNING THIS chainID:", hex.EncodeToString(chainID))
		// fmt.Println("SIGNING THIS cfd:", hex.EncodeToString(cfd))
		sig, err := privKey.Sign(sigDigest)
		if err != nil {
			return nil, err
		}

		tx.Signatures = append(tx.Signatures, sig)
	}

	// tmpcnt, _ := json.Marshal(tx)
	// var newTx *SignedTransaction
	// _ = json.Unmarshal(tmpcnt, &newTx)

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
	if len(chainID) == 0 {
		_, _ = h.Write(make([]byte, 32, 32))
	} else {
		_, _ = h.Write(chainID)
	}
	_, _ = h.Write(payload)

	if len(contextFreeData) > 0 {
		h2 := sha256.New()
		_, _ = h2.Write(contextFreeData)
		_, _ = h.Write(h2.Sum(nil)) // add the hash of CFD to the payload
	} else {
		_, _ = h.Write(make([]byte, 32, 32))
	}
	return h.Sum(nil)
}
