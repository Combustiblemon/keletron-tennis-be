package ReservationModel

import (
	"combustiblemon/keletron-tennis-be/database"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"context"
	"encoding/json"
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
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Court    primitive.ObjectID `bson:"court,omitempty" mod:"trim" validate:"required"`
	Datetime string             `bson:"datetime,omitempty" mod:"trim" validate:"required"`
	Duration int                `bson:"duration,omitempty"`
	Type     string             `bson:"type,omitempty" mod:"trim"`
	Owner    primitive.ObjectID `bson:"owner,omitempty"`
	Status   string             `bson:"status,omitempty"`
	Paid     bool               `bson:"paid,omitempty"`
	Notes    string             `bson:"notes,omitempty" mod:"trim" validate:"max=600"`
	People   []string           `bson:"people,omitempty" mod:"trim" validate:"required,max=4,min=2,dive,max=30"`
}

type ReservationPartial struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Court    primitive.ObjectID `bson:"court,omitempty" mod:"trim"`
	Datetime string             `bson:"datetime,omitempty" mod:"trim"`
	Duration int                `bson:"duration,omitempty" `
	Type     string             `bson:"type,omitempty" mod:"trim"`
	Owner    primitive.ObjectID `bson:"owner,omitempty"`
	Status   string             `bson:"status,omitempty"`
	Paid     bool               `bson:"paid,omitempty" `
	Notes    string             `bson:"notes,omitempty" mod:"trim" validate:"max=600"`
	People   []string           `bson:"people,omitempty" mod:"trim" validate:"max=4,min=2,dive,max=30"`
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

func (r *Reservation) UnmarshalJSON(data []byte) error {
	// Define a temporary struct with ID as string
	type Alias Reservation
	aux := &struct {
		ID    string `json:"_id"`
		Court string
		Owner string
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	// Unmarshal into the temporary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert ID from string to ObjectID
	if aux.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(aux.ID)
		if err != nil {
			return fmt.Errorf("ID is an invalid ObjectID: %w", err)
		}
		r.ID = objectID
	}

	if aux.Court != "" {
		objectID, err := primitive.ObjectIDFromHex(aux.Court)
		if err != nil {
			return fmt.Errorf("Court is an invalid ObjectID: %w", err)
		}
		r.Court = objectID
	}

	if aux.Owner != "" {
		objectID, err := primitive.ObjectIDFromHex(aux.Owner)
		if err != nil {
			return fmt.Errorf("Owner is an invalid ObjectID: %w", err)
		}
		r.Owner = objectID
	}

	return nil
}

const COLLECTION string = "reservations"

func (r *Reservation) Save(new *Reservation) error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)

	if new != nil {
		_, err = coll.UpdateByID(context.TODO(), r.ID, Reservation{
			Court:    helpers.Condition(new.Court.IsZero(), r.Court, new.Court),
			Datetime: helpers.Condition(new.Datetime == "", r.Datetime, new.Datetime),
			Duration: helpers.Condition(new.Duration > 0, r.Duration, new.Duration),
			Type:     helpers.Condition(new.Type == "", r.Type, new.Type),
			Owner:    helpers.Condition(new.Owner.IsZero(), r.Owner, new.Owner),
			Status:   helpers.Condition(new.Status == "", r.Status, new.Status),
			Notes:    helpers.Condition(new.Notes == "", r.Notes, new.Notes),
			People:   helpers.Condition(len(new.People) == 0, r.People, new.People),
		})
	}

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

func Create(r *Reservation) (*Reservation, error) {
	client, err := database.GetClient()

	if err != nil {
		return nil, err
	}

	id := primitive.NewObjectIDFromTimestamp(time.Now())

	r.ID = id

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.InsertOne(context.TODO(), Reservation{
		ID:       id,
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

	return r, err
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
