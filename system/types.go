package system

import (
	eos "github.com/eosioca/eosapi"
	"github.com/eosioca/eosapi/ecc"
)

// UpdateAuth represents the hard-coded `updateauth` action.
//
// If you change the `active` permission, `owner` is the required parent.
//
// If you change the `owner` permission, there should be no parent.
type UpdateAuth struct {
	Account    eos.AccountName    `json:"account"`
	Permission eos.PermissionName `json:"permission"`
	Parent     eos.PermissionName `json:"parent"`
	Data       eos.Authority      `json:"data"`
	Delay      uint32             `json:"delay"` // this represents what exactly?
}

// SetPriv sets privileged account status. Used in the bios boot mechanism.
type SetPriv struct {
	Account eos.AccountName `json:"account"`
	IsPriv  bool            `json:"is_priv"`
}

// SetCode represents the hard-coded `setcode` action.
type SetCode struct {
	Account   eos.AccountName `json:"account"`
	VMType    byte            `json:"vmtype"`
	VMVersion byte            `json:"vmversion"`
	Code      eos.HexBytes    `json:"bytes"`
}

// SetABI represents the hard-coded `setabi` action.
type SetABI struct {
	Account eos.AccountName `json:"account"`
	ABI     eos.ABI         `json:"abi"`
}

// SetProds is present in `eosio.bios` contract. Used only at boot time.
type SetProds struct {
	Version   uint32        `json:"version"`
	Producers []ProducerKey `json:"producers"`
}

type ProducerKey struct {
	ProducerName    eos.AccountName `json:"producer_name"`
	BlockSigningKey ecc.PublicKey   `json:"block_signing_key"`
}

// NewAccount represents the `newaccount` on the `eosio.system` contract / hard-coded in the chain.
type NewAccount struct {
	Creator  eos.AccountName `json:"creator"`
	Name     eos.AccountName `json:"name"`
	Owner    eos.Authority   `json:"owner"`
	Active   eos.Authority   `json:"active"`
	Recovery eos.Authority   `json:"recovery"`
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
	Producer    eos.AccountName     `json:"producer"`
	ProducerKey []byte              `json:"producer_key"`
	Prefs       eos.EOSIOParameters `json:"eosio_parameters"`
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
