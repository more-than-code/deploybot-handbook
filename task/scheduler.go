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

func (s *Scheduler) ProcessPreTask(pipelineId, taskId primitive.ObjectID, sourceRef *string) {
	body, _ := json.Marshal(model.UpdateTaskStatusInput{
		PipelineId: pipelineId,
		TaskId:     taskId,
		Payload:    model.UpdateTaskStatusInputPayload{Status: model.TaskInProgress}})

	req, _ := http.NewRequest("PUT", s.cfg.ApiBaseUrl+"/taskStatus", bytes.NewReader(body))
	http.DefaultClient.Do(req)

	if sourceRef != nil {
		body, _ = json.Marshal(model.UpdateTaskInput{
			PipelineId: pipelineId,
			Id:         taskId,
			Payload:    model.UpdateTaskInputPayload{Remarks: sourceRef}})

		req, _ = http.NewRequest("PATCH", s.cfg.ApiBaseUrl+"/task", bytes.NewReader(body))
		http.DefaultClient.Do(req)
	}

	body, _ = json.Marshal(model.UpdatePipelineStatusInput{
		PipelineId: pipelineId,
		Payload:    model.UpdatePipelineStatusInputPayload{Status: model.PipelineBusy}})

	req, _ = http.NewRequest("PUT", s.cfg.ApiBaseUrl+"/pipelineStatus", bytes.NewReader(body))
	http.DefaultClient.Do(req)
}

func (s *Scheduler) ProcessPostTask(pipelineId, taskId, nextTaskId primitive.ObjectID, webhook string) {
	body, _ := json.Marshal(model.UpdateTaskStatusInput{
		PipelineId: pipelineId,
		TaskId:     taskId,
		Payload:    model.UpdateTaskStatusInputPayload{Status: model.TaskDone}})

	req, _ := http.NewRequest("PUT", s.cfg.ApiBaseUrl+"/taskStatus", bytes.NewReader(body))
	http.DefaultClient.Do(req)

	if nextTaskId.IsZero() {
		body, _ = json.Marshal(model.UpdatePipelineStatusInput{
			PipelineId: pipelineId,
			Payload:    model.UpdatePipelineStatusInputPayload{Status: model.PipelineIdle}})

		req, _ = http.NewRequest("PUT", s.cfg.ApiBaseUrl+"/pipelineStatus", bytes.NewReader(body))
		http.DefaultClient.Do(req)
	}

	body, _ = json.Marshal(model.StreamWebhook{Payload: model.StreamWebhookPayload{PipelineId: pipelineId, TaskId: nextTaskId}})
	http.Post(webhook, "application/json", bytes.NewReader(body))
}

func (s *Scheduler) StreamWebhookHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)

		var sw model.StreamWebhook
		json.Unmarshal(body, &sw)

		url := fmt.Sprintf("%s/task/%s/%s", s.cfg.ApiBaseUrl, sw.Payload.PipelineId.Hex(), sw.Payload.TaskId.Hex())

		res, _ := http.Get(url)
		body, _ = io.ReadAll(res.Body)

		var ptRes api.GetTaskResponse
		json.Unmarshal(body, &ptRes)

		t := ptRes.Payload.Task

		if t == nil {
			ctx.JSON(api.ExHttpStatusBusinessLogicError, api.WebhookResponse{Code: api.CodeTaskNotFound, Msg: api.MsgTaskNotFound})
			return
		}

		if t.Status == model.TaskInProgress {
			ctx.JSON(api.ExHttpStatusBusinessLogicError, api.WebhookResponse{Code: api.CodeTaskBusy, Msg: api.MsgTaskBusy})
			return
		}

		go func() {
			s.ProcessPreTask(sw.Payload.PipelineId, sw.Payload.TaskId, sw.Payload.Remarks)
			s.runner.DoTask(*t)
			s.ProcessPostTask(sw.Payload.PipelineId, sw.Payload.TaskId, t.Config.DownstreamTaskId, t.Config.DownstreamWebhook)
		}()

		ctx.JSON(http.StatusOK, api.WebhookResponse{})
	}
}

func (s *Scheduler) GhWebhookHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)

		var data model.GitHubHookshot
		json.Unmarshal(body, &data)

		res, _ := http.Get(fmt.Sprintf("%s/pipelines?repoWatched=%s&autoRun=true", s.cfg.ApiBaseUrl, data.Repository.Name))
		body, _ = io.ReadAll(res.Body)

		var plRes api.GetPipelinesResponse
		json.Unmarshal(body, &plRes)

		for _, pl := range plRes.Payload.Pipelines {
			if pl.Status == model.PipelineBusy {
				ctx.JSON(api.ExHttpStatusBusinessLogicError, api.WebhookResponse{Code: api.CodePipelineBusy, Msg: api.MsgPipelineBusy})
				continue
			}

			t := pl.Tasks[0]

			cbs, _ := json.Marshal(data.Commits)
			cbsStr := string(cbs)

			log.Printf("%s", cbsStr)

			if t == nil {
				continue
			}

			pl2 := pl
			go func() {
				s.ProcessPreTask(pl2.Id, t.Id, &cbsStr)
				s.runner.DoTask(*t)
				s.ProcessPostTask(pl2.Id, t.Id, t.Config.DownstreamTaskId, t.Config.DownstreamWebhook)
			}()

		}

		ctx.JSON(http.StatusOK, api.WebhookResponse{})
	}
}

func (s *Scheduler) HealthCheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

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
	body, err := json.Marshal(model.CreateTaskInput{PipelineId: pipelineId, Payload: model.CreateTaskInputPayload{Id: taskId, Config: model.TaskConfig{DownstreamTaskId: downstreamTaskId, DownstreamWebhook: downstreamWebhook, Script: script}}})

	if err != nil {
		return primitive.NilObjectID, err
	}

	res, err := http.Post(s.cfg.ApiBaseUrl+"/task", "application/json", bytes.NewReader(body))
	if err != nil {
		return primitive.NilObjectID, err
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return primitive.NilObjectID, err
	}

	var ptRes api.PostTaskResponse
	err = json.Unmarshal(body, &ptRes)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return ptRes.Payload.Id, err
}
