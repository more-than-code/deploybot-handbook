package repository

import (
	"context"
	"reflect"

	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	tM := reflect.TypeOf(bson.M{})

	reg := bson.NewRegistryBuilder().RegisterTypeMapEntry(bsontype.EmbeddedDocument, tM).Build()

	clientOpts := options.Client().SetRegistry(reg).ApplyURI(cfg.MongoUri)

	mongoClient, err := mongo.Connect(context.TODO(), clientOpts)
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
