package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserModel struct {
	DB *mongo.Client
}

func (m *UserModel) Insert(username, password string) (interface{}, error) {
	// context for db connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.DB.Database("testbase").Collection("users")
	res, err := collection.InsertOne(ctx, bson.D{
		{Key: "username", Value: username},
		{Key: "password", Value: password},
	})

	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}
	id := res.InsertedID

	return id, nil
}
