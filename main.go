package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/api"
	"github.com/more-than-code/deploybot/task"
)

type Config struct {
	JobRole string `envconfig:"JOB_ROLE"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	g := gin.Default()

	if cfg.JobRole == "Runner" {
		t := task.NewScheduler()
		g.POST("/ghWebhook", t.GhWebhookHandler())
		g.POST("/pkStreamWebhook", t.StreamWebhookHandler())
	} else if cfg.JobRole == "Coordinator" {
		api := api.NewApi()
		g.GET("/", api.DashboardHandler())
		g.GET("/api/pipelines", api.GetPipelines())
		g.GET("/api/pipeline/:name", api.GetPipeline())
		g.POST("/api/pipeline", api.PostPipeline())

		g.GET("/api/pipelineTask/:pid/:tid", api.GetPipelineTask())
		g.POST("/api/pipelineTask", api.PostPipelineTask())
		g.PATCH("/api/pipelineTask", api.PatchPipelineTask())
		g.PUT("/api/pipelineTaskStatus", api.PutPipelineTaskStatus())
	} else {
		t := task.NewScheduler()
		g.POST("/ghWebhook", t.GhWebhookHandler())
		g.POST("/pkStreamWebhook", t.StreamWebhookHandler())

		api := api.NewApi()
		g.GET("/", api.DashboardHandler())
		g.GET("/api/pipelines", api.GetPipelines())
		g.GET("/api/pipeline/:name", api.GetPipeline())
		g.POST("/api/pipeline", api.PostPipeline())

		g.GET("/api/pipelineTask/:pid/:tid", api.GetPipelineTask())
		g.POST("/api/pipelineTask", api.PostPipelineTask())
		g.PATCH("/api/pipelineTask", api.PatchPipelineTask())
		g.PUT("/api/pipelineTaskStatus", api.PutPipelineTaskStatus())
	}

	g.Run(":8080")
}
