package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/api"
	"github.com/more-than-code/deploybot/task"
)

type Config struct {
	JobRole    int `envconfig:"JOB_ROLE"`
	ServerPort int `envconfig:"SERVER_PORT"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	g := gin.Default()

	if cfg.JobRole == 0 || cfg.JobRole == 2 {
		t := task.NewScheduler()
		g.POST("/ghWebhook", t.GhWebhookHandler())
		g.POST("/pkStreamWebhook", t.StreamWebhookHandler())
		g.GET("/pkHealthCheck", t.HealthCheckHandler())
	}

	if cfg.JobRole == 1 || cfg.JobRole == 2 {
		api := api.NewApi()

		g.GET("/api/pipelines", api.GetPipelines())
		g.GET("/api/pipeline/:name", api.GetPipeline())
		g.POST("/api/pipeline", api.PostPipeline())
		g.PATCH("/api/pipeline", api.PatchPipeline())
		g.PUT("/api/pipelineStatus", api.PutPipelineStatus())

		g.GET("/api/task/:pid/:tid", api.GetTask())
		g.POST("/api/task", api.PostTask())
		g.PATCH("/api/task", api.PatchTask())
		g.PUT("/api/taskStatus", api.PutTaskStatus())
	}

	g.Run(fmt.Sprintf(":%d", cfg.ServerPort))
}
