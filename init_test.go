package eos

import (
	"time"

	"github.com/streamingfast/logging"
)

func init() {
	logging.InstantiateLoggers()

	time.Local, _ = time.LoadLocation("America/New_York")
}
