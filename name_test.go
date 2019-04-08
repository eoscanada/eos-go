package eos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtendedStringToName(t *testing.T) {
	tests := []struct {
		in            string
		expectedValue uint64
		expectedErr   error
	}{
		{"eosio.token", 0x5530ea033482a600, nil},
		{"1,CUSD", 293455872769, nil},
		{"4,EOS", 1397703940, nil},
		{"CUSD", 1146312003, nil},
		{"KARMA", 280470110539, nil},
		{"IQ", 20809, nil},
		{"EOS", 5459781, nil},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d_%s", i, test.in), func(t *testing.T) {
			actual, err := ExtendedStringToName(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expectedValue, uint64(actual))
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}
