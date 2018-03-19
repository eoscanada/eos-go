package ecc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignatureSerialization(t *testing.T) {
	privkey, err := NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	require.NoError(t, err)
	sig, err := privkey.Sign([]byte("payload"))
	require.NoError(t, err)
	assert.Equal(t, `EOSJwRwHEd6yxEpBDYixJM67UtCbpyBzm9rXWNpjs2RE7WYDBsBbk2wgAHik8Bwzbd6fSEWuKoe4TA45sRgRMd5jiQxXdhYnk`, sig.String())
}

func TestSignatureMarshalUnmarshal(t *testing.T) {
	fromEOSIOC := "EOSK5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnNK"
	sig, err := NewSignature(fromEOSIOC)
	require.NoError(t, err)
	assert.Equal(t, fromEOSIOC, sig.String())
}

func TestSignatureUnmarshalChecksum(t *testing.T) {
	fromEOSIOC := "EOSK5yY5ehsnDMc6xcRhsLYzFuZGUaKwb4hc8oLmP5HA1EhU42NRo3ygx3zvLRJ1nkw1NA5nCSegwcYkSfkZBQBzqMDsCGnZZ" // simply checked the last 2 bytes
	_, err := NewSignature(fromEOSIOC)
	require.Equal(t, "signature checksum failed, found 02c9bc70 expected 02c9befc", err.Error())
}
