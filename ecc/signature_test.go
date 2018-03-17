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
