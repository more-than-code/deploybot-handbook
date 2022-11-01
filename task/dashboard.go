package task

import (
	"context"
	"net/http"
	"text/template"

	"github.com/more-than-code/deploybot/model"
)

type TaskInfo struct {
	Id  string
	Log interface{}
}

type TaskCollection struct {
	Title1      string
	BuildTasks  []TaskInfo
	Title2      string
	DeployTasks []TaskInfo
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
	bTaskInfos := []TaskInfo{}
	dTaskInfos := []TaskInfo{}

	for _, t := range bTasks {
		bTaskInfos = append(bTaskInfos, TaskInfo{Id: t.Id.Hex(), Log: t})
	}

	for _, t := range dTasks {
		dTaskInfos = append(dTaskInfos, TaskInfo{Id: t.Id.Hex(), Log: t})
	}

	tmpl := template.Must(template.ParseFiles("asset/tasks.html"))

	tmpl.Execute(w, TaskCollection{Title1: "Build Tasks", Title2: "Deploy Tasks", BuildTasks: bTaskInfos, DeployTasks: dTaskInfos})
}
