package repository

import (
	"context"
	"testing"
)

func TestGetTasks(t *testing.T) {
	r, _ := NewRepository()

	tasks, err := r.GetPipelines(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tasks)
}
