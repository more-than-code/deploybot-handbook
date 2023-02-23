package task

import (
	"encoding/json"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/util"
)

type RunnerConfig struct {
	ProjectsPath string `envconfig:"PROJECTS_PATH"`
	DockerHost   string `envconfig:"DOCKER_HOST"`
}

type Runner struct {
	cfg     RunnerConfig
	cHelper *util.ContainerHelper
}

func NewRunner() *Runner {
	var cfg RunnerConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	return &Runner{cfg: cfg, cHelper: util.NewContainerHelper(cfg.DockerHost)}
}

func (r *Runner) DoTask(t model.Task, arguments []string) error {
	if t.Type == model.TaskBuild {
		var c model.BuildConfig

		bs, err := json.Marshal(t.Config)

		if err != nil {
			return err
		}

		err = json.Unmarshal(bs, &c)

		if err != nil {
			return err
		}

		path := r.cfg.ProjectsPath + "/" + c.RepoName + "/"

		os.RemoveAll(path)
		err = util.CloneRepo(path, c.RepoUrl)

		if err != nil {
			return err
		}

		files, err := util.TarFiles(path)

		if err != nil {
			return err
		}

		imageNameTag := c.ImageName + ":" + c.ImageTag

		err = r.cHelper.BuildImage(files, &types.ImageBuildOptions{Dockerfile: c.Dockerfile, Tags: []string{imageNameTag}, BuildArgs: c.Args})

		if err != nil {
			return err
		}

		r.cHelper.PushImage(c.ImageName)
	} else if t.Type == model.TaskDeploy {
		var c model.DeployConfig

		bs, err := json.Marshal(t.Config)

		if err != nil {
			return err
		}

		err = json.Unmarshal(bs, &c)

		if err != nil {
			return err
		}

		go func() {
			r.cHelper.StartContainer(&c)
		}()
	} else if t.Type == model.TaskRestart {
		var c model.RestartConfig

		bs, err := json.Marshal(t.Config)

		if err != nil {
			return err
		}

		err = json.Unmarshal(bs, &c)

		if err != nil {
			return err
		}

		r.cHelper.RestartContainer(&c)
	}

	return nil
}
