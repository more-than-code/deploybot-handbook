package task

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/more-than-code/deploybot/model"
)

type Runner struct {
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) DoTask(t model.Task, args []string) error {
	const pipe = "/var/opt/mypipe"
	if t.Config.Script != "" {
		envVarStr := strings.Join(args, " ")

		data, _ := json.Marshal(t)
		callback := fmt.Sprintf("printf '%s\n' > ./mypipe", data)
		err := os.WriteFile(pipe, []byte(fmt.Sprintf("%s; %s; %s", envVarStr, t.Config.Script, callback)), 0644)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(pipe, os.O_CREATE, os.ModeNamedPipe)
		if err != nil {
			return err
		}

		reader := bufio.NewReader(file)

		for {
			res, err := reader.ReadBytes('\n')
			if err == nil {
				log.Println(string(res))
				break
			}
		}
	}

	return nil
}
