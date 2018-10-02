package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewExpire is an action to expire a proposal ahead of its natural death.
func NewExpire(proposer eos.AccountName, proposalName eos.Name) *eos.Action {
	a := &eos.Action{
		Account: ForumAN,
		Name:    ActN("expire"),
		Authorization: []eos.PermissionLevel{
			{Actor: proposer, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Expire{
			ProposalName: proposalName,
		}),
	}
	return a
}

// Expire represents the `eosio.forum::propose` action.
type Expire struct {
	ProposalName eos.Name `json:"proposal_name"`
}
