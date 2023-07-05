package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Thread struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title"`
}
