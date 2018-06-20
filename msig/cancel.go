package msig

import (
	eos "github.com/eoscanada/eos-go"
)

// NewCancel returns a `cancel` action that lives on the
// `eosio.msig` contract.
func NewCancel(proposer eos.AccountName, proposalName eos.Name, canceler eos.AccountName) *eos.Action {
	return &eos.Action{
		Account: eos.AccountName("eosio.msig"),
		Name:    eos.ActionName("cancel"),
		// TODO: double check in this package that the `Actor` is always the `proposer`..
		Authorization: []eos.PermissionLevel{
			{Actor: canceler, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Cancel{proposer, proposalName, canceler}),
	}
}

type Cancel struct {
	Proposer     eos.AccountName `json:"proposer"`
	ProposalName eos.Name        `json:"proposal_name"`
	Canceler     eos.AccountName `json:"canceler"`
}
