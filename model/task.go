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
	Payload *DeployConfigPayload
}

type DeployTask struct {
	Id          primitive.ObjectID `bson:"_id"`
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
	ExecutedAt  primitive.DateTime
	ScheduledAt primitive.DateTime
	Status      string
	Config      *DeployConfig
}

type UpdateDeployTaskInput struct {
	Id          primitive.ObjectID
	ScheduledAt primitive.DateTime `bson:",omitempty"`
	Config      *DeployConfig
}

type DeployStatusFilter struct {
	Option string
}

type DeployTasksInput struct {
	StatusFilter *DeployStatusFilter
}
