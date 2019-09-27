package eos

import (
	"os"

	"go.uber.org/zap"
)

func init() {
	if os.Getenv("DEBUG") != "" {
		encoderLog, _ = zap.NewDevelopment()
		decoderLog, _ = zap.NewDevelopment()
		abiEncoderLog, _ = zap.NewDevelopment()
		abiDecoderLog, _ = zap.NewDevelopment()
		loggingEnabled = true
	}
}
