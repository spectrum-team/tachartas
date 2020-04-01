package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Description string             `bson:"description" json:"description"`
}
