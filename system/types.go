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

// SetABI represents the hard-coded `setabi` action.
type SetABI struct {
	Account eos.AccountName `json:"account"`
	ABI     eos.ABI         `json:"abi"`
}

// SetProds is present in `eosio.bios` contract. Used only at boot time.
type SetProds struct {
	Schedule []ProducerKey `json:"schedule"`
}

type ProducerKey struct {
	ProducerName    eos.AccountName `json:"producer_name"`
	BlockSigningKey ecc.PublicKey   `json:"block_signing_key"`
}

// EOSIOParameters are all the params that can be set on the system contract.
type EOSIOParameters struct {
	BasePerTransactionNetUsage     uint32 `json:"base_per_transaction_net_usage" yaml:"base_per_transaction_net_usage"`
	BasePerTransactionCPUUsage     uint32 `json:"base_per_transaction_cpu_usage" yaml:"base_per_transaction_cpu_usage"`
	BasePerActionCPUUsage          uint32 `json:"base_per_action_cpu_usage" yaml:"base_per_action_cpu_usage"`
	BaseSetcodeCPUUsage            uint32 `json:"base_setcode_cpu_usage" yaml:"base_setcode_cpu_usage"`
	PerSignatureCPUUsage           uint32 `json:"per_signature_cpu_usage" yaml:"per_signature_cpu_usage"`
	PerLockNetUsage                uint32 `json:"per_lock_net_usage" yaml:"per_lock_net_usage"`
	ContextFreeDiscountCPUUsageNum uint64 `json:"context_free_discount_cpu_usage_num" yaml:"context_free_discount_cpu_usage_num"`
	ContextFreeDiscountCPUUsageDen uint64 `json:"context_free_discount_cpu_usage_den" yaml:"context_free_discount_cpu_usage_den"`
	MaxTransactionCPUUsage         uint32 `json:"max_transaction_cpu_usage" yaml:"max_transaction_cpu_usage"`
	MaxTransactionNetUsage         uint32 `json:"max_transaction_net_usage" yaml:"max_transaction_net_usage"`

	MaxBlockCPUUsage       uint64 `json:"max_block_cpu_usage" yaml:"max_block_cpu_usage"`
	TargetBlockCPUUsagePct uint32 `json:"target_block_cpu_usage_pct" yaml:"target_block_cpu_usage_pct"` //< the target percent (1% == 100, 100%= 10,000) of maximum cpu usage; exceeding this triggers congestion handling
	MaxBblockNetUsage      uint64 `json:"max_block_net_usage" yaml:"max_block_net_usage"`               //< the maxiumum net usage in instructions for a block
	TargetBlockNetUsagePct uint32 `json:"target_block_net_usage_pct" yaml:"target_block_net_usage_pct"` //< the target percent (1% == 100, 100%= 10,000) of maximum net usage; exceeding this triggers congestion handling

	MaxTransactionLifetime       uint32 `json:"max_transaction_lifetime" yaml:"max_transaction_lifetime"`
	MaxTransactionExecTime       uint32 `json:"max_transaction_exec_time" yaml:"max_transaction_exec_time"`
	MaxAuthorityDepth            uint16 `json:"max_authority_depth" yaml:"max_authority_depth"`
	MaxInlineDepth               uint16 `json:"max_inline_depth" yaml:"max_inline_depth"`
	MaxInlineActionSize          uint32 `json:"max_inline_action_size" yaml:"max_inline_action_size"`
	MaxGeneratedTransactionCount uint32 `json:"max_generated_transaction_count" yaml:"max_generated_transaction_count"`

	// FIXME: does not appear in the `abi` for `eosio.system`.
	// MaxStorageSize uint64 `json:"max_storage_size" yaml:"max_storage_size"`
	PercentOfMaxInflationRate uint32 `json:"percent_of_max_inflation_rate" yaml:"percent_of_max_inflation_rate"`
	StorageReserveRatio       uint32 `json:"storage_reserve_ratio" yaml:"storage_reserve_ratio"`
}

type EOSIOGlobalState struct {
	EOSIOParameters
	TotalStorageBytesReserved uint64 `json:"total_storage_bytes_reserved"`
	TotalStorageStake         uint64 `json:"total_storage_stake"`
	PaymentPerBlock           uint64 `json:"payment_per_block"`
}

type DelegatedBandwidth struct {
	From      eos.AccountName `json:"from"`
	To        eos.AccountName `json:"to"`
	NetWeight eos.Asset       `json:"net_weight"`
	CPUWeight eos.Asset       `json:"cpu_weight"`
	// TODO: whooops, please review this..
	RAMBytes int64 `json:"ram_bytes"`
}

type TotalResources struct {
}

// DelegateBW represents the `eosio.system::delegatebw` action.
type DelegateBW struct {
	From         eos.AccountName `json:"from"`
	Receiver     eos.AccountName `json:"receiver"`
	StakeNet     eos.Asset       `json:"stake_net"`
	StakeCPU     eos.Asset       `json:"stake_cpu"`
	StakeStorage eos.Asset       `json:"stake_storage"`
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

// RegProducer represents the `eosio.system::regproducer` action
type RegProducer struct {
	Producer    eos.AccountName `json:"producer"`
	ProducerKey ecc.PublicKey   `json:"producer_key"`
	URL         string          `json:"url"`
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
