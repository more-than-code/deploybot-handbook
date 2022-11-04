package repository

import (
	"context"
	"time"

	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreateTask(ctx context.Context, input *model.CreateTaskInput) (primitive.ObjectID, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	filter := bson.M{"_id": input.PipelineId}

	doc := util.StructToBsonDoc(input.Payload)
	if input.Payload.Id.IsZero() {
		doc["id"] = primitive.NewObjectID()
	}

	doc["status"] = model.TaskPending
	doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	update := bson.M{"$push": bson.M{"tasks": doc}}
	_, err := coll.UpdateOne(ctx, filter, update)

	return doc["id"].(primitive.ObjectID), err
}

func (r *Repository) GetTask(ctx context.Context, input *model.GetTaskInput) (*model.Task, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	filter := bson.M{"_id": input.PipelineId, "tasks.id": input.Id}

	opts := options.FindOneOptions{Projection: bson.M{"tasks.$": 1}}
	var pipeline model.Pipeline
	err := coll.FindOne(ctx, filter, &opts).Decode(&pipeline)

	if err != nil {
		return nil, err
	}

	return pipeline.Tasks[0], nil
}

func (r *Repository) DeleteTask(ctx context.Context, input *model.DeleteTaskInput) error {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	filter := bson.M{"_id": input.PipelineId}
	update := bson.M{"$pull": bson.M{"tasks": bson.M{"id": input.TaskId}}}
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdateTask(ctx context.Context, input *model.UpdateTaskInput) error {
	filter := bson.M{"_id": input.PipelineId, "tasks.id": input.Id}

	doc := bson.M{}
	doc["tasks.$.updatedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	if input.Payload.Name != nil {
		doc["tasks.$.name"] = input.Payload.Name
	}
	if input.Payload.ScheduledAt != nil {
		doc["tasks.$.scheduledat"] = input.Payload.ScheduledAt
	}
	if input.Payload.Config != nil {
		doc["tasks.$.config"] = input.Payload.Config
	}
	if input.Payload.Remarks != nil {
		doc["tasks.$.remarks"] = input.Payload.Remarks
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdateTaskStatus(ctx context.Context, input *model.UpdateTaskStatusInput) error {
	filter := bson.M{"_id": input.PipelineId, "tasks.id": input.TaskId}

	doc := bson.M{"tasks.$.status": input.Payload.Status}

	switch input.Payload.Status {
	case model.TaskInProgress:
		doc["tasks.$.executedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
		doc["tasks.$.stoppedat"] = nil
	case model.TaskDone, model.TaskFailed, model.TaskCanceled:
		doc["tasks.$.stoppedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}
