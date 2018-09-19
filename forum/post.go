package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewPost is an action representing a simple message to be posted
// through the chain network.
func NewPost(poster eos.AccountName, postUUID, content string, replyToPoster eos.AccountName, replyToPostUUID string, certify bool, jsonMetadata string) *eos.Action {
	a := &eos.Action{
		Account: ForumAN,
		Name:    ActN("post"),
		Authorization: []eos.PermissionLevel{
			{Actor: poster, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Post{
			Poster:          poster,
			PostUUID:        postUUID,
			Content:         content,
			ReplyToPoster:   replyToPoster,
			ReplyToPostUUID: replyToPostUUID,
			Certify:         certify,
			JSONMetadata:    jsonMetadata,
		}),
	}
	return a
}

// Post represents the `eosio.forum::post` action.
type Post struct {
	Poster          eos.AccountName `json:"poster"`
	PostUUID        string          `json:"post_uuid"`
	Content         string          `json:"content"`
	ReplyToPoster   eos.AccountName `json:"reply_to_poster"`
	ReplyToPostUUID string          `json:"reply_to_post_uuid"`
	Certify         bool            `json:"certify"`
	JSONMetadata    string          `json:"json_metadata"`
}
