package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewPropose is an action to submit a proposal for vote.
func NewPropose(proposer eos.AccountName, proposalName eos.Name, title string, proposalJSON string) *eos.Action {
	a := &eos.Action{
		Account: AN("eosforumdapp"),
		Name:    ActN("propose"),
		Authorization: []eos.PermissionLevel{
			{Actor: proposer, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Propose{
			Proposer:     proposer,
			ProposalName: proposalName,
			Title:        title,
			ProposalJSON: proposalJSON,
		}),
	}
	return a
}

// Propose represents the `eosforumdapp::propose` action.
type Propose struct {
	Proposer     eos.AccountName `json:"proposer"`
	ProposalName eos.Name        `json:"proposal_name"`
	Title        string          `json:"title"`
	ProposalJSON string          `json:"proposal_json"`
}
