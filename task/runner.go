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

	return &Runner{cfg: cfg, cHelper: util.NewContainerHelper("unix:///var/run/docker.sock")}
}

func (r *Runner) DoTask(t model.Task, arguments []string) error {
	if t.Type == model.TaskBuild {
		var c model.BuildConfig

		if conf, ok := t.Config.(model.BuildConfig); ok {
			c = conf
		} else {
			m := map[string]interface{}{}

			list := t.Config.([]interface{})

			for _, e := range list {
				e2 := e.(map[string]interface{})
				m[e2["Key"].(string)] = e2["Value"]
			}

			bs, err := json.Marshal(m)

			if err != nil {
				return err
			}

			err = json.Unmarshal(bs, &c)

			if err != nil {
				return err
			}
		}

		path := r.cfg.ProjectsPath + "/" + c.RepoName + "/"

		os.RemoveAll(path)
		util.CloneRepo(path, c.RepoUrl)
		files, err := util.TarFiles(path)

		if err != nil {
			return err
		}

		imageNameTag := c.ImageName + ":" + c.ImageTag

		err = r.cHelper.BuildImage(files, &types.ImageBuildOptions{Dockerfile: c.Dockerfile, Tags: []string{imageNameTag}})

		if err != nil {
			return err
		}

		err = r.cHelper.PushImage(imageNameTag)

		if err != nil {
			return err
		}
	} else if t.Type == model.EventDeploy {
		var c model.DeployConfig

		if conf, ok := t.Config.(model.DeployConfig); ok {
			c = conf
		} else {
			m := map[string]interface{}{}

			list := t.Config.([]interface{})

			for _, e := range list {
				e2 := e.(map[string]interface{})
				m[e2["Key"].(string)] = e2["Value"]
			}

			bs, err := json.Marshal(m)

			if err != nil {
				return err
			}

			err = json.Unmarshal(bs, &c)

			if err != nil {
				return err
			}
		}

		r.cHelper.StartContainer(&c)
	}

	return nil
}
