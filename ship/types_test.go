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

func TestGetSmallBlockResult(t *testing.T) {
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

func TestGetSmallBlockWithPackedTransactions(t *testing.T) {
	result, err := hex.DecodeString(blockdataForBlock4)
	require.NoError(t, err)

	out := &SignedBlock{}
	err = eos.UnmarshalBinary(result, &out)
	require.NoError(t, err)

}

// block with some failures
func TestBlock33Complete(t *testing.T) {
	result, err := hex.DecodeString(block33Block)
	require.NoError(t, err)

	out, err := ParseGetBlockResultV0(result)
	require.NoError(t, err)
	for _, tt := range out.Traces.AsTransactionTracesV0() {
		fmt.Println(tt.ID)
	}
}

func TestTableDeltasInBlock(t *testing.T) {
	result, err := hex.DecodeString("0120ac01000001ac201bfa206d4a92c1663fc9e8cc69e843ff91f49d396f03bf53761708011fac01000001ac1faccb177cfa87f069aa202b6cd1af243366c84ae510618bc8db8cd568012000000000000020f5cbb6126d5138940e5da68fe0671406e4b28cf96285e98061f0311a011f0000000000001fb376dd73c1676c67c3ad68c8e5c7806951964670adb80dc20f36e986000001a6030400106163636f756e745f6d6574616461746100000c636f6e74726163745f726f770101c501000000000000ea30550000000000ea3055000000004473686400000000447368640000000000ea30559a010000100000000000e8030000000008000c000000f40100001400000064000000400d0300e8030000f049020064000000100e00005802000080533b0000100000040006000000000010000000db1f000000000000a3040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e7265736f757263655f757361676501013b000000000000ea305500375df348a45c0d0000000000010000000000000000375df3485eee0800000000006500000000000000400a12000000000000157265736f757263655f6c696d6974735f73746174650101530000200000006949e27600000000cb0700000000000000200000002872165a000000004b0600000000000010270000000000001027000000000000db1f000000000000f780100000000000cf25030000000000")
	require.NoError(t, err)

	out, err := ParseGetBlockResultV0(result)
	require.NoError(t, err)

	ValidateBlock32TableDeltas(t, out.Deltas.AsTableDeltasV0())
}

func TestTableDeltasOnly(t *testing.T) {
	result, err := hex.DecodeString(`0400106163636f756e745f6d6574616461746100000c636f6e74726163745f726f770101c501000000000000ea30550000000000ea3055000000004473686400000000447368640000000000ea30559a010000100000000000e8030000000008000c000000f40100001400000064000000400d0300e8030000f049020064000000100e00005802000080533b0000100000040006000000000010000000db1f000000000000a3040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e7265736f757263655f757361676501013b000000000000ea305500375df348a45c0d0000000000010000000000000000375df3485eee0800000000006500000000000000400a12000000000000157265736f757263655f6c696d6974735f73746174650101530000200000006949e27600000000cb0700000000000000200000002872165a000000004b0600000000000010270000000000001027000000000000db1f000000000000f780100000000000cf25030000000000`)
	require.NoError(t, err)
	out := []*TableDelta{}
	err = eos.UnmarshalBinary(result, &out)

	var outV0 []*TableDeltaV0
	for _, d := range out {
		toV0, ok := d.Impl.(*TableDeltaV0)
		require.True(t, ok)
		outV0 = append(outV0, toV0)
	}
	ValidateBlock32TableDeltas(t, outV0)
}

func ValidateBlock32TableDeltas(t *testing.T, deltas []*TableDeltaV0) {
	t.Helper()

	require.Equal(t, 4, len(deltas))
	d1 := deltas[1]
	assert.Equal(t, "contract_row", d1.Name)
	assert.Equal(t, 1, len(d1.Rows))
	assert.Equal(t, mustDecodeString("000000000000ea30550000000000ea3055000000004473686400000000447368640000000000ea30559a010000100000000000e8030000000008000c000000f40100001400000064000000400d0300e8030000f049020064000000100e00005802000080533b0000100000040006000000000010000000db1f000000000000a3040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"), d1.Rows[0].Data)
}

func mustDecodeString(s string) []byte {
	out, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return out
}

func TestBlock5Complete(t *testing.T) {
	result, err := hex.DecodeString(block5Complete)
	require.NoError(t, err)

	out, err := ParseGetBlockResultV0(result)
	require.NoError(t, err)

	validateBlock5TransactionsTraces(t, out.Traces.AsTransactionTracesV0())
}

func TestBlock5TransactionTracesOnly(t *testing.T) {
	result, err := hex.DecodeString(block5Traces)
	require.NoError(t, err)

	out := []*TransactionTrace{}
	err = eos.UnmarshalBinary(result, &out)
	require.NoError(t, err)

	var outV0 []*TransactionTraceV0
	for _, tr := range out {
		asV0, ok := tr.Impl.(*TransactionTraceV0)
		require.True(t, ok)
		outV0 = append(outV0, asV0)
	}

	validateBlock5TransactionsTraces(t, outV0)
}

func validateBlock5TransactionsTraces(t *testing.T, traces []*TransactionTraceV0) {
	t.Helper()

	tr := traces[0]
	assert.Equal(t, "0b503a7c970b54f074949b6943090e2e8a1c19761e6365c5247b57693b6e446f", tr.ID.String())
	assert.Equal(t, eos.TransactionStatus(0x0), tr.Status)
	assert.Equal(t, eos.Int64(48), tr.Elapsed)
	assert.EqualValues(t, 100, tr.CPUUsageUS)
	assert.EqualValues(t, 0, tr.NetUsageWords)
	assert.Equal(t, false, tr.Scheduled)
	assert.Equal(t, 1, len(tr.ActionTraces))
	atr, ok := tr.ActionTraces[0].Impl.(*ActionTraceV0)
	require.True(t, ok)

	assert.EqualValues(t, "eosio", atr.Receiver)
	assert.EqualValues(t, "onblock", atr.Act.Name)
	assert.EqualValues(t, "eosio", atr.Act.Authorization[0].Actor)
	assert.EqualValues(t, "active", atr.Act.Authorization[0].Permission)
	assert.EqualValues(t, mustDecodeString("1b5df3480000000000ea30550000000000034a6b4b706a39e744906988c175f5e4913db2f604d5c7fbbe798861ccc1e4f101d7f75fbd3ee0c3281960f0d9de6acb9a55ffebbb574cb2831739c8de7e339e0eb9f19ebee872673876f8f49fb9364cb6f42ec9fe17e31c1cfa106a7d000000000000"), atr.Act.Data)
	assert.EqualValues(t, 18, atr.Elapsed)

	receipt, ok := atr.Receipt.Impl.(*ActionReceiptV0)
	require.True(t, ok)
	assert.EqualValues(t, 36, receipt.GlobalSequence)
	assert.EqualValues(t, "820b0ec15dd980cc3f8fd535dc997df5aedecaab244af62ed279a34b31aab2f2", receipt.ActDigest.String())
	assert.EqualValues(t, 25, receipt.AuthSequence[0].Sequence)
	assert.EqualValues(t, "eosio", receipt.AuthSequence[0].Account)
	assert.EqualValues(t, 34, receipt.RecvSequence)
	assert.EqualValues(t, 1, receipt.CodeSequence)
	assert.EqualValues(t, 1, receipt.ABISequence)
}
