package msig

import (
	eos "github.com/eoscanada/eos-go"
)

// NewPropose returns a `propose` action that lives on the
// `eosio.msig` contract.
func NewPropose(proposer eos.AccountName, proposalName eos.Name, requested []eos.PermissionLevel, transaction *eos.Transaction) *eos.Action {
	return &eos.Action{
		Account: eos.AccountName("eosio.msig"),
		Name:    eos.ActionName("propose"),
		Authorization: []eos.PermissionLevel{
			{Actor: proposer, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Propose{proposer, proposalName, requested, transaction}),
	}
}

type Propose struct {
	Proposer     eos.AccountName       `json:"proposer"`
	ProposalName eos.Name              `json:"proposal_name"`
	Requested    []eos.PermissionLevel `json:"requested"`
	Transaction  *eos.Transaction      `json:"trx"`
}
