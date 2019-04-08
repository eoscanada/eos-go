package eos

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionID(t *testing.T) {
	tests := []struct {
		packedTx string
		expectID string
	}{
		{
			// From block 01a528f7ef7412049150d7097acb361051bd2ac3b079ac077912987cc3e94f0b on mainnet
			`{"signatures":["SIG_K1_K4gsBzrZ5dTPrK2dv1bvwttcA7aTuFFyi4X43NDPPxExLvnDxGFpkHx8tmte22sEMKgopcBYT7dvoZgVJ7HFpyQJsrZDuo"],"compression":"none","packed_context_free_data":"","packed_trx":"897ef15ba927136993dd000000000100a6823403ea3055000000572d3ccdcd0190dd39e69a64a64100000000a8ed32322190dd39e69a64a641e05b3597d15cfd45640000000000000004454f530000000000000000000000000000000000000000000000000000000000000000000000000000"}`, // <- 64 bytes of padding in this transaction, which doesn't go in signature.
			"7f84ff0d833c5965f73fb9651881c5b233ab07b45b1a2646b82946f66b78ff92", // the SEEMINGLY GOOD one
			//"ff46c68d0c7fc1e0216dbeb16f52ac27932fdfbff88eb5503ede267bb37f5311", // the WRONG one
		},
	}

	for _, test := range tests {
		var packedTx *PackedTransaction
		err := json.Unmarshal([]byte(test.packedTx), &packedTx)
		require.NoError(t, err)

		id, err := packedTx.ID()
		require.NoError(t, err)

		trxID := hex.EncodeToString(id)
		assert.Equal(t, test.expectID, trxID)
	}
}

func TestTransaction_UnmarshalPacked_Compression(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		expected    CompressionType
		expectedErr error
	}{
		{"string/none", `{"compression": "none"}`, CompressionNone, nil},
		{"string/zlib", `{"compression": "zlib"}`, CompressionZlib, nil},
		{"string/unknown", `{"compression": "random"}`, 0, errors.New("unknown compression type random")},

		{"int/none", `{"compression": 0}`, CompressionNone, nil},
		{"int/zlib", `{"compression": 1}`, CompressionZlib, nil},
		{"int/unknown", `{"compression": 3}`, 0, errors.New("unknown compression type 3")},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var tx *PackedTransaction
			err := json.Unmarshal([]byte(test.in), &tx)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, tx.Compression)
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestTransaction_ExtensionMarshalJSON(t *testing.T) {
	tests := []struct {
		in       *Extension
		expected string
	}{
		{&Extension{1, HexBytes([]byte{0x60, 0x01, 0x2a, 0xbd, 0xef})}, `[1, "60012abdef"]`},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := json.Marshal(test.in)
			require.NoError(t, err)

			assert.JSONEq(t, test.expected, string(actual), string(actual))
		})
	}
}

func TestTransaction_ExtensionUnmarshalJSON(t *testing.T) {
	tests := []struct {
		in       string
		expected Extension
	}{
		{`[1, "60012abdef"]`, Extension{1, HexBytes([]byte{0x60, 0x01, 0x2a, 0xbd, 0xef})}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var actual Extension
			err := json.Unmarshal([]byte(test.in), &actual)
			require.NoError(t, err)

			assert.Equal(t, test.expected, actual)
		})
	}
}
