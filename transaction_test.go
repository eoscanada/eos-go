package eos

import (
	"encoding/hex"
	"encoding/json"
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
