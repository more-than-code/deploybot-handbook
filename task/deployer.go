package task

import (
	"os"
	"os/exec"

	"github.com/more-than-code/deploybot/container"
	"github.com/more-than-code/deploybot/model"
)

type Deployer struct {
	cfg *model.DeployConfigPayload
}

func NewDeployer(cfg *model.DeployConfigPayload) *Deployer {
	return &Deployer{cfg: cfg}
}

func (d *Deployer) Start() error {
	if d.cfg.MountTarget != "" {
		path, _ := os.Getwd()

		sourceDir := path + "/data/" + d.cfg.ServiceName
		err := os.MkdirAll(sourceDir, 0644)

		if err != nil {
			return err
		}

		d.cfg.MountSource = sourceDir
		d.cfg.MountType = "bind"
		d.cfg.AutoRemove = true
	}

	helper := container.NewContainerHelper("unix:///var/run/docker.sock")

	err := helper.StartContainer(d.cfg)

	if d.cfg.MountTarget != "" {
		cmd := exec.Command("sudo", "chmod", "u+x", "-R", d.cfg.MountSource)
		err = cmd.Run()
	}
	return err
}
