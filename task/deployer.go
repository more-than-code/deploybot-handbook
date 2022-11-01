package task

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

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

func (d *Deployer) Start(cfg model.DeployConfig) error {
	if cfg.ContainerConfig != nil {
		if cfg.ContainerConfig.MountTarget != "" {
			path, _ := os.Getwd()

			sourceDir := path + "/data/" + cfg.ContainerConfig.ServiceName
			// err := os.MkdirAll(sourceDir, 0644)

			// cmd := exec.Command("sudo", "mkdir", sourceDir)
			// output, err := cmd.Output()
			// log.Println(string(output))
			// if err != nil {
			// 	return err
			// }

			cfg.ContainerConfig.MountSource = sourceDir
			cfg.ContainerConfig.MountType = "bind"
			cfg.ContainerConfig.AutoRemove = true
		}

		helper := container.NewContainerHelper("unix:///var/run/docker.sock")

		err := helper.StartContainer(cfg.ContainerConfig)

		if err != nil {
			return err
		}

		if cfg.ContainerConfig.MountSource != "" {
			cmd := exec.Command("sudo", "chmod", "u+x", "-R", cfg.ContainerConfig.MountSource)
			output, err := cmd.Output()
			log.Println(string(output))
			if err != nil {
				return err
			}
		}
	}

	if cfg.PostInstall != "" {
		strs := strings.Split(cfg.PostInstall, " ")
		cmd := exec.Command(strs[0], strs[1:]...)
		output, err := cmd.Output()
		log.Println(string(output))
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Deployer) UpdateTask(input *model.UpdateDeployTaskInput) (primitive.ObjectID, error) {
	return d.repo.UpdateDeployTask(context.TODO(), input)
}

func (d *Deployer) UpdateTaskStatus(input *model.UpdateDeployTaskStatusInput) error {
	return d.repo.UpdateDeployTaskStatus(context.TODO(), input)
}
