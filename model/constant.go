package model

const (
	TaskPending    = "Pending"
	TaskInProgress = "InProgress"
	TaskDone       = "Done"
	TaskCanceled   = "Canceled"
	TaskFailed     = "Failed"
	TaskTimedOut   = "TimedOut"
)

const (
	PipelineIdle = "Idle"
	PipelineBusy = "Busy"
)

const (
	EventBuild  = "Build"
	EventDeploy = "Deploy"
)

const (
	TaskTest    = "Test"
	TaskBuild   = "Build"
	TaskDeploy  = "Deploy"
	TaskRestart = "Restart"
)
