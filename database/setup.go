package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const DatabaseName string = "dev"

func Setup() error {
	bsonOpts := &options.BSONOptions{
		UseJSONStructTags:       true,
		NilSliceAsEmpty:         true,
		ErrorOnInlineDuplicates: true,
		NilMapAsEmpty:           true,
		NilByteSliceAsEmpty:     true,
	}

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		return fmt.Errorf("Set your 'MONGODB_URI' environment variable. " +
			"See: " +
			"www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	_client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri).SetBSONOptions(bsonOpts))

	client = _client

	return err
}

func GetClient() (*mongo.Client, error) {
	if client != nil {
		return client, nil
	}
	err := Setup()

	if err != nil {
		return nil, err
	}

	return client, nil
}

func Teardown() error {
	err := client.Disconnect(context.TODO())

	return err
}
