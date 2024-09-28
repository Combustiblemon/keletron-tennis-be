package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

//revive:next-line:var-declaration
var DatabaseName = "dev"

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

	// coll := client.Database("sample_mflix").Collection("movies")
	// title := "Back to the Future"

	// var result bson.M
	// err = coll.FindOne(context.TODO(), bson.D{{Key: "title", Value: title}}).
	// 	Decode(&result)
	// if err == mongo.ErrNoDocuments {
	// 	fmt.Printf("No document was found with the title %s\n", title)
	// 	return err
	// }
	// if err != nil {
	// 	panic(err)
	// }

	// jsonData, err := json.MarshalIndent(result, "", "    ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", jsonData)

	// return err
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
