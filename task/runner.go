package task

import (
	"log"
	"os/exec"

	"github.com/more-than-code/deploybot/model"
)

type Runner struct {
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) DoTask(t model.Task) error {
	if t.Config.Script != "" {
		cmd := exec.Command("sh", "-c", t.Config.Script)
		output, err := cmd.Output()
		log.Println(string(output))
		if err != nil {
			return err
		}
	}

	return nil
}
