package ecc

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**

$ cleos wallet private_keys -n r1test
password: [[
    "EOS67SCWnz6trqFPCtmxfjYEPSsT9JKRn4zhow8X3VTtgaEzNMULF",
    "5JaKaxySEyjBFGT9K9cYKSFhfojn1RfPcresqRVbmtxnQt1w3qW"
  ],[
    "PUB_R1_6RJ9pXJNe1wk6p2yiJcuJ8QPo7WTudHya9z8vu1VPk44fhBz79",
    "PVT_R1_2o5WfMRU4dTp23pbcbP2yn5MumQzSMy3ayNQ31qi5nUfa2jdWC"
  ],[
    "PUB_R1_7aE3zt3f7cfNuuUwLogDtxSsniQA2uPthATQZ5ErQLuu1nDKFG",
    "PVT_R1_rjKe476v6zXntjC93YAGyqL35NJWshbwcbGRwb27wuKvsRVEa"
  ],[
    "PUB_R1_8KT5dWt33np9V4Nqpdja1GAbkEqVY3pupeYgvCkKTA5FeqePTp",
    "PVT_R1_2FiHVhVjDNjRVAbLg9Cwj1PvVu6Dxn4HKDMFmkyhPZRdAfXwk6"
  ],[
    "PUB_R1_8S4TodyXa9KASMAJgkLbstFYzAWHNjNJPhpHuqqHF9Af8ekV7i",
    "PVT_R1_2sPCnkH6652KFYQZNWuQvgfTTHvqjrhV6pQ8tcVQGqBNsopKZp"
  ]
]

$ echo -n 'banana' | shasum -a 256
b493d48364afe44d11c0165cf470a4164d1e2609911ef998be868d46ade3de4e  -

$ curl --data '["b493d48364afe44d11c0165cf470a4164d1e2609911ef998be868d46ade3de4e","PUB_R1_6RJ9pXJNe1wk6p2yiJcuJ8QPo7WTudHya9z8vu1VPk44fhBz79"]'
http://127.0.0.1:8900/v1/wallet/sign_digest
"SIG_R1_KJmGMknL29w1jTDbkm4wCB5Lr7UXLLWQrfdyurw8dGoTeHggoVbB9wErfUeFhJXwbihuQHK4G4VeaWoNdW7fdScF92Ctx5"

*/

func TestK1PrivateToPublic(t *testing.T) {
	wif := "5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss"
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	pubKey := privKey.PublicKey()

	pubKeyString := pubKey.String()
	assert.Equal(t, "EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM", pubKeyString)
}

func TestR1PrivateToPublic(t *testing.T) {
	encoded_privKey := "PVT_R1_2o5WfMRU4dTp23pbcbP2yn5MumQzSMy3ayNQ31qi5nUfa2jdWC"
	_, err := NewPrivateKey(encoded_privKey)
	require.NoError(t, err)

	// FIXME: Actual retrieval of publicKey from privateKey for R1 is not done yet, disable this check
	// pubKey := privKey.PublicKey()

	//pubKeyString := pubKey.String()
	//assert.Equal(t, "PUB_R1_0000000000000000000000000000000000000000000000", pubKeyString)
}

func TestNewPublicKeyAndSerializeCompress(t *testing.T) {
	// Copied test from eosjs(-.*)?
	key, err := NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
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
		{"EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM", nil},
		{"MMM859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM", fmt.Errorf("public key should start with [\"PUB_K1_\" | \"PUB_R1_\"] (or the old \"EOS\")")},
		{"EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhTo", fmt.Errorf("checkDecode: invalid checksum")},
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
