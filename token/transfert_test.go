package token

import (
	"encoding/hex"
	"fmt"
	"testing"

	eos "github.com/eoscanada/eos-go"
	"github.com/stretchr/testify/assert"
)

func FixmeTestPackedTransaction_Unpack(t *testing.T) {

	transfer := Transfer{}
	fmt.Println(transfer)

	hexString := "7353f05a0000f23fc87d8a98000000000100a6823403ea3055000000572d3ccdcd010000000000ea305500000000a8ed3232210000000000ea305500000039ab18dd41a08601000000000004454f530000000000"
	data, err := hex.DecodeString(hexString)
	assert.NoError(t, err)
	tx := eos.PackedTransaction{
		PackedTransaction: data,
	}

	signedTx, err := tx.Unpack()
	assert.NoError(t, err)
	fmt.Println(signedTx)
}
