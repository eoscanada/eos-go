package blockslog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMe(t *testing.T) {
	t.Skip("Update me so that it's not tied to a particular machine!")

	require.NoError(t, Process("/home/abourget/dfuse/dfuse-eosio/proj/mainnet/mindreader/data/blocks/blocks.log"))
}
