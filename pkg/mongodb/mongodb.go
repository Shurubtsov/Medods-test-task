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

func (m *UserModel) Login(username, password string) (string, error) {
	// context for db connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.D

	collection := m.DB.Database("testbase").Collection("users")
	err := collection.FindOne(ctx, bson.D{
		{Key: "username", Value: username},
		{Key: "password", Value: password},
	}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return "User does not exist", err
	} else if err != nil {
		return err.Error(), err
	}

	return "u have been logged", nil
}
