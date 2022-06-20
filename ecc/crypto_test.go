package ecc

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestK1PrivateToPublic(t *testing.T) {
	wif := "5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss"
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	pubKey := privKey.PublicKey()

	pubKeyString := pubKey.String()
	assert.Equal(t, PublicKeyPrefixCompat+"859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM", pubKeyString)
}

func TestPrefixedK1PrivateToPublic(t *testing.T) {
	wif := "PVT_K1_9FN3K4JhzaMsw2Duzr1ijHzVecHtqg1QG4ZCX9udh69Z7QGTk"
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	pubKey := privKey.PublicKey()

	pubKeyString := pubKey.String()
	assert.Equal(t, PublicKeyPrefixCompat+"7LrH8N3f3BTCLRHQWeo9gVfuBB6XgEqtjksKoN9jhjFjbaGQES", pubKeyString)
}

func TestR1PrivateToPublic(t *testing.T) {
	encodedPrivKey := "PVT_R1_2o5WfMRU4dTp23pbcbP2yn5MumQzSMy3ayNQ31qi5nUfa2jdWC"
	_, err := NewPrivateKey(encodedPrivKey)
	require.NoError(t, err)

	// FIXME: Actual retrieval of publicKey from privateKey for R1 is not done yet, disable this check
	// pubKey := privKey.PublicKey()

	// pubKeyString := pubKey.String()
	// assert.Equal(t, "PUB_R1_0000000000000000000000000000000000000000000000", pubKeyString)
}

func TestNewPublicKeyAndSerializeCompress(t *testing.T) {
	// Copied test from eosjs(-.*)?
	key, err := NewPublicKey(PublicKeyPrefixCompat + "6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	require.NoError(t, err)
	assert.Equal(t, "02c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf", hex.EncodeToString(key.Content))
}

func TestNewRandomPrivateKey(t *testing.T) {
	key, err := NewRandomPrivateKey()
	require.NoError(t, err)
	// taken from eosiojs-ecc:common.test.js:12
	assert.Regexp(t, "^5[HJK].*", key.String())
}

func TestPrivateKeyValidity(t *testing.T) {
	tests := []struct {
		in    string
		valid bool
	}{
		{"5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss", true},
		{"5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjsm", false},
	}

	for _, test := range tests {
		_, err := NewPrivateKey(test.in)
		if test.valid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
			assert.Equal(t, "checksum mismatch", err.Error())
		}
	}
}

func TestPublicKeyValidity(t *testing.T) {
	tests := []struct {
		in  string
		err error
	}{
		{PublicKeyPrefixCompat + "859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM", nil},
		{"MMM859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM", fmt.Errorf(`public key should start with "PUB_K1_", "PUB_R1_", "PUB_WA_" or the old "` + PublicKeyPrefixCompat + `"`)},
		{PublicKeyPrefixCompat + "859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhTo", fmt.Errorf("public key checksum failed, found 0e2e1094 but expected 169c2652")},
	}

	for idx, test := range tests {
		_, err := NewPublicKey(test.in)
		if test.err == nil {
			assert.NoError(t, err, fmt.Sprintf("test %d with key %q", idx, test.in))
		} else {
			assert.Error(t, err)
			assert.Equal(t, test.err.Error(), err.Error())
		}
	}
}

func TestK1Signature(t *testing.T) {
	wif := "5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss"
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	cnt := []byte("hi")
	digest := sigDigest([]byte{}, cnt, nil)
	signature, err := privKey.Sign(digest)
	require.NoError(t, err)

	assert.True(t, signature.Verify(digest, privKey.PublicKey()))
}

func TestR1Signature(t *testing.T) {
	encodedPrivKey := "PVT_R1_2o5WfMRU4dTp23pbcbP2yn5MumQzSMy3ayNQ31qi5nUfa2jdWC"
	privKey, err := NewPrivateKey(encodedPrivKey)
	require.NoError(t, err)

	cnt := []byte("hi")
	digest := sigDigest([]byte{}, cnt, nil)
	_, err = privKey.Sign(digest)
	assert.Error(t, err)
	assert.Equal(t, "R1 not supported", err.Error())
}
