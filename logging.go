package eos

import (
	"fmt"

	"github.com/dfuse-io/logging"
	"go.uber.org/zap"
)

var traceEnabled = logging.IsTraceEnabled("eos-go", "github.com/eoscanada/eos-go")
var zlog = zap.NewNop()

func init() {
	logging.Register("github.com/eoscanada/eos-go", &zlog)
}

func EnableDebugLogging(l *zap.Logger) {
	traceEnabled = true
	zlog = l
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
