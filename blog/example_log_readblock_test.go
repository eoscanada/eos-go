package blog_test

import (
	"encoding/json"
	"fmt"

	"github.com/eoscanada/eos-go/blog"
	"github.com/pkg/errors"
)

func ExampleBlockLog_ReadBlock() {
	logFilename := "/tmp/blocks.log"
	indexFilename := "/tmp/blocks.index"
	blockLog := blog.NewFileBlockLog(logFilename, indexFilename)

	err := blockLog.Open()
	if err != nil {
		panic(errors.Wrap(err, "unable to open block log"))
	}

	defer blockLog.Close()

	blockNum := uint32(3)
	block, err := blockLog.ReadBlock(blockNum)
	if err != nil {
		panic(errors.Wrapf(err, "unable to read block %d", blockNum))
	}

	jsonBlock, err := json.Marshal(block)
	if err != nil {
		panic(errors.Wrap(err, "unable to marshal block to JSON"))
	}

	fmt.Printf("Block\n%s\n", jsonBlock)
}
