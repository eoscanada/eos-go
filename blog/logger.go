package blog

import (
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

func SetLogger(logger *zap.Logger) {
	zlog = logger
}
