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

func (r *Repository) UpdateDeployTask(ctx context.Context, input *model.UpdateDeployTaskInput) (primitive.ObjectID, error) {
	doc := util.StructToBsonDoc(input)

	if input.Id.IsZero() {
		input.Id = primitive.NewObjectID()
		doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	} else {
		doc["updatedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	delete(doc, "id")

	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	update := bson.M{"$set": doc}
	filter := bson.M{"_id": input.Id, "status": model.TaskPending}

	coll := r.mongoClient.Database("pipeline").Collection("deployTasks")
	result, err := coll.UpdateOne(ctx, filter, update, &opts)

	if err != nil {
		return primitive.NilObjectID, err
	}

	if result.UpsertedID == nil {
		return primitive.NilObjectID, nil
	}

	return result.UpsertedID.(primitive.ObjectID), nil
}

func (r *Repository) GetDeployTasks(ctx context.Context, input model.DeployTasksInput) ([]*model.DeployTask, error) {
	coll := r.mongoClient.Database("pipeline").Collection("deployTasks")

	filter := bson.M{}
	if input.StatusFilter != nil {
		filter["status"] = input.StatusFilter.Option
	}

	opts := options.Find().SetSort(bson.D{{"createdat", -1}})
	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	var tasks []*model.DeployTask
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) DeleteDeployTasks(ctx context.Context, id primitive.ObjectID) error {
	coll := r.mongoClient.Database("pipeline").Collection("deployTasks")
	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *Repository) UpdateDeployTaskStatus(ctx context.Context, input *model.UpdateDeployTaskStatusInput) error {
	filter := bson.M{"_id": input.DeployTaskId, "status": bson.M{"$in": bson.A{model.TaskPending, model.TaskInProgress}}}

	doc := bson.M{"status": input.Status}
	switch input.Status {
	case model.TaskInProgress:
		doc["executedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	case model.TaskDone, model.TaskFailed, model.TaskCanceled:
		doc["stoppedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("deployTasks")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdateBuildTask(ctx context.Context, input *model.UpdateBuildTaskInput) (primitive.ObjectID, error) {
	doc := util.StructToBsonDoc(input)

	if input.Id.IsZero() {
		input.Id = primitive.NewObjectID()
		doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	} else {
		doc["updatedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	delete(doc, "id")

	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	update := bson.M{"$set": doc}
	filter := bson.M{"_id": input.Id, "status": model.TaskPending}

	coll := r.mongoClient.Database("pipeline").Collection("buildTasks")
	result, err := coll.UpdateOne(ctx, filter, update, &opts)

	if err != nil {
		return primitive.NilObjectID, err
	}

	if result.UpsertedID == nil {
		return primitive.NilObjectID, nil
	}

	return result.UpsertedID.(primitive.ObjectID), nil
}

func (r *Repository) GetBuildTasks(ctx context.Context, input model.BuildTasksInput) ([]*model.BuildTask, error) {
	coll := r.mongoClient.Database("pipeline").Collection("buildTasks")

	filter := bson.M{}
	if input.StatusFilter != nil {
		filter["status"] = input.StatusFilter.Option
	}

	opts := options.Find().SetSort(bson.D{{"createdat", -1}})
	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	var tasks []*model.BuildTask
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) DeleteBuildTasks(ctx context.Context, id primitive.ObjectID) error {
	coll := r.mongoClient.Database("pipeline").Collection("buildTasks")
	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *Repository) UpdateBuildTaskStatus(ctx context.Context, input *model.UpdateBuildTaskStatusInput) error {
	filter := bson.M{"_id": input.BuildTaskId, "status": bson.M{"$in": bson.A{model.TaskPending, model.TaskInProgress}}}

	doc := bson.M{"status": input.Status}
	switch input.Status {
	case model.TaskInProgress:
		doc["executedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	case model.TaskDone, model.TaskFailed, model.TaskCanceled:
		doc["stoppedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("buildTasks")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}
