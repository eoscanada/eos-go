package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewVote is an action representing a simple vote to be broadcast
// through the chain network.
func NewVote(voter eos.AccountName, proposition, propositionHash, voteValue string) *eos.Action {
	a := &eos.Action{
		Account: AN("eosforumtest"),
		Name:    ActN("vote"),
		Authorization: []eos.PermissionLevel{
			{Actor: voter, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Vote{
			Voter:           voter,
			Proposition:     proposition,
			PropositionHash: propositionHash,
			VoteValue:       voteValue,
		}),
	}
	return a
}

// Vote represents the `eosforumtest::vote` action.
type Vote struct {
	Voter           eos.AccountName `json:"voter"`
	Proposition     string          `json:"proposition"`
	PropositionHash string          `json:"proposition_hash"`
	VoteValue       string          `json:"vote_value"`
}
