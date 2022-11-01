package repository

import (
	"context"
	"testing"

	"github.com/more-than-code/deploybot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetDeployTasks(t *testing.T) {
	r, _ := NewRepository()

	tasks, err := r.GetDeployTasks(context.TODO(), model.DeployTasksInput{
		StatusFilter: &model.DeployTaskStatusFilter{Option: model.TaskPending},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tasks)
}

func TestCreateDeployTask(t *testing.T) {
	r, _ := NewRepository()

	id, _ := primitive.ObjectIDFromHex("635d38a1988fd51e865a5244")
	task := model.UpdateDeployTaskInput{
		Id: id,
		Config: model.DeployConfig{
			Webhook: "https://geoy.appsive.com/deploy",
			ContainerConfig: &model.ContainerConfig{
				ImageName:   "binartist/geoy-webapp",
				ImageTag:    ":latest",
				ServiceName: "geoy_webapp",
				MountTarget: "/var/www",
				AutoRemove:  true,
			},
		},
	}
	_, err := r.UpdateDeployTask(context.TODO(), &task)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(task)
}

func TestDeleteDeployTask(t *testing.T) {
	r, _ := NewRepository()

	task := model.UpdateDeployTaskInput{}
	id, _ := r.UpdateDeployTask(context.TODO(), &task)

	err := r.DeleteDeployTasks(context.TODO(), id)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(err)
}
