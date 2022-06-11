package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dshurubtsov/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ErrorNoValideRefreshToken error = errors.New("indalid refresh token")

type UserModel struct {
	DB *mongo.Client
}

func (m *UserModel) CreateUser(username, password string) (string, error) {
	// context for db connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.DB.Database("testbase").Collection("users")
	err := collection.FindOne(ctx, bson.D{
		{Key: "username", Value: username},
	}).Err()
	if err != mongo.ErrNoDocuments {
		fmt.Println(err)
		return "", errors.New("user already exists")
	}

	res, err := collection.InsertOne(ctx, bson.D{
		{Key: "username", Value: username},
		{Key: "password", Value: password},
	})
	if err != nil {
		return "", err
	}

	id := fmt.Sprintf("%v", res.InsertedID)
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

	if err == mongo.ErrNoDocuments {
		return user, err
	} else if err != nil {
		return user, err
	}

	return user, nil
}

func (m *UserModel) FindByRefreshToken(refreshToken string) (models.User, error) {

	// context for db connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := models.User{}
	collection := m.DB.Database("testbase").Collection("users")

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return user, err
	}
	defer cur.Close(ctx)

	// TODO: optimize request to database for search with method FindeOne()
	for cur.Next(ctx) {
		err := cur.Decode(&user)
		if err != nil {
			return user, err
		}

		match := checkPasswordHash(refreshToken, user.RefreshToken)
		if !match {
			continue
		} else {
			return user, nil
		}
	}

	return user, ErrorNoValideRefreshToken
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
