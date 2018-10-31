package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewUnPost is an action undoing a post that is active
func NewUnPost(poster eos.AccountName, postUUID string) *eos.Action {
	a := &eos.Action{
		Account: ForumAN,
		Name:    ActN("unpost"),
		Authorization: []eos.PermissionLevel{
			{Actor: poster, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(UnPost{
			Poster:   poster,
			PostUUID: postUUID,
		}),
	}
	return a
}

// UnPost represents the `eosio.forum::unpost` action.
type UnPost struct {
	Poster   eos.AccountName `json:"poster"`
	PostUUID string          `json:"post_uuid"`
}
