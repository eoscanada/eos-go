package eos

import (
	"encoding/hex"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestPackedTransaction_UnPack(t *testing.T) {

	hexString := "7353f05a0000f23fc87d8a98000000000100a6823403ea3055000000572d3ccdcd010000000000ea305500000000a8ed3232210000000000ea305500000039ab18dd41a08601000000000004454f530000000000"
	data, err := hex.DecodeString(hexString)
	assert.NoError(t, err)
	tx := PackedTransaction{
		PackedTransaction: data,
	}

	signedTx, err := tx.UnPack()
	assert.NoError(t, err)
	fmt.Println(signedTx)
}
