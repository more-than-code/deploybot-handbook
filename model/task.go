package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskConfig struct {
	UpstreamTaskId    primitive.ObjectID `bson:",omitempty"`
	DownstreamTaskId  primitive.ObjectID `bson:",omitempty"`
	UpstreamWebhook   string             `bson:",omitempty"`
	DownstreamWebhook string             `bson:",omitempty"`
	Script            string
	Secrets           []string `bson:",omitempty"`
	AutoRun           bool
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
	Remarks     string
}

type UpdateTaskInputPayload struct {
	Name        *string
	ScheduledAt *primitive.DateTime
	Config      *TaskConfig
	Remarks     *string
}
type UpdateTaskInput struct {
	PipelineId primitive.ObjectID
	Id         primitive.ObjectID
	Payload    UpdateTaskInputPayload
}

type UpdateTaskStatusInputPayload struct {
	Status string
}
type UpdateTaskStatusInput struct {
	PipelineId primitive.ObjectID
	TaskId     primitive.ObjectID
	Payload    UpdateTaskStatusInputPayload
}

type CreateTaskInputPayload struct {
	Id          primitive.ObjectID
	Name        string
	ScheduledAt primitive.DateTime `bson:",omitempty"`
	Config      TaskConfig
}
type CreateTaskInput struct {
	PipelineId primitive.ObjectID
	Payload    CreateTaskInputPayload
}

type GetTaskInput struct {
	PipelineId primitive.ObjectID
	Id         primitive.ObjectID
}

type DeleteTaskInput struct {
	PipelineId primitive.ObjectID
	TaskId     primitive.ObjectID
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
