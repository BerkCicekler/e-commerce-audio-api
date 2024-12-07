package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Value string             `json:"value" bson:"value"`
}
