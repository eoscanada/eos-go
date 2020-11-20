package ecc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PublicKeyMarshalUnmarshal(t *testing.T) {
	cases := []struct {
		name        string
		key         string
		expectedKey string
	}{
		{
			name:        "K1-EOS",
			key:         PublicKeyPrefixCompat + "6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expectedKey: PublicKeyPrefixCompat + "6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
		},
		{
			name:        "K1",
			key:         "PUB_K1_6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expectedKey: PublicKeyPrefixCompat + "6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
		},
		{
			name:        "R1",
			key:         "PUB_R1_78rbUHSk87e7eCBoccgWUkhNTCZLYdvJzerDRHg6fxj2SQy6Xm",
			expectedKey: "PUB_R1_78rbUHSk87e7eCBoccgWUkhNTCZLYdvJzerDRHg6fxj2SQy6Xm",
		},
		{
			name:        "WA",
			key:         "PUB_WA_5hyixc7vkMbKiThWi1TnFtXw7HTDcHfjREj2SzxCtgw3jQGepa5T9VHEy1Tunjzzj",
			expectedKey: "PUB_WA_5hyixc7vkMbKiThWi1TnFtXw7HTDcHfjREj2SzxCtgw3jQGepa5T9VHEy1Tunjzzj",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			key, err := NewPublicKey(c.key)
			require.NoError(t, err)
			assert.Equal(t, c.expectedKey, key.String())
		})
	}
}

func TestPublicKey_MarshalJSON(t *testing.T) {
	cases := []struct {
		name         string
		key          string
		expectedJSON string
	}{
		{
			name:         "K1-EOS",
			key:          PublicKeyPrefixCompat + "6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expectedJSON: `"` + PublicKeyPrefixCompat + `6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"`,
		},
		{
			name:         "K1",
			key:          "PUB_K1_6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expectedJSON: `"` + PublicKeyPrefixCompat + `6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"`,
		},
		{
			name:         "R1",
			key:          "PUB_R1_78rbUHSk87e7eCBoccgWUkhNTCZLYdvJzerDRHg6fxj2SQy6Xm",
			expectedJSON: `"PUB_R1_78rbUHSk87e7eCBoccgWUkhNTCZLYdvJzerDRHg6fxj2SQy6Xm"`,
		},
		{
			name:         "WA",
			key:          "PUB_WA_5hyixc7vkMbKiThWi1TnFtXw7HTDcHfjREj2SzxCtgw3jQGepa5T9VHEy1Tunjzzj",
			expectedJSON: `"PUB_WA_5hyixc7vkMbKiThWi1TnFtXw7HTDcHfjREj2SzxCtgw3jQGepa5T9VHEy1Tunjzzj"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			key, err := NewPublicKey(c.key)
			require.NoError(t, err)

			result, err := key.MarshalJSON()
			require.NoError(t, err)

			assert.Equal(t, []byte(c.expectedJSON), result)
		})
	}
}

func TestPublicKey_ToKey_WA(t *testing.T) {
	key, err := NewPublicKey("PUB_WA_5hyixc7vkMbKiThWi1TnFtXw7HTDcHfjREj2SzxCtgw3jQGepa5T9VHEy1Tunjzzj")
	require.NoError(t, err)

	btcecKey, err := key.Key()
	require.NoError(t, err)

	ecdsaKey := btcecKey.ToECDSA()
	assert.Equal(t, "364ee3c86f1a4159576e46078431a9906b44ec2bdc720ec4dbea4afae0ac643b", hex.EncodeToString(ecdsaKey.X.Bytes()))
	assert.Equal(t, "87ae0fc0799400f8e1320692d5a7bc0cf51190d182ea4ec69a60f38177568550", hex.EncodeToString(ecdsaKey.Y.Bytes()))
}
