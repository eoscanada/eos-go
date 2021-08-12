package eos_test

import (
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
)

func ExamplePackedTransaction_Unpack() {
	var packedTrx *eos.PackedTransaction
	err := json.Unmarshal(packedTrxData(), &packedTrx)
	if err != nil {
		panic(fmt.Errorf("unmarshaling to PackedTransaction: %w", err))
	}

	var signedTrx *eos.SignedTransaction
	signedTrx, err = packedTrx.Unpack()
	if err != nil {
		panic(fmt.Errorf("unpacking transaction: %w", err))
	}

	fmt.Printf("%#v\n", signedTrx.Actions)
}

func packedTrxData() []byte {
	return []byte(`
		{
		"signatures": [
		  "SIG_K1_K8VSYk76oK4Hdy23UtAJwwRHtBNP8mbu8uo9TVKsT3si5cujPbRqif8eqxqTwLbKREDFm7eK7YG3skLg9LVXZ54KrEoTuJ"
		],
		"compression": "none",
		"packed_context_free_data": "",
		"packed_trx": "a67a815c0d358ee0065800000000011082422e6575305500405647ed48b1ba0140a7c3066575305500000000489aa6b94a1c88ee2531ab18a800201ee9053cde8078023ba1229389f58a0c72ef7fe9ee942e6be7705021630a03e206b016a9711064ee11cc894100701a1160f12c37000903729a1b60f3c7b0117900"
	  }
	`)
}
