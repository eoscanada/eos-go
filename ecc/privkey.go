package ecc

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func NewRandomPrivateKey() (*PrivateKey, error) {
	return newRandomPrivateKey(cryptorand.Reader)
}

func NewDeterministicPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	return newRandomPrivateKey(randSource)
}

func newRandomPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	rawPrivKey := make([]byte, 32)
	written, err := io.ReadFull(randSource, rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("error feeding crypto-rand numbers to seed ephemeral private key: %s", err)
	}
	if written != 32 {
		return nil, fmt.Errorf("couldn't write 32 bytes of randomness to seed ephemeral private key")
	}

	h := sha256.New()
	h.Write(rawPrivKey)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), h.Sum(nil))

	return &PrivateKey{privKey: privKey}, nil
}

func NewPrivateKey(wif string) (*PrivateKey, error) {
	wifObj, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{privKey: wifObj.PrivKey}, nil
}

type PrivateKey struct {
	privKey *btcec.PrivateKey
}

func (p *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{pubKey: p.privKey.PubKey()}
}

func (p *PrivateKey) Sign(payload []byte) (Signature, error) {
	h := sha256.New()
	h.Write(payload)
	hash := h.Sum(nil)

	compactSig, err := btcec.SignCompact(btcec.S256(), p.privKey, hash, false)
	if err != nil {
		return nil, err
	}

	return Signature(compactSig), nil
}

func (p *PrivateKey) String() string {
	wif, _ := btcutil.NewWIF(p.privKey, &chaincfg.Params{PrivateKeyID: '\x80'}, false) // no error possible
	return wif.String()
}
