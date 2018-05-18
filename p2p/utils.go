package p2p

import (
	"encoding/hex"
	"fmt"

	eos "github.com/eoscanada/eos-go"
)

func DecodeHex(hexString string) (data []byte) {
	data, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println("decodeHex error: ", err)
	}
	return data
}

var println = func(args ...interface{}) {
	print(fmt.Sprintf("%s\n", args...))
}

var print = func(s string) {
	if eos.Debug {
		fmt.Print(s)
	}
}
