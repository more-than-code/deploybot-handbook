package main

import (
	"log"
	"net/http"

	"github.com/more-than-code/deploybot/task"
)

func main() {
	s := task.NewScheduler()

	http.HandleFunc("/ghWebhook", s.GhWebhookHandler)
	http.HandleFunc("/pkBuild", s.BuildHandler)
	http.HandleFunc("/pkPostBuild", s.PostBuildHandler)
	http.HandleFunc("/pkDeploy", s.DeployHandler)
	http.HandleFunc("/pkPostDeploy", s.PostDeployHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
