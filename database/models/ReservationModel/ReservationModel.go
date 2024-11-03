package ReservationModel

import (
	"combustiblemon/keletron-tennis-be/database"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationType interface {
	Sanitize() ReservationSanitized
	SanitizeOwner() ReservationSanitizedOwner
}

type ReservationSanitized struct {
	ID       primitive.ObjectID `bson:"_id"`
	Court    primitive.ObjectID
	Datetime string
	Duration int
	Type     string
}

type ReservationSanitizedOwner struct {
	ID       primitive.ObjectID `bson:"_id"`
	Court    primitive.ObjectID
	Datetime string
	Duration int
	Type     string
	Notes    string
	People   []string
	Status   string
}

type Reservation struct {
	ID       primitive.ObjectID `bson:"_id"`
	Court    primitive.ObjectID `mod:"trim" validate:"required,mongodb"`
	Datetime string             `mod:"trim" validate:"required"`
	Duration int                ``
	Type     string             `mod:"trim"`
	Owner    primitive.ObjectID ``
	Status   string             ``
	Paid     bool               ``
	Notes    string             `mod:"trim" validate:"max=600"`
	People   []string           `mod:"trim" validate:"required,max=4,min=2,dive,max=30"`
}

func (r *Reservation) Sanitize() ReservationSanitized {
	return ReservationSanitized{
		Court:    r.Court,
		Datetime: r.Datetime,
		Duration: r.Duration,
		Type:     r.Type,
	}
}

func (r *Reservation) SanitizeOwner() ReservationSanitizedOwner {
	return ReservationSanitizedOwner{
		ID:       r.ID,
		Court:    r.Court,
		Datetime: r.Datetime,
		Duration: r.Duration,
		Type:     r.Type,
		Notes:    r.Notes,
		People:   r.People,
		Status:   r.Status,
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

func (r *Reservation) Delete() error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: r.ID}})

	return err
}

func (r *Reservation) Date() time.Time {
	return helpers.ParseDate(r.Datetime)
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

func Create(r Reservation) error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.InsertOne(context.TODO(), Reservation{
		Court:    r.Court,
		Datetime: r.Datetime,
		Duration: r.Duration,
		Type:     r.Type,
		Owner:    r.Owner,
		Status:   r.Status,
		Paid:     false,
		Notes:    r.Notes,
		People:   r.People,
	})

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
