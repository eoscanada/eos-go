package msig

import (
	eos "github.com/eoscanada/eos-go"
)

// NewExec returns a `exec` action that lives on the
// `eosio.msig` contract.
func NewExec(proposer eos.AccountName, proposalName eos.Name, executer eos.AccountName) *eos.Action {
	return &eos.Action{
		Account: eos.AccountName("eosio.msig"),
		Name:    eos.ActionName("exec"),
		// TODO: double check in this package that the `Actor` is always the `proposer`..
		Authorization: []eos.PermissionLevel{
			{Actor: executer, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Exec{proposer, proposalName, executer}),
	}
}

type Exec struct {
	Proposer     eos.AccountName `json:"proposer"`
	ProposalName eos.Name        `json:"proposal_name"`
	Executer     eos.AccountName `json:"executer"`
}
