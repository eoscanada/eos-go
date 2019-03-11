package blog_test

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/eoscanada/eos-go"

	"github.com/eoscanada/eos-go/blog"
	"github.com/pkg/errors"
)

func init() {
	logger, _ := zap.NewDevelopment()

	blog.SetLogger(logger)
	eos.EnableDecoderLogging()
}

func ExampleBlockLog_ForEach() {
	logFilename := "/tmp/blocks.log"
	indexFilename := "/tmp/blocks.index"
	blockLog := blog.NewFileBlockLog(logFilename, indexFilename)

	err := blockLog.Open()
	if err != nil {
		panic(errors.Wrap(err, "unable to open block log"))
	}

	defer blockLog.Close()

	err = blockLog.ForEach(func(block *eos.SignedBlock) error {
		blockID, err := block.BlockID()
		if err != nil {
			return errors.Wrap(err, "unable to compute block id")
		}

		fmt.Printf("Read block %s (#%d)\n", blockID, block.BlockNumber())
		return nil
	})

	if err != nil {
		panic(errors.Wrap(err, "unable to read all blocks"))
	}
}
