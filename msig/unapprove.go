package msig

import (
	eos "github.com/eoscanada/eos-go"
)

// NewUnapprove returns a `unapprove` action that lives on the
// `eosio.msig` contract.
func NewUnapprove(proposer eos.AccountName, proposalName eos.Name, level eos.PermissionLevel) *eos.Action {
	return &eos.Action{
		Account:       eos.AccountName("eosio.msig"),
		Name:          eos.ActionName("unapprove"),
		Authorization: []eos.PermissionLevel{level},
		ActionData:    eos.NewActionData(Unapprove{proposer, proposalName, level}),
	}
}

type Unapprove struct {
	Proposer     eos.AccountName     `json:"proposer"`
	ProposalName eos.Name            `json:"proposal_name"`
	Level        eos.PermissionLevel `json:"level"`
}
