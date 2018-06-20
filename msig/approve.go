package msig

import (
	eos "github.com/eoscanada/eos-go"
)

// NewApprove returns a `approve` action that lives on the
// `eosio.msig` contract.
func NewApprove(proposer eos.AccountName, proposalName eos.Name, level eos.PermissionLevel) *eos.Action {
	return &eos.Action{
		Account:       eos.AccountName("eosio.msig"),
		Name:          eos.ActionName("approve"),
		Authorization: []eos.PermissionLevel{level},
		ActionData:    eos.NewActionData(Approve{proposer, proposalName, level}),
	}
}

type Approve struct {
	Proposer     eos.AccountName     `json:"proposer"`
	ProposalName eos.Name            `json:"proposal_name"`
	Level        eos.PermissionLevel `json:"level"`
}
