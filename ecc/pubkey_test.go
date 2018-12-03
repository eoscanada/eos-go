package ecc

import (
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
			key:         "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expectedKey: "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
		},
		{
			name:        "K1",
			key:         "PUB_K1_6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expectedKey: "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
		},
		{
			name:        "R1",
			key:         "PUB_R1_78rbUHSk87e7eCBoccgWUkhNTCZLYdvJzerDRHg6fxj2SQy6Xm",
			expectedKey: "PUB_R1_78rbUHSk87e7eCBoccgWUkhNTCZLYdvJzerDRHg6fxj2SQy6Xm",
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
			key:          "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expectedJSON: `"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"`,
		},
		{
			name:         "K1",
			key:          "PUB_K1_6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expectedJSON: `"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"`,
		},
		{
			name:         "R1",
			key:          "PUB_R1_78rbUHSk87e7eCBoccgWUkhNTCZLYdvJzerDRHg6fxj2SQy6Xm",
			expectedJSON: `"PUB_R1_78rbUHSk87e7eCBoccgWUkhNTCZLYdvJzerDRHg6fxj2SQy6Xm"`,
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
