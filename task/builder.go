package task

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/more-than-code/deploybot/container"
	"github.com/more-than-code/deploybot/util"
)

type BuildConfig struct {
	RepoCloneUrl   string
	RepoName       string
	RepoUsername   string
	RepoToken      string
	ImageTagPrefix string
}

type Builder struct {
	cfg *BuildConfig
}

func NewBuilder(cfg *BuildConfig) *Builder {
	return &Builder{cfg: cfg}
}

func (t *Builder) Start() error {
	defer os.RemoveAll(t.cfg.RepoName)

	err := util.CloneRepo(t.cfg.RepoName, t.cfg.RepoCloneUrl, t.cfg.RepoUsername, t.cfg.RepoToken)

	if err != nil {
		return err
	}

	helper := container.NewContainerHelper("unix:///var/run/docker.sock")

	path, err := os.Getwd()

	if err != nil {
		return err
	}

	buf, err := util.TarFiles(fmt.Sprintf("%s/%s/", path, t.cfg.RepoName))

	if err != nil {
		return err
	}

	tag := t.cfg.ImageTagPrefix + t.cfg.RepoName
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
