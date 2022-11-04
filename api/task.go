package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/more-than-code/deploybot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.CreateTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostTaskResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		id, err := a.repo.CreateTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostTaskResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostTaskResponse{Payload: PostTaskResponsePayload{id}})
	}
}

func (a *Api) GetTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pidStr := ctx.Param("pid")
		tidStr := ctx.Param("tid")

		pid, _ := primitive.ObjectIDFromHex(pidStr)
		tid, _ := primitive.ObjectIDFromHex(tidStr)

		input := model.GetTaskInput{PipelineId: pid, Id: tid}
		task, err := a.repo.GetTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetTaskResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetTaskResponse{Payload: GetTaskResponsePayload{task}})
	}
}

func (a *Api) PatchTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdateTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchTaskResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchTaskResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchTaskResponse{})
	}
}

func (a *Api) PutTaskStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdateTaskStatusInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutTaskStatusResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateTaskStatus(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutTaskStatusResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PutTaskStatusResponse{})
	}
}