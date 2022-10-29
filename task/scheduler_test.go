package task

import (
	"context"
	"testing"

	"github.com/more-than-code/deploybot/model"
)

func TestHandleEvent(t *testing.T) {
	s := NewScheduler()
	err := s.HandleBuildEvent(context.TODO(), &model.Event{Name: "build", Data: "geoy-webapp"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDispatchDeployTask(t *testing.T) {
	s := NewScheduler()
	s.DispatchDeployTask(context.TODO())
}
