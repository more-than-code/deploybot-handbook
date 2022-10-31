package task

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/more-than-code/deploybot/container"
	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/repository"
	"github.com/more-than-code/deploybot/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Builder struct {
	repo *repository.Repository
}

func NewBuilder() *Builder {
	r, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}
	return &Builder{repo: r}
}

func (t *Builder) Start(cfg *model.BuildConfigPayload) error {
	defer os.RemoveAll(cfg.RepoName)

	err := util.CloneRepo(cfg.RepoName, cfg.RepoCloneUrl, cfg.RepoUsername, cfg.RepoToken)

	if err != nil {
		return err
	}

	helper := container.NewContainerHelper("unix:///var/run/docker.sock")

	path, err := os.Getwd()

	if err != nil {
		return err
	}

	buf, err := util.TarFiles(fmt.Sprintf("%s/%s/", path, cfg.RepoName))

	if err != nil {
		return err
	}

	tag := cfg.ImageTagPrefix + cfg.RepoName
	err = helper.BuildImage(buf, &types.ImageBuildOptions{Tags: []string{tag}})

	if err != nil {
		return err
	}

	// TODO: figure out the right way of using the SDK API instead of the CMD workaround
	cmd := exec.Command("docker", "push", tag)
	log.Printf("Pushing image %s", tag)
	err = cmd.Run()
	log.Printf("Pushing finished with error: %v", err)

	return err
}

func (b *Builder) UpdateTask(input *model.UpdateBuildTaskInput) (primitive.ObjectID, error) {
	return b.repo.UpdateBuildTask(context.TODO(), input)
}

func (b *Builder) UpdateTaskStatus(input *model.UpdateBuildTaskStatusInput) error {
	return b.repo.UpdateBuildTaskStatus(context.TODO(), input)
}
