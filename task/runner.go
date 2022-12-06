package task

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/util"
)

type RunnerConfig struct {
	ProjectsPath string `envconfig:"PROJECTS_PATH"`
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

func (r *Runner) DoTask(t model.Task, arguments []string) error {
	helper := util.NewContainerHelper("unix:///var/run/docker.sock")
	if t.Type == model.TaskBuild {

		bs, err := json.Marshal(t.Config)

		if err != nil {
			return err
		}

		var c model.BuildConfig
		err = json.Unmarshal(bs, &c)

		if err != nil {
			return err
		}

		path := r.cfg.ProjectsPath + "/" + c.RepoName
		util.CloneRepo(path, c.RepoUrl)
		r, err := util.TarFiles(path)

		if err != nil {
			return err
		}

		err = helper.BuildImage(r, &types.ImageBuildOptions{Tags: []string{c.ImageTag}})

		if err != nil {
			return nil
		}

		helper.PushImage(c.ImageName + "/" + c.ImageTag)
	} else if t.Type == model.EventDeploy {
		bs, err := json.Marshal(t.Config)

		if err != nil {
			return err
		}

		var c model.DeployConfig
		err = json.Unmarshal(bs, &c)

		if err != nil {
			return err
		}

		helper.StartContainer(&c)
	}

	return nil
}
