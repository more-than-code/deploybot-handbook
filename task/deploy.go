package task

import (
	"github.com/more-than-code/deploybot/container"
	"github.com/more-than-code/deploybot/model"
)

type DeployTask struct {
	cfg *model.DeployConfig
}

func NewDeployTask(cfg *model.DeployConfig) *DeployTask {
	return &DeployTask{cfg: cfg}
}

func (t *DeployTask) Start() error {
	helper := container.NewContainerHelper("unix:///var/run/docker.sock")

	return helper.StartContainer(t.cfg)
}
