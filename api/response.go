package api

import (
	"github.com/more-than-code/deploybot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CodeClientError = 1000
	CodeServerError = 2000
)

type PostPipelineResponsePayload struct {
	Id primitive.ObjectID
}
type PostPipelineResponse struct {
	Code    int
	Msg     string
	Payload PostPipelineResponsePayload
}

type GetPipelinesResponsePayload struct {
	Pipelines []*model.Pipeline
}

type GetPipelinesResponse struct {
	Code    int
	Msg     string
	Payload GetPipelinesResponsePayload
}

type GetPipelineResponsePayload struct {
	Pipeline *model.Pipeline
}

type GetPipelineResponse struct {
	Code    int
	Msg     string
	Payload GetPipelineResponsePayload
}

type PostPipelineTaskResponsePayload struct {
	Id primitive.ObjectID
}
type PostPipelineTaskResponse struct {
	Code    int
	Msg     string
	Payload PostPipelineTaskResponsePayload
}

type GetPipelineTaskResponsePayload struct {
	Task *model.Task
}
type GetPipelineTaskResponse struct {
	Code    int
	Msg     string
	Payload GetPipelineTaskResponsePayload
}

type PatchPipelineTaskResponse struct {
	Code int
	Msg  string
}

type PutPipelineTaskStatusResponse struct {
	Code int
	Msg  string
}
