package UserModel

import (
	"combustiblemon/keletron-tennis-be/database"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserType interface {
	ComparePasswords(string) bool
	CompareResetKey(string) bool
	CompareSessions(string) bool
	Sanitize() UserSanitized
}

type UserSanitized struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string
	Role  string
	Email string
}

type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string
	Role        string
	Email       string
	Password    string
	ResetKey    string
	FCMTokens   []string
	Session     string
	AccountType string
}

func (u *User) ComparePasswords(password string) bool {
	return u.Password == password
}

func (u *User) CompareResetKey(resetKey string) bool {
	return u.ResetKey == resetKey
}

func (u *User) CompareSessions(session string) bool {
	return u.Session == session
}

func (u *User) Sanitize() UserSanitized {
	return UserSanitized{
		ID:    u.ID,
		Name:  u.Name,
		Role:  u.Role,
		Email: u.Email,
	}
}

const COLLECTION string = "users"

func FindOne(filter primitive.D) (*User, error) {
	client, err := database.GetClient()

	if err != nil {
		return nil, err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	var result User
	err = coll.FindOne(context.TODO(), filter).
		Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Find(filter primitive.D) (*[]User, error) {
	client, err := database.GetClient()

	if err != nil {
		return nil, err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []User
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func Create(u User) error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.InsertOne(context.TODO(), u)

	return err
}

func (u *User) Save() error {
	client, err := database.GetClient()

	if err != nil {
		return err
	}

	coll := client.Database(database.DatabaseName).Collection(COLLECTION)
	_, err = coll.UpdateByID(context.TODO(), u.ID, u)

	return err
}
