package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskConfig struct {
	UpstreamTaskId    primitive.ObjectID `bson:",omitempty"`
	DownstreamTaskId  primitive.ObjectID `bson:",omitempty"`
	UpstreamWebhook   string             `bson:",omitempty"`
	DownstreamWebhook string             `bson:",omitempty"`
	Script            string
}

type Task struct {
	Id          primitive.ObjectID
	Name        string
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
	ExecutedAt  primitive.DateTime
	StoppedAt   primitive.DateTime
	ScheduledAt primitive.DateTime
	Status      string
	Config      TaskConfig
}

func (t Task) Id2Hex() string {
	return t.Id.Hex()
}

func (t Task) UpstreamtaskId2Hex() string {
	if t.Config.UpstreamTaskId.IsZero() {
		return ""
	}

	return t.Config.UpstreamTaskId.Hex()
}

func (t Task) DownstreamtaskId2Hex() string {
	if t.Config.DownstreamTaskId.IsZero() {
		return ""
	}

	return t.Config.DownstreamTaskId.Hex()
}

func (t Task) CreatedAt2Str() string {
	return t.CreatedAt.Time().String()
}

func (t Task) ExecutedAt2Str() string {
	if t.ExecutedAt == 0 {
		return ""
	}

	return t.ExecutedAt.Time().String()
}

func (t Task) StoppedAt2Str() string {
	if t.StoppedAt == 0 {
		return ""
	}

	return t.StoppedAt.Time().String()
}

type UpdatePipelineTaskInputPayload struct {
	ScheduledAt primitive.DateTime `bson:",omitempty"`
	Config      TaskConfig
}

type CreatePipelineTaskInputPayload struct {
	Id          primitive.ObjectID
	ScheduledAt primitive.DateTime `bson:",omitempty"`
	Config      TaskConfig
}

type UpdatePipelineTaskStatusInputPayload struct {
	Status string
}

type GetPipelineTaskInput struct {
	PipelineId primitive.ObjectID
	Id         primitive.ObjectID
}

type Pipeline struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
	Tasks     []*Task
}

type CreatePipelineInputPayload struct {
	Name string
}

type CreatePipelineInput struct {
	Payload CreatePipelineInputPayload
}

type CreatePipelineTaskInput struct {
	PipelineId primitive.ObjectID
	Payload    CreatePipelineTaskInputPayload
}

type UpdatePipelineTaskInput struct {
	PipelineId primitive.ObjectID
	Id         primitive.ObjectID
	Payload    UpdatePipelineTaskInputPayload
}

type DeletePipelineTaskInput struct {
	PipelineId primitive.ObjectID
	TaskId     primitive.ObjectID
}

type UpdatePipelineTaskStatusInput struct {
	PipelineId primitive.ObjectID
	TaskId     primitive.ObjectID
	Payload    UpdatePipelineTaskStatusInputPayload
}

type GetPipelineInput struct {
	Name string
}
