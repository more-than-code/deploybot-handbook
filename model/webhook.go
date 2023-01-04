package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Commit struct {
	Id        string
	TreeId    string `json:"tree_id"`
	Distinct  bool
	Message   string
	Timestamp string
	Url       string
	Author    struct {
		Name  string
		Email string
	}
	Committer struct {
		Name     string
		Email    string
		Username string
	}
	Added    []string
	Removed  []string
	Modified []string
}

type GitHubHookshot struct {
	Ref        string
	Repository struct {
		Id       int
		NodeId   string `json:"node_id"`
		Name     string
		FullName string `json:"full_name"`
		Private  bool
		Owner    struct {
			Name       string
			Email      string
			Login      string
			Id         int
			NodeId     string `json:"node_id"`
			AvatarUrl  string `json:"avatar_url"`
			GravatarId string `json:"gravatar_id"`
			Url        string
			HtmlUrl    string `json:"html_url"`
			Type       string
			SiteAdmin  bool `json:"site_admin"`
		}
		HtmlUrl       string `json:"html_url"`
		Url           string
		CreatedAt     uint64 `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
		PushedAt      uint64 `json:"pushed_at"`
		CloneUrl      string `json:"clone_url"`
		Size          int
		Language      string
		DefaultBranch string `json:"default_branch"`
		MasterBranch  string `json:"master_branch"`
		Organization  string
	}
	Pusher struct {
		Name  string
		Email string
	}
	Organization struct {
		Login     string
		Id        int
		NodeId    string `json:"node_id"`
		Url       string
		AvatarUrl string `json:"avatar_url"`
	}
	Sender struct {
		Login     string
		Id        int
		NodeId    string `json:"node_id"`
		AvatarUrl string `json:"avatar_url"`
		Url       string
		HtmlUrl   string `json:"html_url"`
		Type      string
		SiteAdmin bool `json:"site_admin"`
	}
	Commits    []Commit
	HeadCommit Commit `json:"head_commit"`
}

type StreamWebhookPayload struct {
	PipelineId primitive.ObjectID
	TaskId     primitive.ObjectID
	Arguments  []string
}

type StreamWebhook struct {
	Payload StreamWebhookPayload
}
