package system

import eos "github.com/eosioca/eosapi"

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
	IsPriv  int8            `json:"is_priv"`
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
