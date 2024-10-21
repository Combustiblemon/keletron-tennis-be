package AnnouncementModel

import (
	"combustiblemon/keletron-tennis-be/database"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Announcement struct {
	ID         primitive.ObjectID `bson:"_id"`
	Title      string
	ValidUntil string
	Visible    string
}

const COLLECTION string = "announcements"

func FindOne(filter primitive.D) (*Announcement, error) {
	client, err := database.GetClient()

	if err != nil {
		return nil, err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	var result Announcement
	err = coll.FindOne(context.TODO(), filter).
		Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Find(filter primitive.D) (*[]Announcement, error) {
	client, err := database.GetClient()

	if err != nil {
		return nil, err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []Announcement
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func Create(u Announcement) error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.InsertOne(context.TODO(), u)

	return err
}
func (c *Announcement) Save() error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.UpdateByID(context.TODO(), c.ID, c)

	return err
}
