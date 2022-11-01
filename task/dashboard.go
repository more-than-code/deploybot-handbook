package task

import (
	"context"
	"net/http"
	"text/template"

	"github.com/more-than-code/deploybot/model"
)

type TaskInfo struct {
	Title string
	Log   interface{}
}

type Dashboard struct {
	builder  *Builder
	deployer *Deployer
}

func NewDashboard() *Dashboard {
	return &Dashboard{builder: NewBuilder(), deployer: NewDeployer()}
}

func (d *Dashboard) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	bTasks, _ := d.builder.repo.GetBuildTasks(context.TODO(), model.BuildTasksInput{})
	// dTasks, _ := d.deployer.repo.GetDeployTasks(context.TODO(), model.DeployTasksInput{})
	data := []TaskInfo{}

	for _, t := range bTasks {
		data = append(data, TaskInfo{Title: t.Id.Hex(), Log: t})
	}

	tmpl := template.Must(template.ParseFiles("asset/tasks.html"))

	tmpl.Execute(w, data)
}
