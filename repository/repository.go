package repository

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CacheGetHandlerFunc func(ctx context.Context) interface{}
type CacheUpdateHandlerFunc func(ctx context.Context)
type CacheClearHandlerFunc func(ctx context.Context)

type Repository struct {
	mongoClient *mongo.Client
}

type Config struct {
	MongoUri string `envconfig:"MONGODB_URI"`
}

func NewRepository() (*Repository, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoUri))
	if err != nil {
		panic(err)
	}

	return &Repository{mongoClient: mongoClient}, nil
}

func (r *Repository) Disconnect() {
	if err := r.mongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
