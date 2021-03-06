package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UUID         string             `bson:"uuid,omitempty" json:"uuid,omitempty"`
	Username     string             `bson:"username,omitempty" json:"username,omitempty"`
	Password     string             `bson:"password,omitempty" json:"password,omitempty"`
	RefreshToken string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
}
