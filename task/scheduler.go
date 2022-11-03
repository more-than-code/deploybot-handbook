package task

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (s *Scheduler) ProcessPostTask(pipelineId, taskId, nextTaskId primitive.ObjectID, webhook string) {
	body, _ := json.Marshal(model.UpdatePipelineTaskStatusInput{
		PipelineId: pipelineId,
		TaskId:     taskId,
		Payload:    model.UpdatePipelineTaskStatusInputPayload{Status: model.TaskDone}})
	http.Post(s.cfg.ApiBaseUrl+"/pipelineTaskStatus", "application/json", bytes.NewReader(body))

	body, _ = json.Marshal(model.StreamWebhook{Payload: model.StreamWebhookPayload{PipelineId: pipelineId, TaskId: nextTaskId}})
	http.Post(webhook, "application/json", bytes.NewReader(body))
}

func (s *Scheduler) StreamWebhookHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)

		var sw model.StreamWebhook
		json.Unmarshal(body, &sw)

		url := fmt.Sprintf("%s/pipelineTask/%s/%s", s.cfg.ApiBaseUrl, sw.Payload.PipelineId.Hex(), sw.Payload.TaskId.Hex())

		res, _ := http.Get(url)
		body, _ = io.ReadAll(res.Body)

		var ptRes api.GetPipelineTaskResponse
		json.Unmarshal(body, &ptRes)

		t := ptRes.Payload.Task

		if t != nil {
			go func() {
				s.runner.DoTask(*t)
				s.ProcessPostTask(sw.Payload.PipelineId, sw.Payload.TaskId, t.Config.DownstreamTaskId, t.Config.DownstreamWebhook)
			}()
		}
	}
}

func (s *Scheduler) GhWebhookHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)

		var data model.GitHubHookshot
		json.Unmarshal(body, &data)

		log.Printf("%+v", data.Commits)

		res, _ := http.Get(fmt.Sprintf("%s/pipeline/%s", s.cfg.ApiBaseUrl, data.Repository.Name))
		body, _ = io.ReadAll(res.Body)

		var plRes api.GetPipelineResponse
		json.Unmarshal(body, &plRes)

		t := plRes.Payload.Pipeline.Tasks[0]

		if t != nil {
			go func() {
				s.runner.DoTask(*t)
				s.ProcessPostTask(plRes.Payload.Pipeline.Id, t.Id, t.Config.DownstreamTaskId, t.Config.DownstreamWebhook)
			}()
		}
	}
}

func (s *Scheduler) CreatePipeline(name string) (primitive.ObjectID, error) {
	body, _ := json.Marshal(model.CreatePipelineInput{Payload: model.CreatePipelineInputPayload{Name: name}})
	res, _ := http.Post(s.cfg.ApiBaseUrl+"/pipeline", "application/json", bytes.NewReader(body))
	body, _ = io.ReadAll(res.Body)
	var plRes api.PostPipelineResponse
	err := json.Unmarshal(body, &plRes)

	return plRes.Payload.Id, err
}

func (s *Scheduler) CreateTask(pipelineId, taskId, downstreamTaskId primitive.ObjectID, script, upstreamWebhook, downstreamWebhook string) (primitive.ObjectID, error) {
	body, err := json.Marshal(model.CreatePipelineTaskInput{PipelineId: pipelineId, Payload: model.CreatePipelineTaskInputPayload{Id: taskId, Config: model.TaskConfig{DownstreamTaskId: downstreamTaskId, DownstreamWebhook: downstreamWebhook, Script: script}}})

	if err != nil {
		return primitive.NilObjectID, err
	}

	res, err := http.Post(s.cfg.ApiBaseUrl+"/pipelineTask", "application/json", bytes.NewReader(body))
	if err != nil {
		return primitive.NilObjectID, err
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return primitive.NilObjectID, err
	}

	var ptRes api.PostPipelineTaskResponse
	err = json.Unmarshal(body, &ptRes)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return ptRes.Payload.Id, err
}
