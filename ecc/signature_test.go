package ecc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignatureSerialization(t *testing.T) {
	privkey, err := NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	require.NoError(t, err)

	payload := []byte("payload")
	sig, err := privkey.Sign(payload)
	require.NoError(t, err)
	assert.Equal(t, `EOSJwRwHEd6yxEpBDYixJM67UtCbpyBzm9rXWNpjs2RE7WYDBsBbk2wgAHik8Bwzbd6fSEWuKoe4TA45sRgRMd5jiQxXdhYnk`, sig.String())
}

func TestSignatureMarshalUnmarshal(t *testing.T) {
	fromEOSIOC := "EOSK5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnNK"
	sig, err := NewSignature(fromEOSIOC)
	require.NoError(t, err)
	assert.Equal(t, fromEOSIOC, sig.String())

}

func TestSignaturePublicKeyExtraction(t *testing.T) {
	// was signed with EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	fromEOSIOC := "EOSK5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnNK"
	sig, err := NewSignature(fromEOSIOC)
	require.NoError(t, err)

	//payload, err := hex.DecodeString("0000000000ea305500001059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100")
	payload, err := hex.DecodeString("20d8af5a0000b32bcc0e37eb0000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500001059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100")
	require.NoError(t, err)

	// The
	pubKey, err := sig.PublicKey(payload)
	require.NoError(t, err)

	// Ok, we'd need to find values where we know the signature is valid, and comes from the given key.
	assert.Equal(t, "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", pubKey.String())
}

func TestEOSIOCSigningComparison(t *testing.T) {
	wif := "5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss" // corresponds to EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	payload, err := hex.DecodeString("20d8af5a0000b32bcc0e37eb0000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500001059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100")
	require.NoError(t, err)

	sig, err := privKey.Sign(payload)
	require.NoError(t, err)

	fromEOSIOC := "EOSK5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnNK"
	assert.Equal(t, fromEOSIOC, sig.String())
}

func TestSignatureUnmarshalChecksum(t *testing.T) {
	fromEOSIOC := "EOSK5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnZZ" // simply checked the last 2 bytes
	_, err := NewSignature(fromEOSIOC)
	require.Equal(t, "signature checksum failed, found 02c9bc70 expected 02c9befc", err.Error())
}
