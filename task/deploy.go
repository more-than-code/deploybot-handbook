package task

import (
	"os"

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
	if t.cfg.MountTarget != "" {
		path, _ := os.Getwd()

		sourceDir := path + "/data/" + t.cfg.ContainerName
		err := os.MkdirAll(sourceDir, 0744)

		if err != nil {
			return err
		}

		t.cfg.MountSource = sourceDir
		t.cfg.MountType = "bind"
		t.cfg.AutoRemove = true
	}

	helper := container.NewContainerHelper("unix:///var/run/docker.sock")

	return helper.StartContainer(t.cfg)
}
