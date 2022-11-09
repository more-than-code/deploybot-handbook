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
	doc["status"] = model.PipelineIdle

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

func (r *Repository) GetPipelines(ctx context.Context, input model.GetPipelinesInput) ([]*model.Pipeline, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")

	filter := bson.M{}
	if input.RepoWatched != nil {
		filter["repowatched"] = input.RepoWatched
	}
	if input.AutoRun != nil {
		filter["autorun"] = input.AutoRun
	}

	opts := options.Find().SetSort(bson.D{{"executedat", -1}})
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

func (r *Repository) GetPipeline(ctx context.Context, input model.GetPipelineInput) (*model.Pipeline, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")

	filter := bson.M{}
	if input.Id != nil {
		filter["_id"] = input.Id
	}
	if input.Name != nil {
		filter["name"] = input.Name
	}

	taskFilter := bson.M{}
	if input.TaskFilter.UpstreamTaskId != nil {
		taskFilter["upstreamtaskid"] = input.TaskFilter.UpstreamTaskId

	}
	if input.TaskFilter.AutoRun != nil {
		taskFilter["autorun"] = input.TaskFilter.AutoRun
	}

	opts := options.FindOneOptions{}
	if len(taskFilter) > 0 {
		opts.SetProjection(bson.M{"tasks": bson.M{"$elemMatch": taskFilter}})
	}

	var pipeline model.Pipeline
	err := coll.FindOne(ctx, filter, &opts).Decode(&pipeline)

	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (r *Repository) UpdatePipeline(ctx context.Context, input model.UpdatePipelineInput) error {
	filter := bson.M{"_id": input.Id}

	doc := bson.M{}
	doc["updatedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	if input.Payload.Name != nil {
		doc["name"] = input.Payload.Name
	}
	if input.Payload.ScheduledAt != nil {
		doc["scheduledat"] = input.Payload.ScheduledAt
	}
	if input.Payload.AutoRun != nil {
		doc["autorun"] = input.Payload.AutoRun
	}
	if input.Payload.Arguments != nil {
		doc["arguments"] = input.Payload.Arguments
	}
	if input.Payload.RepoWatched != nil {
		doc["repowatched"] = input.Payload.RepoWatched
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdatePipelineStatus(ctx context.Context, input model.UpdatePipelineStatusInput) error {
	filter := bson.M{"_id": input.PipelineId}

	doc := bson.M{"status": input.Payload.Status}

	switch input.Payload.Status {
	case model.PipelineBusy:
		doc["executedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
		doc["stoppedat"] = nil
	case model.PipelineIdle:
		doc["stoppedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}
