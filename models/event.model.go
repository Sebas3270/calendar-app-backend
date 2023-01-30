package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	Id          primitive.ObjectID `bson:"_id" json:"id" validate:"omitempty"`
	Title       string             `bson:"title" validate:"required" json:"title"`
	Description string             `bson:"description" validate:"omitempty" json:"description"`
	Start       time.Time          `bson:"start" validate:"required" json:"start"`
	End         time.Time          `bson:"end" validate:"required" json:"end"`
	User        string             `bson:"user" validate:"omitempty" json:"user"`
}
