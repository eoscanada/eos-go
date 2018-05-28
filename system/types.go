package system

import (
	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

// SetPriv sets privileged account status. Used in the bios boot mechanism.
type SetPriv struct {
	Account eos.AccountName `json:"account"`
	IsPriv  bool            `json:"is_priv"`
}

// SetProds is present in `eosio.bios` contract. Used only at boot time.
type SetProds struct {
	Schedule []ProducerKey `json:"schedule"`
}

type ProducerKey struct {
	ProducerName    eos.AccountName `json:"producer_name"`
	BlockSigningKey ecc.PublicKey   `json:"block_signing_key"`
}

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

// UndelegateBW represents the `eosio.system::undelegatebw` action.
type UndelegateBW struct {
	From         eos.AccountName `json:"from"`
	Receiver     eos.AccountName `json:"receiver"`
	UnstakeNet   eos.Asset       `json:"unstake_net"`
	UnstakeCPU   eos.Asset       `json:"unstake_cpu"`
	UnstakeBytes uint64          `json:"unstake_bytes"`
}

// Refund represents the `eosio.system::refund` action
type Refund struct {
	Owner eos.AccountName `json:"owner"`
}

// UnregProducer represents the `eosio.system::unregprod` action
type UnregProducer struct {
	Producer eos.AccountName `json:"producer"`
}

// RegProxy represents the `eosio.system::regproxy` action
type RegProxy struct {
	Proxy eos.AccountName `json:"proxy"`
}

// UnregProxy represents the `eosio.system::unregproxy` action
type UnregProxy struct {
	Proxy eos.AccountName `json:"proxy"`
}

// VoteProducer represents the `eosio.system::voteproducer` action
type VoteProducer struct {
	Voter     eos.AccountName   `json:"voter"`
	Proxy     eos.AccountName   `json:"proxy"`
	Producers []eos.AccountName `json:"producers"`
}

// ClaimRewards repreents the `eosio.system::claimrewards` action
type ClaimRewards struct {
	Owner eos.AccountName `json:"owner"`
}

// Nonce represents the `eosio.system::nonce` action. It is used to
// add variability in a transaction, so you can send the same many
// times in the same block, without it having the same Tx hash.
type Nonce struct {
	Value string `json:"value"`
}
