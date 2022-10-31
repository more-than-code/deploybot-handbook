package main

import (
	"log"
	"net/http"

	"github.com/more-than-code/deploybot/task"
)

func main() {
	s := task.NewScheduler()

	http.HandleFunc("/gh_webhook", s.GhWebhookHandler)
	http.HandleFunc("/build", s.BuildHandler)
	http.HandleFunc("/postBuild", s.PostBuildHandler)
	http.HandleFunc("/deploy", s.DeployHandler)
	http.HandleFunc("/postDeploy", s.PostDeployHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
