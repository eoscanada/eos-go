package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewVote is an action representing a simple vote to be broadcast
// through the chain network.
func NewVote(voter, proposer eos.AccountName, proposalName eos.Name, proposalHash string, voteValue uint8, voteJSON string) *eos.Action {
	a := &eos.Action{
		Account: AN("eosforumdapp"),
		Name:    ActN("vote"),
		Authorization: []eos.PermissionLevel{
			{Actor: voter, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Vote{
			Voter:        voter,
			Proposer:     proposer,
			ProposalName: proposalName,
			ProposalHash: proposalHash,
			Vote:         voteValue,
			VoteJSON:     voteJSON,
		}),
	}
	return a
}

// Vote represents the `eosforumtest::vote` action.
type Vote struct {
	Voter        eos.AccountName `json:"voter"`
	Proposer     eos.AccountName `json:"proposer"`
	ProposalName eos.Name        `json:"proposal_name"`
	ProposalHash string          `json:"proposal_hash"`
	Vote         uint8           `json:"vote"`
	VoteJSON     string          `json:"vote_json"`
}
