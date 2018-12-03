package p2p

import (
	"encoding/hex"
)

func DecodeHex(hexString string) (data []byte) {
	data, err := hex.DecodeString(hexString)
	if err != nil {
		logErr("decodeHexErr", err)
	}
	return data
}
