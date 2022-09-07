package cfutil

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConfig struct {
	ConnectionUri string
	DatabaseName  string
}

func NewMongoConfigFromEnv() (*MongoConfig, error) {
	mongoDbUrl, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		return nil, errors.New("MONGODB_URI is not set")
	}

	mongoDbName, ok := os.LookupEnv("MONGODB_DB_NAME")
	if !ok {
		return nil, errors.New("MONGODB_DB_NAME is not set")
	}

	return &MongoConfig{
		ConnectionUri: mongoDbUrl,
		DatabaseName:  mongoDbName,
	}, nil
}

type MongoUtil struct {
	mgo         *mongo.Client
	collections map[string]*mongo.Collection
	config      *MongoConfig
}

var collections = map[string]*mongo.Collection{}
var mgo *mongo.Client
var mgoConfig *MongoConfig

func GetMongoConfig() (*MongoConfig, error) {
	if mgoConfig == nil {
		cfg, err := NewMongoConfigFromEnv()
		if err != nil {
			return nil, err
		}
		mgoConfig = cfg
	}

	return mgoConfig, nil
}

func GetMongoClient(ctx context.Context) (*mongo.Client, error) {
	if mgo == nil {
		cfg, err := GetMongoConfig()
		if err != nil {
			return nil, err
		}

		mgoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.ConnectionUri))
		if err != nil {
			return nil, err
		}
		mgo = mgoClient
	}

	return mgo, nil
}

func GetCollection(ctx context.Context, collectionName string) (*mongo.Collection, error) {
	if mgo == nil || mgo.Ping(ctx, readpref.Primary()) != nil {
		var err error
		mgo, err = GetMongoClient(ctx)
		if err != nil {
			return nil, err
		}
	}

	if _, ok := collections[collectionName]; !ok {
		cfg, err := GetMongoConfig()
		if err != nil {
			return nil, err
		}

		collections[collectionName] = mgo.Database(cfg.DatabaseName).Collection(collectionName)
	}

	return collections[collectionName], nil
}
