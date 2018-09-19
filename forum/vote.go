package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewVote is an action representing a simple vote to be broadcast
// through the chain network.
func NewVote(voter eos.AccountName, proposalName eos.Name, voteValue uint8, voteJSON string) *eos.Action {
	a := &eos.Action{
		Account: ForumAN,
		Name:    ActN("vote"),
		Authorization: []eos.PermissionLevel{
			{Actor: voter, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Vote{
			Voter:        voter,
			ProposalName: proposalName,
			Vote:         voteValue,
			VoteJSON:     voteJSON,
		}),
	}
	return a
}

// Vote represents the `eosio.forum::vote` action.
type Vote struct {
	Voter        eos.AccountName `json:"voter"`
	ProposalName eos.Name        `json:"proposal_name"`
	Vote         uint8           `json:"vote"`
	VoteJSON     string          `json:"vote_json"`
}
