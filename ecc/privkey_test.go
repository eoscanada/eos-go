package ecc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewPrivateKey(t *testing.T) {
	cases := []struct {
		name        string
		key         string
		expectedKey string
	}{
		{
			name:        "K1",
			key:         "PVT_K1_9FN3K4JhzaMsw2Duzr1ijHzVecHtqg1QG4ZCX9udh69Z7QGTk",
			expectedKey: "5HxXwim9PAZZctKJG7Sk6mURD6UXW2hkjDKqnNZu9WYjKD6fF5a",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			key, err := NewPrivateKey(c.key)
			require.NoError(t, err)
			assert.Equal(t, c.expectedKey, key.String())
		})
	}
}
