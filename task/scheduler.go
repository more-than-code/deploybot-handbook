package task

import (
	"bytes"
	"container/list"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/model"
)

type Config struct {
	RepoUsername      string `envconfig:"REPO_USERNAME"`
	RepoToken         string `envconfig:"REPO_TOKEN"`
	BuildWebhook      string `envconfig:"BUILD_WEBHOOK"`
	PostBuildWebhook  string `envconfig:"POST_BUILD_WEBHOOK"`
	DeployWebhook     string `envconfig:"DEPLOY_WEBHOOK"`
	PostDeployWebhook string `envconfig:"POST_DEPLOY_WEBHOOK"`
}

var gBuildTaskTicker *time.Ticker
var gDeployTaskTicker *time.Ticker

var gEventQueue = list.New()

type Scheduler struct {
	builder  *Builder
	deployer *Deployer
	cfg      *Config
}

func NewScheduler() *Scheduler {
	var cfg Config
	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Println(err)
		return nil
	}

	return &Scheduler{builder: NewBuilder(), deployer: NewDeployer(), cfg: &cfg}
}

func (s *Scheduler) PushEvent(e model.Event) {
	gEventQueue.PushBack(e)
}

func (s *Scheduler) PullEvent() model.Event {
	e := gEventQueue.Front()

	gEventQueue.Remove(e)

	return e.Value.(model.Event)
}

func (s *Scheduler) ProcessEvents() {
	e := s.PullEvent()

	val, ok := e.Value.(string)
	if e.Key == model.EventDeploy {
		if ok {
			switch val {
			case "geoy-webapp":
			}
		}
	}
}

func (s *Scheduler) ProcessBuildTasks() {
	tasks, err := s.builder.repo.GetBuildTasks(context.TODO(), &model.BuildTasksInput{StatusFilter: &model.BuildTaskStatusFilter{Option: model.TaskPending}})

	if err != nil {
		log.Println(err)
		return
	}

	for _, t := range tasks {
		go func(t2 *model.BuildTask) {
			err := s.builder.Start(t2.Config.SourceConfig)

			if err != nil {
				log.Println(err)
				return
			}

			t2.Status = model.TaskDone
			data, _ := json.Marshal(t2)
			_, _ = http.Post(t2.Config.Webhook, "application/json", bytes.NewReader(data))

		}(t)
	}
}

func (s *Scheduler) ProcessDeployTasks() {
	tasks, err := s.deployer.repo.GetDeployTasks(context.TODO(), &model.DeployTasksInput{StatusFilter: &model.DeployTaskStatusFilter{Option: model.TaskPending}})

	if err != nil {
		log.Println(err)
		return
	}

	for _, t := range tasks {
		go func(t2 *model.DeployTask) {
			err := s.deployer.Start(t2.Config)

			if err != nil {
				log.Println(err)
				return
			}

			t2.Status = model.TaskDone
			data, _ := json.Marshal(t2)
			_, _ = http.Post(t2.Config.Webhook, "application/json", bytes.NewReader(data))

		}(t)
	}
}

func (s *Scheduler) StartBuild() {
	if gBuildTaskTicker != nil {
		return
	}

	gTicker := time.NewTicker(time.Minute)
	go func() {
		for range gTicker.C {
			s.ProcessBuildTasks()
		}
	}()
}

func (s *Scheduler) StartDeploy() {
	if gDeployTaskTicker != nil {
		return
	}

	gTicker := time.NewTicker(time.Minute)
	go func() {
		for range gTicker.C {
			s.ProcessDeployTasks()
		}
	}()
}

func (s *Scheduler) BuildHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var data model.BuildTask
	json.Unmarshal(body, &data)

	log.Printf("%+v", data)

	go func() {
		s.builder.UpdateTaskStatus(&model.UpdateBuildTaskStatusInput{BuildTaskId: data.Id, Status: model.TaskInProgress})
		err := s.builder.Start(data.Config.SourceConfig)

		if err != nil {
			s.builder.UpdateTaskStatus(&model.UpdateBuildTaskStatusInput{BuildTaskId: data.Id, Status: model.TaskFailed})
			log.Println(err)
			return
		}

		s.builder.UpdateTaskStatus(&model.UpdateBuildTaskStatusInput{BuildTaskId: data.Id, Status: model.TaskDone})
		data.Status = model.TaskDone

		bs, _ := json.Marshal(data)
		_, _ = http.Post(data.Config.Webhook, "application/json", bytes.NewReader(bs))
	}()
}

func (s *Scheduler) PostBuildHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var data model.BuildTask
	json.Unmarshal(body, &data)

	log.Printf("%+v", data)

	var cfg model.DeployConfig
	switch data.Config.SourceConfig.RepoName {
	case "geoy-webapp":
		cfg = model.DeployConfig{
			Webhook:     s.cfg.PostDeployWebhook,
			PostInstall: "docker restart swag",
			ContainerConfig: &model.ContainerConfig{
				ImageName:   "binartist/geoy-webapp",
				ImageTag:    ":latest",
				ServiceName: "geoy_webapp",
				MountTarget: "/var/www",
				AutoRemove:  true,
			},
		}
	}

	deployTaskId, _ := s.deployer.UpdateTask(&model.UpdateDeployTaskInput{Config: cfg, BuildTaskId: data.Id})

	task, _ := json.Marshal(&model.DeployTask{Id: deployTaskId, BuildTaskId: data.Id, Config: cfg})
	_, err := http.Post(s.cfg.DeployWebhook, "application/json", bytes.NewReader(task))

	if err == nil {
		s.deployer.UpdateTaskStatus(&model.UpdateDeployTaskStatusInput{DeployTaskId: deployTaskId, Status: model.TaskInProgress})
	}
}

func (s *Scheduler) DeployHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var data model.DeployTask
	json.Unmarshal(body, &data)

	log.Printf("%+v", data)

	go func() {
		err := s.deployer.Start(data.Config)

		if err != nil {
			log.Println(err)
			data.Status = model.TaskFailed
		} else {
			data.Status = model.TaskDone
		}

		bs, _ := json.Marshal(data)
		_, _ = http.Post(data.Config.Webhook, "application/json", bytes.NewReader(bs))
	}()

}

func (s *Scheduler) PostDeployHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var data model.DeployTask
	json.Unmarshal(body, &data)

	log.Printf("%+v", data)

	s.deployer.UpdateTaskStatus(&model.UpdateDeployTaskStatusInput{DeployTaskId: data.Id, Status: data.Status})
}

func (s *Scheduler) GhWebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var data model.GitHubHookshot
	json.Unmarshal(body, &data)

	log.Printf("%+v", data)

	input := &model.UpdateBuildTaskInput{
		Config: model.BuildConfig{
			Webhook:      s.cfg.PostBuildWebhook,
			SourceConfig: model.SourceConfig{RepoCloneUrl: data.Repository.CloneUrl, RepoName: data.Repository.Name, RepoUsername: s.cfg.RepoUsername, RepoToken: s.cfg.RepoToken, ImageTagPrefix: "binartist/"}}}

	buildTaskId, _ := s.builder.UpdateTask(input)

	bs, _ := json.Marshal(model.BuildTask{Id: buildTaskId, Config: input.Config})
	_, _ = http.Post(s.cfg.BuildWebhook, "application/json", bytes.NewReader(bs))
}
