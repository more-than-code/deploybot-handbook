package task

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/repository"
)

type Scheduler struct {
	repo *repository.Repository
}

func NewScheduler() *Scheduler {
	r, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}

	return &Scheduler{r}
}

func (s *Scheduler) HandleBuildEvent(ctx context.Context, e *model.Event) error {
	str, ok := e.Data.(string)

	if ok {
		switch str {
		case "geoy-webapp":
			task := model.UpdateDeployTaskInput{
				Config: &model.DeployConfig{
					Webhook: "https://geoy.appsive.com/deploy",
					Payload: &model.DeployConfigPayload{
						ImageName:   "binartist/geoy-webapp",
						ImageTag:    ":latest",
						ServiceName: "geoy_webapp",
						MountTarget: "/var/www",
						AutoRemove:  true,
					},
				},
			}
			s.repo.UpdateDeployTask(ctx, &task)

			s.DispatchDeployTask(ctx)
		}
	}

	return nil
}

func (s *Scheduler) DispatchDeployTask(ctx context.Context) error {
	tasks, err := s.repo.GetDeployTasks(ctx, &model.DeployTasksInput{StatusFilter: &model.DeployStatusFilter{Option: "pending"}})

	if err != nil {
		return err
	}

	for _, t := range tasks {
		data, err := json.Marshal(t.Config.Payload)

		if err != nil {
			return nil
		}

		_, err = http.Post(t.Config.Webhook, "application/json", bytes.NewReader(data))

		if err != nil {
			return nil
		}

	}

	return nil
}
