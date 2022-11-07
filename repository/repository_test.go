package repository

import (
	"context"
	"testing"

	"github.com/more-than-code/deploybot/model"
)

func TestGetTasks(t *testing.T) {
	r, _ := NewRepository()

	tasks, err := r.GetPipelines(context.TODO(), model.GetPipelinesInput{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tasks)
}
