package p2p

import (
	"encoding/hex"

	"go.uber.org/zap"
)

func DecodeHex(hexString string) (data []byte) {
	data, err := hex.DecodeString(hexString)
	if err != nil {
		zlog.Error("decode hex err", zap.Error(err))
	}
	return data
}
