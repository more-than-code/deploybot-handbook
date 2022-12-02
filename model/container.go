package model

import "github.com/docker/docker/api/types/mount"

type ContainerConfig struct {
	ImageName   string
	ImageTag    string        `bson:",omitempty"`
	ServiceName string        `bson:",omitempty"`
	Mounts      []mount.Mount `bson:",omitempty"`
	AutoRemove  bool          `bson:",omitempty"`
}
