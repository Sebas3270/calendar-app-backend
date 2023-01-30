package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Name     string             `bson:"name" validate:"required" json:"name"`
	Email    string             `bson:"email" validate:"required,email" json:"email"`
	Password string             `bson:"password" validate:"required" json:"password"`
}

type LogIn struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
