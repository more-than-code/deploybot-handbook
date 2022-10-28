package task

import (
	"github.com/more-than-code/deploybot/container"
)

type DeployConfig struct {
	ImageName     string
	ImageTag      string
	ContainerName string
}

type DeployTask struct {
	cfg *DeployConfig
}

func NewDeployTask(cfg *DeployConfig) *DeployTask {
	return &DeployTask{cfg: cfg}
}

func (t *DeployTask) Start() error {
	helper := container.NewContainerHelper("unix:///var/run/docker.sock")

	return helper.StartContainer(t.cfg.ImageName+t.cfg.ImageTag, t.cfg.ContainerName)
}
