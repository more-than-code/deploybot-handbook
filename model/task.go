package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeployConfigPayload struct {
	ImageName   string
	ImageTag    string `bson:",omitempty"`
	ServiceName string `bson:",omitempty"`
	MountSource string `bson:",omitempty"`
	MountTarget string `bson:",omitempty"`
	MountType   string `bson:",omitempty"`
	AutoRemove  bool   `bson:",omitempty"`
}

type DeployConfig struct {
	Webhook string
	Payload DeployConfigPayload
}

type DeployTask struct {
	Id          primitive.ObjectID `bson:"_id"`
	BuildTaskId primitive.ObjectID
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
	ExecutedAt  primitive.DateTime
	StoppedAt   primitive.DateTime
	ScheduledAt primitive.DateTime
	Status      string
	Config      DeployConfig
}

type UpdateDeployTaskInput struct {
	Id          primitive.ObjectID
	BuildTaskId primitive.ObjectID
	ScheduledAt primitive.DateTime `bson:",omitempty"`
	Config      DeployConfig
}

type DeployTaskStatusFilter struct {
	Option string
}

type DeployTasksInput struct {
	StatusFilter *DeployTaskStatusFilter
}

type UpdateDeployTaskStatusInput struct {
	DeployTaskId primitive.ObjectID
	Status       string
}

type BuildConfigPayload struct {
	RepoCloneUrl   string
	RepoName       string
	RepoUsername   string `bson:",omitempty"`
	RepoToken      string `bson:",omitempty"`
	ImageTagPrefix string
}

type BuildConfig struct {
	Webhook string
	Payload BuildConfigPayload
}

type BuildTask struct {
	Id          primitive.ObjectID `bson:"_id"`
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
	ExecutedAt  primitive.DateTime
	StoppedAt   primitive.DateTime
	ScheduledAt primitive.DateTime
	Status      string
	Config      BuildConfig
}

type UpdateBuildTaskInput struct {
	Id          primitive.ObjectID
	ScheduledAt primitive.DateTime `bson:",omitempty"`
	Config      BuildConfig
}

type UpdateBuildTaskStatusInput struct {
	BuildTaskId primitive.ObjectID
	Status      string
}

type BuildTaskStatusFilter struct {
	Option string
}

type BuildTasksInput struct {
	StatusFilter *BuildTaskStatusFilter
}
