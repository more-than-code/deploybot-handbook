package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PipelineConfig struct {
	RepoUrlWatched string
	AutoRun        bool
}

type Pipeline struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
	ExecutedAt  primitive.DateTime
	StoppedAt   primitive.DateTime
	ScheduledAt primitive.DateTime
	Status      string
	Config      PipelineConfig

	Tasks []*Task
}

type CreatePipelineInputPayload struct {
	Name string
}

type CreatePipelineInput struct {
	Payload CreatePipelineInputPayload
}

type GetPipelineInput struct {
	RepoWatched string
}

type UpdatePipelineInputPayload struct {
	Name        *string
	ScheduledAt *primitive.DateTime `bson:",omitempty"`
	Config      *PipelineConfig
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
