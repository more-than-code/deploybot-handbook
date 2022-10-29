package main

import (
	"context"
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

func main() {
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
		var cfg Config
		err := envconfig.Process("", &cfg)

		if err != nil {
			log.Println(err)
			return
		}

		t := task.NewBuilder(&task.BuildConfig{RepoCloneUrl: data.Repository.CloneUrl, RepoName: data.Repository.Name, RepoUsername: cfg.RepoUsername, RepoToken: cfg.RepoToken, ImageTagPrefix: "binartist/"})

		err = t.Start()

		if err != nil {
			log.Println(err)
			return
		}

		s := task.NewScheduler()
		s.HandleBuildEvent(context.TODO(), &model.Event{Name: "build", Data: data.Repository.Name})
	}()
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var data model.DeployConfigPayload
	json.Unmarshal(body, &data)

	log.Printf("%+v", data)

	go func() {
		t := task.NewDeployer(&data)

		err := t.Start()
		log.Println(err)
	}()
}
