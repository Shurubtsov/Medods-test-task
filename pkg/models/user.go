package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username,omitempty"`
	Password     string             `bson:"password,omitempty"`
	RefreshToken string             `bson:"refresh_token,omitempty"`
}
