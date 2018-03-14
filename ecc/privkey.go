package ecc

import (
	cryptorand "crypto/rand"
	"fmt"
	"io"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

/**
 * PrivateKey
 */

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

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), rawPrivKey)

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

func (p *PrivateKey) String() string {
	//     var private_key = toBuffer();
	// // checksum includes the version
	// private_key = Buffer.concat([new Buffer([0x80]), private_key]);
	// var checksum = hash.sha256(private_key);
	// checksum = hash.sha256(checksum);
	// checksum = checksum.slice(0, 4);
	// var private_wif = Buffer.concat([private_key, checksum]);
	// return base58.encode(private_wif);

	wif, _ := btcutil.NewWIF(p.privKey, &chaincfg.Params{PrivateKeyID: '\x80'}, true) // no error possible
	return wif.String()
}
