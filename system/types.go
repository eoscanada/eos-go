package system

// BlockchainParameters are all the params we can set through `setparams`.
type BlockchainParameters struct {
	MaxBlockNetUsage               uint64 `json:"max_block_net_usage"`
	TargetBlockNetUsagePct         uint32 `json:"target_block_net_usage_pct"`
	MaxTransactionNetUsage         uint32 `json:"max_transaction_net_usage"`
	BasePerTransactionNetUsage     uint32 `json:"base_per_transaction_net_usage"`
	NetUsageLeeway                 uint32 `json:"net_usage_leeway"`
	ContextFreeDiscountNetUsageNum uint32 `json:"context_free_discount_net_usage_num"`
	ContextFreeDiscountNetUsageDen uint32 `json:"context_free_discount_net_usage_den"`
	MaxBlockCPUUsage               uint32 `json:"max_block_cpu_usage"`
	TargetBlockCPUUsagePct         uint32 `json:"target_block_cpu_usage_pct"`
	MaxTransactionCPUUsage         uint32 `json:"max_transaction_cpu_usage"`
	MinTransactionCPUUsage         uint32 `json:"min_transaction_cpu_usage"`
	MaxTransactionLifetime         uint32 `json:"max_transaction_lifetime"`
	DeferredTrxExpirationWindow    uint32 `json:"deferred_trx_expiration_window"`
	MaxTransactionDelay            uint32 `json:"max_transaction_delay"`
	MaxInlineActionSize            uint32 `json:"max_inline_action_size"`
	MaxInlineActionDepth           uint16 `json:"max_inline_action_depth"`
	MaxAuthorityDepth              uint16 `json:"max_authority_depth"`
	MaxGeneratedTransactionCount   uint32 `json:"max_generated_transaction_count"`

	// replace-regexp \(\w\)_\(\w\) -> \1\,(upcase \2)
	// then Cpu -> CPU
}

type EOSIOGlobalState struct {
	BlockchainParameters
	TotalStorageBytesReserved uint64 `json:"total_storage_bytes_reserved"`
	TotalStorageStake         uint64 `json:"total_storage_stake"`
	PaymentPerBlock           uint64 `json:"payment_per_block"`
}

// Nonce represents the `eosio.system::nonce` action. It is used to
// add variability in a transaction, so you can send the same many
// times in the same block, without it having the same Tx hash.
type Nonce struct {
	Value string `json:"value"`
}
