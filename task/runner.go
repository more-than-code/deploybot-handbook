package task

import (
	"io"
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
		stdout, _ := cmd.StdoutPipe()

		output, _ := io.ReadAll(stdout)
		log.Println(string(output))

		err := cmd.Start()
		if err != nil {
			return err
		}
	}

	return nil
}
