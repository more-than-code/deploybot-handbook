package task

import (
	"context"
	"net/http"
	"text/template"

	"github.com/more-than-code/deploybot/model"
)

type TaskCollection struct {
	Title1      string
	BuildTasks  []*model.BuildTask
	Title2      string
	DeployTasks []*model.DeployTask
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
	dTasks, _ := d.deployer.repo.GetDeployTasks(context.TODO(), model.DeployTasksInput{})

	coll := TaskCollection{Title1: "Build Tasks", Title2: "Deploy Tasks", BuildTasks: bTasks, DeployTasks: dTasks}

	tmpl := template.Must(template.ParseFiles("asset/tasks.html"))

	tmpl.Execute(w, coll)
}
