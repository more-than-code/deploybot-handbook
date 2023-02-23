package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuildConfig struct {
	ImageName  string             `json:"imageName"`
	ImageTag   string             `json:"imageTag" bson:",omitempty"`
	Args       map[string]*string `json:"args" bson:",omitempty"`
	Dockerfile string             `json:"dockerfile" bson:",omitempty"`
	RepoUrl    string             `json:"repoUrl"`
	RepoName   string             `json:"repoName"`
}

type DeployConfig struct {
	ImageName   string   `json:"imageName"`
	ImageTag    string   `json:"imageTag" bson:",omitempty"`
	ServiceName string   `json:"serviceName" bson:",omitempty"`
	MountSource string   `json:"mountSource" bson:",omitempty"`
	MountTarget string   `json:"mountTarget" bson:",omitempty"`
	AutoRemove  bool     `json:"autoRemove"`
	Env         []string `json:"env" bson:",omitempty"`
	HostPort    string   `json:"hostPort" bson:",omitempty"`
	ExposedPort string   `json:"exposedPort" bson:",omitempty"`
	NetworkId   string   `json:"networkId" bson:",omitempty"`
	NetworkName string   `json:"networkName" bson:",omitempty"`
}

type RestartConfig struct {
	ServiceName string `json:"serviceName"`
}

type Task struct {
	Id             primitive.ObjectID `json:"id"`
	Name           string             `json:"name"`
	CreatedAt      primitive.DateTime `json:"createdAt"`
	UpdatedAt      primitive.DateTime `json:"updatedAt"`
	ExecutedAt     primitive.DateTime `json:"executedAt"`
	StoppedAt      primitive.DateTime `json:"stoppedAt"`
	ScheduledAt    primitive.DateTime `json:"scheduledAt"`
	Status         string             `json:"status"`
	UpstreamTaskId primitive.ObjectID `json:"upstreamTaskId" bson:",omitempty"`
	StreamWebhook  string             `json:"streamWebhook" bson:",omitempty"`
	Config         interface{}        `json:"config"`
	Remarks        string             `json:"remarks"`
	AutoRun        bool               `json:"autoRun"`
	Timeout        int64              `json:"timeout"` // minutes
	Type           string             `json:"type"`
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
