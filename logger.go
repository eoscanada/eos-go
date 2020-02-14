package eos

import (
	"fmt"

	"go.uber.org/zap"
)

var coreLog = zap.NewNop()
var encoderLog = zap.NewNop()
var decoderLog = zap.NewNop()
var abiEncoderLog = zap.NewNop()
var abiDecoderLog = zap.NewNop()
var loggingEnabled = false

func EnableDebugLogging(l *zap.Logger) {
	coreLog = l
	encoderLog = l
	decoderLog = l
	abiEncoderLog = l
	abiDecoderLog = l
	loggingEnabled = true
}

func EnableCoreLogging() {
	coreLog = newLogger(false)
	enableLogging(coreLog)
}

func EnableEncoderLogging() {
	encoderLog = newLogger(false)
	enableLogging(encoderLog)
}

func EnableDecoderLogging() {
	decoderLog = newLogger(false)
	enableLogging(decoderLog)
}

func EnableABIEncoderLogging() {
	abiEncoderLog = newLogger(false)
	enableLogging(abiEncoderLog)
}

func EnableABIDecoderLogging() {
	abiDecoderLog = newLogger(false)
	enableLogging(abiDecoderLog)
}

func enableLogging(logger *zap.Logger) {
	if loggingEnabled == false {
		logger.Warn("Enabling logs. Expect performance hits for high throughput")
		loggingEnabled = true
	}
}

type logStringerFunc func() string

func (f logStringerFunc) String() string { return f() }

func typeField(field string, v interface{}) zap.Field {
	return zap.Stringer(field, logStringerFunc(func() string {
		return fmt.Sprintf("%T", v)
	}))
}

func newLogger(production bool) (l *zap.Logger) {
	if production {
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	return
}

// NewLogger a wrap to newLogger
func NewLogger(production bool) *zap.Logger {
	return newLogger(production)
}
