package task

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/api"
	"github.com/more-than-code/deploybot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var gTaskTicker *time.Ticker
var gEventQueue = list.New()

type Config struct {
	ApiBaseUrl string `envconfig:"API_BASE_URL"`
	PkUsername string `envconfig:"PK_USERNAME"`
	PkPassword string `envconfig:"PK_PASSWORD"`
}

type Scheduler struct {
	runner *Runner
	cfg    Config
}

func NewScheduler() *Scheduler {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	return &Scheduler{runner: NewRunner(), cfg: cfg}
}

func (s *Scheduler) PushEvent(e model.Event) {
	gEventQueue.PushBack(e)
}

func (s *Scheduler) PullEvent() model.Event {
	e := gEventQueue.Front()

	gEventQueue.Remove(e)

	return e.Value.(model.Event)
}

func (s *Scheduler) ProcessPreTask(pipelineId, taskId primitive.ObjectID) {
	body, _ := json.Marshal(model.UpdateTaskStatusInput{
		PipelineId: pipelineId,
		TaskId:     taskId,
		Payload:    model.UpdateTaskStatusInputPayload{Status: model.TaskInProgress}})

	req, _ := http.NewRequest("PUT", s.cfg.ApiBaseUrl+"/taskStatus", bytes.NewReader(body))
	http.DefaultClient.Do(req)
}

func (s *Scheduler) ProcessPostTask(pipelineId, taskId primitive.ObjectID, status string) {
	body, _ := json.Marshal(model.UpdateTaskStatusInput{
		PipelineId: pipelineId,
		TaskId:     taskId,
		Payload:    model.UpdateTaskStatusInputPayload{Status: status}})

	req, _ := http.NewRequest("PUT", s.cfg.ApiBaseUrl+"/taskStatus", bytes.NewReader(body))
	http.DefaultClient.Do(req)
}

func (s *Scheduler) StreamWebhookHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)

		var sw model.StreamWebhook
		json.Unmarshal(body, &sw)

		log.Println(sw.Payload)

		go func() {
			s.ProcessPreTask(sw.Payload.PipelineId, sw.Payload.Task.Id)
			err := s.runner.DoTask(sw.Payload.Task, sw.Payload.Arguments)
			if err != nil {
				log.Println(err)
				s.ProcessPostTask(sw.Payload.PipelineId, sw.Payload.Task.Id, model.TaskFailed)
			} else {
				s.ProcessPostTask(sw.Payload.PipelineId, sw.Payload.Task.Id, model.TaskDone)
			}
		}()

		ctx.JSON(http.StatusOK, api.WebhookResponse{})
	}
}

func (s *Scheduler) GhWebhookHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)

		var data model.GitHubHookshot
		json.Unmarshal(body, &data)

		comps := strings.Split(data.Ref, "/")
		branch := ""
		imageTag := "latest"
		if len(comps) == 3 {
			if comps[1] == "tags" {
				imageTag = comps[2]
			} else {
				branch = comps[2]
			}
		}

		res, _ := http.Get(fmt.Sprintf("%s/pipelines?repoWatched=%s&branchWatched=%s&autoRun=true", s.cfg.ApiBaseUrl, data.Repository.Name, branch))
		body, _ = io.ReadAll(res.Body)

		var plRes api.GetPipelinesResponse
		json.Unmarshal(body, &plRes)

		for _, pl := range plRes.Payload.Pipelines {
			if pl.Status == model.PipelineBusy {
				continue
			}

			if len(pl.Tasks) == 0 {
				continue
			}

			t := pl.Tasks[0]

			if t.Status == model.TaskInProgress {
				continue
			}

			// update task
			cbs, _ := json.Marshal(data.Commits)
			cbsStr := string(cbs)

			log.Printf("%s", cbsStr)

			body, _ = json.Marshal(model.UpdateTaskInput{
				PipelineId: pl.Id,
				Id:         t.Id,
				Payload:    model.UpdateTaskInputPayload{Remarks: &cbsStr}})
			req, _ := http.NewRequest("PATCH", s.cfg.ApiBaseUrl+"/task", bytes.NewReader(body))
			http.DefaultClient.Do(req)

			// update pipeline
			args := []string{fmt.Sprintf("IMAGE_TAG=%s", imageTag)}

			body, _ = json.Marshal(model.UpdatePipelineInput{
				Id:      pl.Id,
				Payload: model.UpdatePipelineInputPayload{Arguments: args}})
			req, _ = http.NewRequest("PATCH", s.cfg.ApiBaseUrl+"/pipeline", bytes.NewReader(body))
			http.DefaultClient.Do(req)

			// call stream webhook
			body, _ = json.Marshal(model.StreamWebhook{Payload: model.StreamWebhookPayload{PipelineId: pl.Id, Task: t, Arguments: args}})
			req, _ = http.NewRequest("POST", t.StreamWebhook, bytes.NewReader(body))
			req.SetBasicAuth(s.cfg.PkUsername, s.cfg.PkPassword)
			res, _ := http.DefaultClient.Do(req)

			if res != nil {
				log.Println(res.Status)
			}
		}

		ctx.JSON(http.StatusOK, api.WebhookResponse{})
	}
}

func (s *Scheduler) HealthCheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
