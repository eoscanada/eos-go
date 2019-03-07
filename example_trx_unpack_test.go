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
		panic(fmt.Errorf("unmarshaling to PackedTransaction: %s", err))
	}

	signedTrx, err := packedTrx.Unpack()
	if err != nil {
		panic(fmt.Errorf("unpacking transaction: %s", err))
	}

	fmt.Printf("%#v\n", signedTrx.Actions)
}

func packedTrxData() []byte {
	return []byte(`
		{
			"signatures": [
				"SIG_K1_KcVC8F2bH5ETYRNeZwK27PQW6WVGmPcB1kGHYkT7sqH91JaY3YuLt5UuFo5w5o2QzLfMEXRdwWaH8qesSiD3MaKyW732Jq"
			],
			"compression": "none",
			"packed_context_free_data": "",
			"packed_trx": "714d7d5c5149e2bb30a9000000000100a6823403ea3055000000572d3ccdcd01000000000060b64a00000000a8ed323227000000000060b64a000000000070b64a307500000000000004454f53000000000674657374203100"
		}
	`)
}
