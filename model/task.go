package model

import (
	"github.com/docker/docker/api/types/mount"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContainerConfig struct {
	ImageName   string
	ImageTag    string        `bson:",omitempty"`
	ServiceName string        `bson:",omitempty"`
	Mounts      []mount.Mount `bson:",omitempty"`
	AutoRemove  bool          `bson:",omitempty"`
}

type DeployConfig struct {
	Webhook         string
	Script          string
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

func (t DeployTask) Id2Hex() string {
	return t.Id.Hex()
}

func (t DeployTask) BuildTaskId2Hex() string {
	return t.BuildTaskId.Hex()
}

func (t DeployTask) CreatedAt2Str() string {
	return t.CreatedAt.Time().String()
}

func (t DeployTask) ExecutedAt2Str() string {
	if t.ExecutedAt == 0 {
		return ""
	}

	return t.ExecutedAt.Time().String()
}

func (t DeployTask) StoppedAt2Str() string {
	if t.StoppedAt == 0 {
		return ""
	}

	return t.StoppedAt.Time().String()
}

type UpdateDeployTaskInputPayload struct {
	BuildTaskId primitive.ObjectID `bson:",omitempty"`
	ScheduledAt primitive.DateTime `bson:",omitempty"`
	Config      DeployConfig
}

type UpdateDeployTaskInput struct {
	Id      primitive.ObjectID
	Payload UpdateDeployTaskInputPayload
}

type DeployTaskStatusFilter struct {
	Option string
}

type DeployTasksInput struct {
	StatusFilter *DeployTaskStatusFilter
}

type UpdateDeployTaskStatusInputPayload struct {
	Status string
}

type UpdateDeployTaskStatusInput struct {
	DeployTaskId primitive.ObjectID
	Payload      UpdateDeployTaskStatusInputPayload
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
	Script       string
	SourceConfig SourceConfig
}

type BuildTask struct {
	Id           primitive.ObjectID `bson:"_id"`
	DeployTaskId primitive.ObjectID
	CreatedAt    primitive.DateTime
	UpdatedAt    primitive.DateTime
	ExecutedAt   primitive.DateTime
	StoppedAt    primitive.DateTime
	ScheduledAt  primitive.DateTime
	Status       string
	Config       BuildConfig
}

func (t BuildTask) Id2Hex() string {
	return t.Id.Hex()
}

func (t BuildTask) DeployTaskId2Hex() string {
	return t.DeployTaskId.Hex()
}

func (t BuildTask) CreatedAt2Str() string {
	return t.CreatedAt.Time().String()
}

func (t BuildTask) ExecutedAt2Str() string {
	if t.ExecutedAt == 0 {
		return ""
	}

	return t.ExecutedAt.Time().String()
}

func (t BuildTask) StoppedAt2Str() string {
	if t.StoppedAt == 0 {
		return ""
	}

	return t.StoppedAt.Time().String()
}

type UpdateBuildTaskInputPayload struct {
	ScheduledAt primitive.DateTime `bson:",omitempty"`
	Config      BuildConfig
}

type UpdateBuildTaskInput struct {
	Id      primitive.ObjectID
	Payload UpdateBuildTaskInputPayload
}

type UpdateBuildTaskStatusInputPayload struct {
	DeployTaskId primitive.ObjectID `bson:",omitempty"`
	Status       string
}

type UpdateBuildTaskStatusInput struct {
	BuildTaskId primitive.ObjectID
	Payload     UpdateBuildTaskStatusInputPayload
}

type BuildTaskStatusFilter struct {
	Option string
}

type BuildTasksInput struct {
	StatusFilter *BuildTaskStatusFilter
}
