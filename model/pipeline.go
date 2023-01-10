package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pipeline struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Name          string             `json:"name"`
	CreatedAt     primitive.DateTime `json:"createAt"`
	UpdatedAt     primitive.DateTime `json:"updateAt"`
	ExecutedAt    primitive.DateTime `json:"executedAt"`
	StoppedAt     primitive.DateTime `json:"stoppedAt"`
	ScheduledAt   primitive.DateTime `json:"scheduledAt"`
	Status        string             `json:"status"`
	Arguments     []string           `json:"argumentss"`
	Tasks         []Task             `json:"tasks"`
	RepoWatched   string             `json:"repoWatched"`
	BranchWatched string             `json:"branchWatched"`
	AutoRun       bool               `json:"autoRun"`
}

type CreatePipelineInputPayload struct {
	Name          string
	Arguments     []string
	RepoWatched   string
	BranchWatched string
	AutoRun       bool
}

type CreatePipelineInput struct {
	Payload CreatePipelineInputPayload
}

type TaskFilter struct {
	UpstreamTaskId *primitive.ObjectID
	AutoRun        *bool
}

type GetPipelineInput struct {
	Id         *primitive.ObjectID
	Name       *string
	TaskFilter TaskFilter
}

type GetPipelinesInput struct {
	RepoWatched   *string
	BranchWatched *string
	AutoRun       *bool
}

type GetPipelinesOutput struct {
	TotalCount int        `json:"totalCount"`
	Items      []Pipeline `json:"items"`
}

type UpdatePipelineInputPayload struct {
	Name          *string
	ScheduledAt   *primitive.DateTime `bson:",omitempty"`
	Arguments     []string
	RepoWatched   *string
	BranchWatched *string
	AutoRun       *bool
}
type UpdatePipelineInput struct {
	Id      primitive.ObjectID
	Payload UpdatePipelineInputPayload
}

type UpdatePipelineStatusInputPayload struct {
	Status string
}
type UpdatePipelineStatusInput struct {
	PipelineId primitive.ObjectID
	Payload    UpdatePipelineStatusInputPayload
}

func (p Pipeline) Id2Hex() string {
	return p.Id.Hex()
}

func (p Pipeline) CreatedAt2Str() string {
	return p.CreatedAt.Time().String()
}

func (p Pipeline) ExecutedAt2Str() string {
	if p.ExecutedAt == 0 {
		return ""
	}

	return p.ExecutedAt.Time().String()
}

func (p Pipeline) StoppedAt2Str() string {
	if p.StoppedAt == 0 {
		return ""
	}

	return p.StoppedAt.Time().String()
}
