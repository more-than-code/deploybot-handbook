package task

import (
	"testing"

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
	`

	_, err = s.CreateTask(plId, t2Id, primitive.NilObjectID, script, "", "")

	if err != nil {
		t.Fatal(err)
	}
}
