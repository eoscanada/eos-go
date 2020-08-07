package eos

import (
	"os"

	"go.uber.org/zap"
)

func init() {
	if os.Getenv("DEBUG") != "" {
		logger, _ := zap.NewDevelopment()
		EnableDebugLogging(logger)
	}
}
