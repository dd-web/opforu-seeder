package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Identity struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	User  primitive.ObjectID `json:"user" bson:"user,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Style string             `json:"style" bson:"style"`

	// public - mod - creator
	Role string `json:"role" bson:"role"`

	// active - suspended - banned
	Status string `json:"status" bson:"status"`

	Thread    primitive.ObjectID `json:"thread" bson:"thread"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// new empty identity ptr
func NewEmptyIdentity() *Identity {
	return &Identity{}
}

// randomize identity values
func (i *Identity) Randomize(userId primitive.ObjectID, role string) {
	i.Name = GetSlug(8, 10)
	i.Style = GetIdentityStyle()
	i.Status = GetWeightedIdentityStatus()
	i.User = userId
	i.Role = role
	i.CreatedAt = time.Now().UTC()
	i.UpdatedAt = time.Now().UTC()
}
