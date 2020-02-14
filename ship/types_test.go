package ship

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/eoscanada/eos-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func init() {
	if os.Getenv("DEBUG") != "" {
		logger, _ := zap.NewDevelopment()
		eos.EnableDebugLogging(logger)
	}
}
func TestDecodeGetBlocksResultV0(t *testing.T) {
	resultStr := "0120ac01000001ac201bfa206d4a92c1663fc9e8cc69e843ff91f49d396f03bf53761708011fac01000001ac1faccb177cfa87f069aa202b6cd1af243366c84ae510618bc8db8cd568010600000000000006ff3859baa3329bd727cc669799fb9a26b8a4e6e64256bd4029a23a9d010500000000000005889366dad98f29645b419d725d6310cd8d983a6844c998095427d5a801b8011d5df3480000000000ea3055000000000005889366dad98f29645b419d725d6310cd8d983a6844c998095427d5a800000000000000000000000000000000000000000000000000000000000000007f45c72e725b22a1933f1dca9eeb319d6d278e78ba103f612754c7da026a5cc8000000000000002049c27702b2bcd396fc5227e22d5d42c9443b42c3cc868a94ea05f284be1b0eab7423f71d472950419e7f6f6022fe9a3db67a231721efc678ae5b9ec8f4e49131000001cb020100ff6212acca554e91cd6ae901d821672e99927314c457d246c9dfa63c1853c39600640000000065000000000000000000000000000000000100010001000000000000ea30556b51a9907006a417657bf6902f0ef18e063c9211b8f65ad8b7dc78d7d676734f39000000000000003700000000000000010000000000ea30551c0000000000000002020000000000ea30550000000000ea305500000000221acfa4010000000000ea305500000000a8ed3232741c5df3480000000000ea3055000000000004e9c3cfaa93d04d2c331331fc30700b14ba42963cbd42a45309de2e48c9728c2d18bb29abc87d34650cf3ac4f10efd9611d8b2fb45c51d2d88e3bb4949c88ceee48192983b774d4100c27f700e4053a046a3569e0326030940cf330a4000000000000003d00000000000000000000000000000001000ff2e65c0500d98f296400000000000001a6030400106163636f756e745f6d6574616461746100000c636f6e74726163745f726f770101c501000000000000ea30550000000000ea3055000000004473686400000000447368640000000000ea30559a010000100000000000e8030000000008000c000000f40100001400000064000000400d0300e8030000f049020064000000100e00005802000080533b0000100000040006000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e7265736f757263655f757361676501013b000000000000ea3055001d5df348da3e0d00000000000100000000000000001d5df348c2c90600000000006500000000000000200512000000000000157265736f757263655f6c696d6974735f73746174650101530000060000008876fe910000000092090000000000000006000000b4e2115000000000a30500000000000000000000000000000000000000000000000000000000000087141000000000002911030000000000"
	result, err := hex.DecodeString(resultStr)
	require.NoError(t, err)

	out := &Result{}
	err = eos.UnmarshalBinary(result, &out)
	require.NoError(t, err)

	require.Equal(t, GetBlocksResultV0Type, int(out.TypeID))
	res := out.Impl.(*GetBlocksResultV0)

	//fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.Block)
	//	fmt.Printf("%+v\n", res.Traces.Elem[0].Impl.(*TransactionTraceV0))
	//	if len(res.Deltas.Elem) > 0 {
	//		fmt.Printf("%+v\n", res.Deltas.Elem[0].Impl.(*TableDeltaV0))
	//	}

}

func TestGetSmallBlockResult(t *testing.T) {
	t.Skip()
	result, err := hex.DecodeString("0120ac01000001ac201bfa206d4a92c1663fc9e8cc69e843ff91f49d396f03bf53761708011fac01000001ac1faccb177cfa87f069aa202b6cd1af243366c84ae510618bc8db8cd568010500000000000005889366dad98f29645b419d725d6310cd8d983a6844c998095427d5a8010400000000000004e9c3cfaa93d04d2c331331fc30700b14ba42963cbd42a45309de2e48000000")
	require.NoError(t, err)

	out := &Result{}
	err = eos.UnmarshalBinary(result, &out)
	require.NoError(t, err)
	require.Equal(t, GetBlocksResultV0Type, int(out.TypeID))
	bresult := out.Impl.(*GetBlocksResultV0)
	assert.EqualValues(t, bresult.ThisBlock.BlockNum, 5)
	assert.EqualValues(t, bresult.PrevBlock.BlockNum, 4)
	assert.EqualValues(t, bresult.Head.BlockNum, 109600)
	assert.EqualValues(t, bresult.LastIrreversible.BlockNum, 109599)
	assert.Equal(t, bresult.LastIrreversible.BlockID.String(), "0001ac1faccb177cfa87f069aa202b6cd1af243366c84ae510618bc8db8cd568")
}

func TestGetSmallBlockContent(t *testing.T) {
	result, err := hex.DecodeString(blockdataForBlock4)
	require.NoError(t, err)

	out := &SignedBlock{}
	err = eos.UnmarshalBinary(result, &out)
	require.NoError(t, err)

	for _, x := range out.Transactions {
		fmt.Println(x.Status, x.CPUUsageUS, x.NetUsageWords)
		if trx, ok := x.Trx.Impl.(*eos.PackedTransaction); ok {
			untx, err := trx.Unpack()
			assert.NoError(t, err)
			fmt.Println(untx)
		}
	}
}

func TestBiggerGetBlocksResultV0(t *testing.T) {
	// TODO still broken, fix me
	t.Skip()
	result, err := hex.DecodeString(longresult)
	require.NoError(t, err)

	out := &Result{}
	err = eos.UnmarshalBinary(result, &out)
	require.NoError(t, err)

	require.Equal(t, GetBlocksResultV0Type, int(out.TypeID))
	res := out.Impl.(*GetBlocksResultV0)

	_ = res
	//fmt.Printf("%+v\n", res.Head)
	//fmt.Printf("%+v\n", res.LastIrreversible)
	//fmt.Printf("%+v\n", res.ThisBlock)
	//fmt.Printf("%+v\n", res.PrevBlock)
	//fmt.Printf("%+v\n", res.Block)
	//fmt.Printf("%+v\n", res.Traces.Elem[0].Impl.(*TransactionTraceV0))
	//if len(res.Deltas.Elem) > 0 {
	//	fmt.Printf("%+v\n", res.Deltas.Elem[0].Impl.(*TableDeltaV0))
	//}

}
