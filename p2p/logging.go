package p2p

import (
	"fmt"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

func init() {
	logging.Register("github.com/eoscanada/eos-go/p2p", &zlog)
}

// SyncLogger sync logger, should `defer SyncLogger()` when use p2p package
func SyncLogger() {
	err := zlog.Sync()
	if err != nil {
		fmt.Printf("unable to sync p2p logger: %s\n", err)
	}
}
