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

		m := map[string]interface{}{}

		list := t.Config.([]interface{})

		for _, e := range list {
			e2 := e.(map[string]interface{})
			m[e2["Key"].(string)] = e2["Value"]
		}

		bs, err := json.Marshal(m)

		if err != nil {
			return nil
		}

		var c model.BuildConfig
		err = json.Unmarshal(bs, &c)

		if err != nil {
			return nil
		}

		path := r.cfg.ProjectsPath + "/" + c.RepoName + "/"

		os.RemoveAll(path)
		util.CloneRepo(path, c.RepoUrl)
		files, err := util.TarFiles(path)

		if err != nil {
			return nil
		}

		err = r.cHelper.BuildImage(files, &types.ImageBuildOptions{Tags: []string{c.ImageTag}})

		if err != nil {
			return nil
		}

		r.cHelper.PushImage(c.ImageName + "/" + c.ImageTag)
	} else if t.Type == model.EventDeploy {
		m := map[string]interface{}{}

		list := t.Config.([]interface{})

		for _, e := range list {
			e2 := e.(map[string]interface{})
			m[e2["Key"].(string)] = e2["Value"]
		}

		bs, err := json.Marshal(m)

		if err != nil {
			return nil
		}

		var c model.DeployConfig
		err = json.Unmarshal(bs, &c)

		if err != nil {
			return nil
		}

		r.cHelper.StartContainer(&c)
	}

	return nil
}
