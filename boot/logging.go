package boot

import (
	"github.com/eoscanada/eos-go"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

func init() {
	zlog = eos.NewLogger(false)

}
