package common

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func GetDBCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

func OpenDB() error {
	uri := "mongodb://localhost:27017"
	if uri == "" {
		return errors.New("please set the 'MONGO_URI' in the env")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	db = client.Database("local")
	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}
