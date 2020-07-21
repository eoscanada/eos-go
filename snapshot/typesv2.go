package snapshot

import (
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type GenesisState struct {
	//InitialConfiguration ChainConfig //eos.ChainConfig
	InitialTimestamp eos.TimePoint
	InitialKey       ecc.PublicKey
}

type ChainConfig struct {
	MaxBlockNetUsage               eos.Uint64 ///< the maxiumum net usage in instructions for a block
	TargetBlockNetUsagePct         uint32     ///< the target percent (1% == 100, 100%= 10,000) of maximum net usage; exceeding this triggers congestion handling
	MaxTransactionNetUsage         uint32     ///< the maximum objectively measured net usage that the chain will allow regardless of account limits
	BasePerTransactionNetUsage     uint32     ///< the base amount of net usage billed for a transaction to cover incidentals
	NetUsageLeeway                 uint32
	ContextFreeDiscountNetUsageNum uint32 ///< the numerator for the discount on net usage of context-free data
	ContextFreeDiscountNetUsageDen uint32 ///< the denominator for the discount on net usage of context-free data

	MaxBlockCpuUsage       uint32 ///< the maxiumum billable cpu usage (in microseconds) for a block
	TargetBlockCpuUsagePct uint32 ///< the target percent (1% == 100, 100%= 10,000) of maximum cpu usage; exceeding this triggers congestion handling
	MaxTransactionCpuUsage uint32 ///< the maximum billable cpu usage (in microseconds) that the chain will allow regardless of account limits
	MinTransactionCpuUsage uint32 ///< the minimum billable cpu usage (in microseconds) that the chain requires

	MaxTransactionLifetime      uint32 ///< the maximum number of seconds that an input transaction's expiration can be ahead of the time of the block in which it is first included
	DeferredTrxExpirationWindow uint32 ///< the number of seconds after the time a deferred transaction can first execute until it expires
	MaxTransactionDelay         uint32 ///< the maximum number of seconds that can be imposed as a delay requirement by authorization checks
	MaxInlineActionSize         uint32 ///< maximum allowed size (in bytes) of an inline action
	MaxInlineActionDepth        uint16 ///< recursion depth limit on sending inline actions
	MaxAuthorityDepth           uint16 ///< recursion depth limit for checking if an authority is satisfied
}

func (section *Section) readGenesisState(f callbackFunc) error {
	// // THIS SEEMS TO EXIST ONLY IN VERSION 2 OF THE SNAPSHOT FILE FORMAT.
	// // FOR NOW, WE ARE CONCENTRATING ON VERSION 3 (latest)
	// cnt := make([]byte, section.BufferSize)
	// _, err := section.Buffer.Read(cnt)
	// require.NoError(t, err)

	// var state GenesisState
	// assert.NoError(t, eos.UnmarshalBinary(cnt, &state))
	// cnt, _ = json.MarshalIndent(state, "  ", "  ")
	// fmt.Println(string(cnt))
	return nil
}
