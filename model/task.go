package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskRunConfig struct {
	Script  string
	Secrets []string `bson:",omitempty"`
}

type Task struct {
	Id             primitive.ObjectID
	Name           string
	CreatedAt      primitive.DateTime
	UpdatedAt      primitive.DateTime
	ExecutedAt     primitive.DateTime
	StoppedAt      primitive.DateTime
	ScheduledAt    primitive.DateTime
	Status         string
	UpstreamTaskId primitive.ObjectID `bson:",omitempty"`
	StreamWebhook  string             `bson:",omitempty"`
	Config         TaskRunConfig
	Remarks        string
	AutoRun        bool
}

type UpdateTaskInputPayload struct {
	Name           *string
	UpstreamTaskId *primitive.ObjectID
	StreamWebhook  *string
	ScheduledAt    *primitive.DateTime
	Config         *TaskRunConfig
	Remarks        *string
	AutoRun        *bool
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
	Id             primitive.ObjectID
	Name           string
	ScheduledAt    primitive.DateTime `bson:",omitempty"`
	Config         TaskRunConfig
	UpstreamTaskId primitive.ObjectID `bson:",omitempty"`
	StreamWebhook  string
	AutoRun        bool
}
type CreateTaskInput struct {
	PipelineId primitive.ObjectID
	Payload    CreateTaskInputPayload
}

type GetTaskInput struct {
	PipelineId primitive.ObjectID
	Id         primitive.ObjectID
}

type GetTasksInput struct {
	PipelineId     primitive.ObjectID
	UpstreamTaskId *primitive.ObjectID
}

type DeleteTaskInput struct {
	PipelineId primitive.ObjectID
	TaskId     primitive.ObjectID
}

func (t Task) Id2Hex() string {
	return t.Id.Hex()
}

func (t Task) UpstreamtaskId2Hex() string {
	if t.UpstreamTaskId.IsZero() {
		return ""
	}

	return t.UpstreamTaskId.Hex()
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
