package task

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/model"
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
	if t.Config.Script != "" {

	}

	return nil
}
