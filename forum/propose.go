package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewPropose is an action to submit a proposal for vote.
func NewPropose(proposer eos.AccountName, proposalName eos.Name, title string, proposalJSON string, expiresAt eos.JSONTime) *eos.Action {
	a := &eos.Action{
		Account: ForumAN,
		Name:    ActN("propose"),
		Authorization: []eos.PermissionLevel{
			{Actor: proposer, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Propose{
			Proposer:     proposer,
			ProposalName: proposalName,
			Title:        title,
			ProposalJSON: proposalJSON,
			ExpiresAt:    expiresAt,
		}),
	}
	return a
}

// Propose represents the `eosio.forum::propose` action.
type Propose struct {
	Proposer     eos.AccountName `json:"proposer"`
	ProposalName eos.Name        `json:"proposal_name"`
	Title        string          `json:"title"`
	ProposalJSON string          `json:"proposal_json"`
	ExpiresAt    eos.JSONTime    `json:"expires_at"`
}
