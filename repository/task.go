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
	filter := bson.M{"_id": input.Id, "status": "PENDING"}

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

func (r *Repository) GetDeployTasks(ctx context.Context, input *model.DeployTasksInput) ([]*model.DeployTask, error) {
	coll := r.mongoClient.Database("pipeline").Collection("deployTasks")

	filter := bson.M{}
	if input.StatusFilter != nil {
		filter["status"] = input.StatusFilter.Option
	}

	cursor, err := coll.Find(ctx, filter)

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
	filter := bson.M{"_id": input.DeployTaskId, "status": "PENDING"}
	update := bson.M{"$set": bson.M{"status": input.Status}}

	coll := r.mongoClient.Database("pipeline").Collection("deployTasks")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}
