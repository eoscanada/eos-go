package blockslog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMe(t *testing.T) {
	require.NoError(t, Process("/home/abourget/dfuse/dfuse-eosio/proj/mainnet/mindreader/data/blocks/blocks.log"))
}
