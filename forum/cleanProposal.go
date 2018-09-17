package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// CleanProposal is an action to flush proposal and allow RAM used by it.
func NewCleanProposal(cleaner eos.AccountName, proposalName eos.Name, maxCount uint64) *eos.Action {
	a := &eos.Action{
		Account: ForumAN,
		Name:    ActN("clnproposal"),
		Authorization: []eos.PermissionLevel{
			{Actor: cleaner, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(CleanProposal{
			ProposalName: proposalName,
			MaxCount:     maxCount,
		}),
	}
	return a
}

// CleanProposal represents the `eosio.forum::clnproposal` action.
type CleanProposal struct {
	ProposalName eos.Name `json:"proposal_name"`
	MaxCount     uint64   `json:"max_count"`
}
