package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/task"
)

type Config struct {
	RepoUsername string `envconfig:"REPO_USERNAME"`
	RepoToken    string `envconfig:"REPO_TOKEN"`
}

var username string
var token string

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	username = cfg.RepoUsername
	token = cfg.RepoToken

	http.HandleFunc("/build", buildHandler)
	http.HandleFunc("/deploy", deployHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buildHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var data model.GitHubHookshot
	json.Unmarshal(body, &data)

	log.Printf("%+v", data)

	go func() {
		t := task.NewBuildTask(&task.BuildConfig{RepoCloneUrl: data.Repository.CloneUrl, RepoName: data.Repository.Name, RepoUsername: username, RepoToken: token, ImageTagPrefix: "binartist/"})

		err := t.Start()
		log.Println(err)
	}()
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var data model.DeployTarget
	json.Unmarshal(body, &data)

	log.Printf("%+v", data)

	go func() {
		t := task.NewDeployTask(&task.DeployConfig{ImageName: data.ImageName, ImageTag: data.ImageTag, ContainerName: data.ContainerName})

		err := t.Start()
		log.Println(err)
	}()
}
