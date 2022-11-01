package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ContainerConfig struct {
	ImageName   string
	ImageTag    string `bson:",omitempty"`
	ServiceName string `bson:",omitempty"`
	MountSource string `bson:",omitempty"`
	MountTarget string `bson:",omitempty"`
	MountType   string `bson:",omitempty"`
	AutoRemove  bool   `bson:",omitempty"`
}

type DeployConfig struct {
	Webhook         string
	PreInstall      string
	ContainerConfig *ContainerConfig
	PostInstall     string
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

func (d DeployTask) Id2Hex() string {
	return d.Id.Hex()
}

func (d DeployTask) BuildTaskId2Hex() string {
	return d.BuildTaskId.Hex()
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

type SourceConfig struct {
	RepoCloneUrl   string
	RepoName       string
	RepoUsername   string `bson:",omitempty"`
	RepoToken      string `bson:",omitempty"`
	ImageTagPrefix string
	Commits        []Commit `bson:",omitempty"`
}

type BuildConfig struct {
	Webhook      string
	SourceConfig SourceConfig
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

func (t BuildTask) Id2Hex() string {
	return t.Id.Hex()
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
