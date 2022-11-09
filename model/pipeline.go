package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pipeline struct {
	Id            primitive.ObjectID `bson:"_id"`
	Name          string
	CreatedAt     primitive.DateTime
	UpdatedAt     primitive.DateTime
	ExecutedAt    primitive.DateTime
	StoppedAt     primitive.DateTime
	ScheduledAt   primitive.DateTime
	Status        string
	Arguments     []string
	Tasks         []Task
	RepoWatched   string
	BranchWatched string
	AutoRun       bool
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
