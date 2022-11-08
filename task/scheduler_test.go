package task

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/more-than-code/deploybot/api"
	"github.com/more-than-code/deploybot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandleEvent(t *testing.T) {
	s := NewScheduler()
	s.PushEvent(model.Event{Key: "build", Value: "geoy-webapp"})

	e := s.PullEvent()
	if e.Key != "build" {
		t.Fail()
	}
}

func TestCreateTasks(t *testing.T) {
	s := NewScheduler()

	plId, _ := s.CreatePipeline("geoy-webapp")
	tId := primitive.NewObjectID()
	t2Id := primitive.NewObjectID()

	script := `
		rm -rf geoy-webapp
		git clone https://github.com/joe-and-his-friends/geoy-webapp.git
		docker build geoy-webapp -t binartist/geoy-webapp
		docker push binartist/geoy-webapp
	`
	_, err := s.CreateTask(plId, tId, t2Id, script, "", "https://geoy.appsive.com/pkStreamWebook")

	if err != nil {
		t.Fatal(err)
	}

	script = `
		sudo mkdir -p /var/appdata/geoy_webapp
		docker pull binartist/geoy-webapp
		docker run --rm --name geoy_webapp -v /var/appdata/geoy_webapp:/var/www binartist/geoy-webapp
		docker restart swag
	`

	_, err = s.CreateTask(plId, t2Id, primitive.NilObjectID, script, "", "")

	if err != nil {
		t.Fatal(err)
	}

	plId, _ = s.CreatePipeline("geoy-services")
	tId = primitive.NewObjectID()
	t2Id = primitive.NewObjectID()

	script = `
		rm -rf geoy-services
		git clone https://github.com/joe-and-his-friends/geoy-services.git
		docker compose -f geoy-services/docker-compose.yaml build graph
		docker compose -f geoy-services/docker-compose.yaml push graph
	`
	_, err = s.CreateTask(plId, tId, t2Id, script, "", "https://geoy.appsive.com/pkStreamWebook")

	if err != nil {
		t.Fatal(err)
	}

	script = `
		docker compose -f geoy-services/docker-compose.yaml pull graph
		docker compose -f geoy-services/docker-compose.yaml up -d graph
		docker restart swag
	`

	_, err = s.CreateTask(plId, t2Id, primitive.NilObjectID, script, "", "")

	if err != nil {
		t.Fatal(err)
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
	body, err := json.Marshal(model.CreateTaskInput{PipelineId: pipelineId, Payload: model.CreateTaskInputPayload{Id: taskId, Config: model.TaskRunConfig{Script: script}}})

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
