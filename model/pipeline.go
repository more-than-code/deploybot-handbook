package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pipeline struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
	ExecutedAt  primitive.DateTime
	StoppedAt   primitive.DateTime
	ScheduledAt primitive.DateTime
	Status      string

	Tasks []*Task
}

type CreatePipelineInputPayload struct {
	Name string
}

type CreatePipelineInput struct {
	Payload CreatePipelineInputPayload
}

type GetPipelineInput struct {
	Name string
}

type UpdatePipelineInputPayload struct {
	Name        string
	ScheduledAt primitive.DateTime `bson:",omitempty"`
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
