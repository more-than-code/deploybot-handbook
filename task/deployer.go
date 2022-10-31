package task

import (
	"context"
	"os"
	"os/exec"

	"github.com/more-than-code/deploybot/container"
	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Deployer struct {
	repo *repository.Repository
}

func NewDeployer() *Deployer {
	r, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}
	return &Deployer{repo: r}
}

func (d *Deployer) Start(cfg *model.DeployConfigPayload) error {
	if cfg.MountTarget != "" {
		path, _ := os.Getwd()

		sourceDir := path + "/data/" + cfg.ServiceName
		err := os.MkdirAll(sourceDir, 0644)

		if err != nil {
			return err
		}

		cfg.MountSource = sourceDir
		cfg.MountType = "bind"
		cfg.AutoRemove = true
	}

	helper := container.NewContainerHelper("unix:///var/run/docker.sock")

	err := helper.StartContainer(cfg)

	if cfg.MountTarget != "" {
		cmd := exec.Command("sudo", "chmod", "u+x", "-R", cfg.MountSource)
		err = cmd.Run()
	}
	return err
}

func (d *Deployer) UpdateTask(input *model.UpdateDeployTaskInput) (primitive.ObjectID, error) {
	return d.repo.UpdateDeployTask(context.TODO(), input)
}

func (d *Deployer) UpdateTaskStatus(input *model.UpdateDeployTaskStatusInput) error {
	return d.repo.UpdateDeployTaskStatus(context.TODO(), input)
}
