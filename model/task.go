package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuildConfig struct {
	ImageName  string
	ImageTag   string `bson:",omitempty"`
	Dockerfile string `bson:",omitempty"`
	RepoUrl    string
	RepoName   string
}

type DeployConfig struct {
	ImageName   string
	ImageTag    string `bson:",omitempty"`
	ServiceName string `bson:",omitempty"`
	MountSource string `bson:",omitempty"`
	MountTarget string `bson:",omitempty"`
	AutoRemove  bool
	Env         []string `bson:",omitempty"`
	HostPort    string   `bson:",omitempty"`
	ExposedPort string   `bson:",omitempty"`
	NetworkId   string   `bson:",omitempty"`
	NetworkName string   `bson:",omitempty"`
}

type RestartConfig struct {
	ServiceName string
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
	Config         interface{}
	Remarks        string
	AutoRun        bool
	Timeout        int64 // minutes
	Type           string
}

type UpdateTaskInputPayload struct {
	Name           *string
	UpstreamTaskId *primitive.ObjectID
	StreamWebhook  *string
	ScheduledAt    *primitive.DateTime
	Config         *interface{}
	Remarks        *string
	AutoRun        *bool
	Timeout        *int64
	Type           *string
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
	Config         interface{}
	UpstreamTaskId primitive.ObjectID `bson:",omitempty"`
	StreamWebhook  string
	AutoRun        bool
	Timeout        int64
	Type           string
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
