package task

import (
	"testing"

	"github.com/more-than-code/deploybot/model"
)

func TestHandleEvent(t *testing.T) {
	s := NewScheduler()
	s.PushEvent(model.Event{Key: "build", Value: "geoy-webapp"})

	e := s.PullEvent()
	if e.Key != "build" {
		t.Fail()
	}
}
