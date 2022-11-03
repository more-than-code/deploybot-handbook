package api

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/repository"
)

type Api struct {
	repo *repository.Repository
}

type TaskCollection struct {
	Pipelines []*model.Pipeline
}

func NewApi() *Api {
	r, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}
	return &Api{repo: r}
}

func (a *Api) DashboardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pls, _ := a.repo.GetPipelines(ctx)

		tmpl := template.Must(template.ParseFiles("asset/tasks.html"))

		tmpl.Execute(ctx.Writer, pls)
	}
}

func (a *Api) PostPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.CreatePipelineInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		id, err := a.repo.CreatePipeline(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Payload: PostPipelineResponsePayload{id}})
	}

}

func (a *Api) GetPipelines() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pls, err := a.repo.GetPipelines(ctx)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetPipelinesResponse{Code: CodeClientError, Msg: err.Error()})
		}

		ctx.JSON(http.StatusBadRequest, GetPipelinesResponse{Payload: GetPipelinesResponsePayload{pls}})
	}
}

func (a *Api) GetPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.GetPipelineInput
		pl, err := a.repo.GetPipeline(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetPipelineResponse{Code: CodeClientError, Msg: err.Error()})
		}

		ctx.JSON(http.StatusBadRequest, GetPipelineResponse{Payload: GetPipelineResponsePayload{pl}})
	}
}

func (a *Api) PostPipelineTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.CreatePipelineTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineTaskResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		id, err := a.repo.CreatePipelineTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineTaskResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, PostPipelineTaskResponse{Payload: PostPipelineTaskResponsePayload{id}})
	}
}

func (a *Api) GetPipelineTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.GetPipelineTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetPipelineTaskResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		task, err := a.repo.GetPipelineTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetPipelineTaskResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, GetPipelineTaskResponse{Payload: GetPipelineTaskResponsePayload{task}})
	}
}

func (a *Api) PutPipelineTaskStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdatePipelineTaskStatusInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutPipelineTaskStatusResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdatePipelineTaskStatus(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutPipelineTaskStatusResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, PutPipelineTaskStatusResponse{})
	}
}
