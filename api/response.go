package api

import (
	"github.com/more-than-code/deploybot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ExHttpStatusAuthenticationFailure = 460
	ExHttpStatusBusinessLogicError    = 461
)

const (
	CodeClientError = 1000
	CodeServerError = 2000
)

const (
	MsgPipelineBusy = "Pipleline busy"
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

type PatchPipelineResponse struct {
	Code int
	Msg  string
}

type PutPipelineStatusResponse struct {
	Code int
	Msg  string
}

type PostTaskResponsePayload struct {
	Id primitive.ObjectID
}
type PostTaskResponse struct {
	Code    int
	Msg     string
	Payload PostTaskResponsePayload
}

type GetTaskResponsePayload struct {
	Task *model.Task
}
type GetTaskResponse struct {
	Code    int
	Msg     string
	Payload GetTaskResponsePayload
}

type PatchTaskResponse struct {
	Code int
	Msg  string
}

type PutTaskStatusResponse struct {
	Code int
	Msg  string
}

type WebhookResponse struct {
	Code int
	Msg  string
}
