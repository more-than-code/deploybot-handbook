package task

import (
	"errors"
	"log"
	"os/exec"

	"github.com/more-than-code/deploybot/model"
)

type Runner struct {
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) DoTask(t *model.Task) error {
	if t == nil {
		return errors.New("nil pointer")
	}

	if t.Config.Script != "" {
		cmd := exec.Command("sh", t.Config.Script)
		output, err := cmd.Output()
		log.Println(string(output))
		if err != nil {
			return err
		}
	}

	return nil
}
