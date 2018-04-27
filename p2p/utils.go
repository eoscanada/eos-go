package p2p

import (
	"encoding/hex"
	"fmt"
)

func decodeHex(hexString string) (data []byte) {
	data, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println("decodeHex error: ", err)
	}
	return data
}
