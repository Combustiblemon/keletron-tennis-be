package CourtModel

import (
	"combustiblemon/keletron-tennis-be/database"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservedTimes struct {
	ID              primitive.ObjectID `bson:"_id"`
	StartTime       string
	Duration        int
	Type            string
	Repeat          string
	Notes           string
	Days            []string
	datesNotApplied []string
}

type ReservationsInfo struct {
	ID            primitive.ObjectID `bson:"_id"`
	StartTime     string
	EndTime       string
	Duration      int
	ReservedTimes []ReservedTimes
}

type Court struct {
	ID               primitive.ObjectID `bson:"_id"`
	Name             string
	Type             string
	ReservationsInfo ReservationsInfo
}

const COLLECTION string = "courts"

func FindOne(filter primitive.D) (*Court, error) {
	client, err := database.GetClient()

	if err != nil {
		return nil, err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	var result Court
	err = coll.FindOne(context.TODO(), filter).
		Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Find(filter primitive.D) (*[]Court, error) {
	client, err := database.GetClient()

	if err != nil {
		return nil, err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []Court
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func Create(u Court) error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.InsertOne(context.TODO(), u)

	return err
}
func (c *Court) Save() error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.UpdateByID(context.TODO(), c.ID, c)

	return err
}
