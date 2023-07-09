package main

import (
	"context"
	"fmt"
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
	i.ID = primitive.NewObjectID()
	i.Name = GetSlug(8, 10)
	i.Style = GetIdentityStyle()
	i.Status = GetWeightedIdentityStatus()
	i.User = userId
	i.Role = role
	i.CreatedAt = time.Now().UTC()
	i.UpdatedAt = time.Now().UTC()
}

// weighted identity role
func (i *Identity) SetWeightedRole() {
	if RandomIntBetween(0, 100) < 95 {
		i.Role = "public"
	} else {
		i.Role = "mod"
	}
}

// Generate Identity for a thread
func (s *MongoStore) GenerateThreadIdentity(userId primitive.ObjectID, role string) *Identity {
	identity := NewEmptyIdentity()
	identity.Randomize(userId, role)
	return identity
}

// Return a user's identity for a thread - or create one if it doesn't exist
func (s *MongoStore) GetUserThreadIdentity(userId, threadId primitive.ObjectID) *Identity {
	threadIdentity := s.cUserThreadIdentitys[threadId][userId]

	if threadIdentity == nil {
		threadIdentity = s.GenerateThreadIdentity(userId, "public")
		threadIdentity.SetWeightedRole()
		threadIdentity.Thread = threadId
		s.cUserThreadIdentitys[threadId][userId] = threadIdentity
		s.cIdentites = append(s.cIdentites, threadIdentity)
	}

	return threadIdentity
}

// Persist Identitys
func (s *MongoStore) PersistIdentities() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	docs := []interface{}{}

	for _, identity := range s.cIdentites {
		docs = append(docs, identity)
	}

	identityCol := s.DB.Collection("identities")
	response, err := identityCol.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	fmt.Printf(" - Persisted %d %s documents to database\n", len(response.InsertedIDs), "identities")

	return nil
}
