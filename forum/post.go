package forum

import (
	eos "github.com/eoscanada/eos-go"
)

// NewPost is an action representing a simple message to be posted
// through the chain network.
func NewPost(account eos.AccountName, postUUID, title, content string, replyToAccount eos.AccountName, replyToPostUUID string, certify bool, jsonMetadata string) *eos.Action {
	a := &eos.Action{
		Account: AN("eosforumtest"),
		Name:    ActN("post"),
		Authorization: []eos.PermissionLevel{
			{Actor: account, Permission: eos.PermissionName("active")},
		},
		ActionData: eos.NewActionData(Post{
			Account:         account,
			PostUUID:        postUUID,
			Title:           title,
			Content:         content,
			ReplyToAccount:  replyToAccount,
			ReplyToPostUUID: replyToPostUUID,
			Certify:         certify,
			JSONMetadata:    jsonMetadata,
		}),
	}
	return a
}

// Post represents the `eosforumtest::post` action.
type Post struct {
	Account         eos.AccountName `json:"account"`
	PostUUID        string          `json:"post_uuid"`
	Title           string          `json:"title"`
	Content         string          `json:"content"`
	ReplyToAccount  eos.AccountName `json:"reply_to_account"`
	ReplyToPostUUID string          `json:"reply_to_post_uuid"`
	Certify         bool            `json:"certify"`
	JSONMetadata    string          `json:"json_metadata"`
}
