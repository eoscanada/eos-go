package blog_test

import (
	"fmt"

	"github.com/eoscanada/eos-go/blog"
	"github.com/pkg/errors"
)

func ExampleBlockIndex_ReadBlockByteOffset() {
	indexFilename := "/tmp/blocks.index"
	blockIndex := blog.NewFileBlockIndex(indexFilename)

	err := blockIndex.Open()
	if err != nil {
		panic(errors.Wrap(err, "unable to open index"))
	}

	defer blockIndex.Close()

	blockNum := uint32(1)
	blockByteOffset, err := blockIndex.ReadBlockByteOffset(blockNum)
	if err != nil {
		panic(errors.Wrapf(err, "unable to read block byte offset for block num %d", blockNum))
	}

	fmt.Printf("Block byte offset %016x (%d)\n", blockByteOffset, blockByteOffset)
}
