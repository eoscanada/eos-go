package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewUnVote is an action representing the action to undoing a current vote
func NewUnVote(voter eos.AccountName, proposalName eos.Name) *eos.Action {
	a := &eos.Action{
		Account: ForumAN,
		Name:    ActN("unvote"),
		Authorization: []eos.PermissionLevel{
			{Actor: voter, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(UnVote{
			Voter:        voter,
			ProposalName: proposalName,
		}),
	}
	return a
}

// UnVote represents the `eosio.forum::unvote` action.
type UnVote struct {
	Voter        eos.AccountName `json:"voter"`
	ProposalName eos.Name        `json:"proposal_name"`
}
