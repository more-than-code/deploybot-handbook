package task

import (
	"github.com/docker/docker/api/types"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/util"
)

type RunnerConfig struct {
}

type Runner struct {
	cfg RunnerConfig
}

func NewRunner() *Runner {
	var cfg RunnerConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	return &Runner{}
}

func (r *Runner) DoTask(t model.Task, args []string) error {
	helper := util.NewContainerHelper("unix:///var/run/docker.sock")
	if c, ok := t.Config.(model.BuildConfig); ok {

		util.CloneRepo(c.RepoName, c.RepoUrl)
		r, err := util.TarFiles("/var/opt/projects/" + c.RepoName)

		if err != nil {
			return err
		}
		helper.BuildImage(r, &types.ImageBuildOptions{})
	}

	return nil
}
