package repository

import (
	"context"
	"time"

	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *Repository) CreatePipeline(ctx context.Context, input *model.CreatePipelineInput) (primitive.ObjectID, error) {
	doc := util.StructToBsonDoc(input.Payload)

	doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	result, err := coll.InsertOne(ctx, doc)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *Repository) DeletePipeline(ctx context.Context, id primitive.ObjectID) error {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *Repository) GetPipelines(ctx context.Context) ([]*model.Pipeline, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")

	filter := bson.M{}

	opts := options.Find().SetSort(bson.D{{"updatedat", -1}})
	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	var pipelines []*model.Pipeline
	if err = cursor.All(ctx, &pipelines); err != nil {
		return nil, err
	}

	return pipelines, nil
}

func (r *Repository) GetPipeline(ctx context.Context, input *model.GetPipelineInput) (*model.Pipeline, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")

	filter := bson.M{"name": input.Name}

	var pipeline model.Pipeline
	err := coll.FindOne(ctx, filter).Decode(&pipeline)

	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (r *Repository) CreatePipelineTask(ctx context.Context, input *model.CreatePipelineTaskInput) (primitive.ObjectID, error) {
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

func (r *Repository) GetPipelineTask(ctx context.Context, input *model.GetPipelineTaskInput) (*model.Task, error) {
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

func (r *Repository) DeletePipelineTask(ctx context.Context, input *model.DeletePipelineTaskInput) error {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	filter := bson.M{"_id": input.PipelineId}
	update := bson.M{"$pull": bson.M{"tasks": bson.M{"id": input.TaskId}}}
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdatePipelineTaskStatus(ctx context.Context, input *model.UpdatePipelineTaskStatusInput) error {
	filter := bson.M{"_id": input.PipelineId, "tasks.id": input.TaskId}

	doc := bson.M{"tasks.$.status": input.Payload.Status}

	switch input.Payload.Status {
	case model.TaskInProgress:
		doc["tasks.$.executedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	case model.TaskDone, model.TaskFailed, model.TaskCanceled:
		doc["tasks.$.stoppedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}
