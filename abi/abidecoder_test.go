package abi

import (
	"os"
	"testing"

	"encoding/hex"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestABIDecoder_Decode(t *testing.T) {

	//abiFile := "/Volumes/bb/nodeos/eos/contracts/eosio.token/eosio.token.abi"
	abiFile := "/Volumes/bb/nodeos/eos/contracts/eosio.system/eosio.system.abi"

	abiReader, err := os.Open(abiFile)
	assert.NoError(t, err)

	data, err := hex.DecodeString("a09863fc5094be6900000000000000000150cf44982a1aa36a") //eosio.voteproducer
	assert.NoError(t, err)

	decoder := NewABIDecoder(data, abiReader)
	result := map[string]interface{}{}
	err = decoder.Decode(result, "voteproducer")
	assert.NoError(t, err)
	fmt.Println("grr:", result)
}
