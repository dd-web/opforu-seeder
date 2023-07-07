package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Body   string             `json:"body" bson:"body"`
	Author primitive.ObjectID `json:"author" bson:"author"`
}
