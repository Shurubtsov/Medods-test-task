package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dshurubtsov/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserModel struct {
	DB *mongo.Client
}

func (m *UserModel) CreateUser(username, password string) (interface{}, error) {
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

func (m *UserModel) UpdateUserToken(id, refreshToken string) error {
	// context for db connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.DB.Database("testbase").Collection("users")
	objId, _ := primitive.ObjectIDFromHex(id)

	_, err := collection.UpdateOne(ctx, bson.D{
		{Key: "_id", Value: objId},
	}, bson.D{{Key: "$set", Value: bson.D{{Key: "refresh_token", Value: refreshToken}}}})
	if err != nil {
		return err
	}

	if err == mongo.ErrNoDocuments {
		return err
	} else if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) FindById(id string) (models.User, error) {
	// context for db connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := models.User{}

	collection := m.DB.Database("testbase").Collection("users")
	objId, _ := primitive.ObjectIDFromHex(id)

	fmt.Println("id is :", objId)

	err := collection.FindOne(ctx, bson.D{
		{Key: "_id", Value: objId},
	}).Decode(&user)

	fmt.Println(user)

	if err == mongo.ErrNoDocuments {
		return models.User{}, err
	} else if err != nil {
		return models.User{}, err
	}

	return user, nil
}
