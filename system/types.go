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
