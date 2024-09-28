package ReservationModel

import (
	"combustiblemon/keletron-tennis-be/database"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationType interface {
	Sanitize() ReservationSanitized
}

type ReservationSanitized struct {
	ID       primitive.ObjectID `bson:"_id"`
	Court    string
	Datetime string
	Duration int
	Type     string
}

type Reservation struct {
	ID       primitive.ObjectID `bson:"_id"`
	Court    string
	Datetime string
	Duration int
	Type     string
	Owner    string
	Status   string
	Paid     bool
	Notes    string
	People   []string
}

func (r *Reservation) Sanitize() ReservationSanitized {
	return ReservationSanitized{
		Court:    r.Court,
		Datetime: r.Datetime,
		Duration: r.Duration,
		Type:     r.Type,
	}
}

const COLLECTION string = "reservations"

func (r *Reservation) Save() error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.UpdateByID(context.TODO(), r.ID, r)

	return err
}

func FindOne(filter primitive.D) (*Reservation, error) {
	client, err := database.GetClient()

	if err != nil {
		return nil, err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	var result Reservation
	err = coll.FindOne(context.TODO(), filter).
		Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Find(filter primitive.D) (*[]Reservation, error) {
	client, err := database.GetClient()

	if err != nil {
		return nil, err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []Reservation
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func Create(u Reservation) error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.InsertOne(context.TODO(), u)

	return err
}

func DeleteOne(id string) error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}})

	return err
}

func DeleteMany() error {
	return fmt.Errorf("Not implemented")
}
