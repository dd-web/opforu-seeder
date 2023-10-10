package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Identity struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Account primitive.ObjectID `json:"account" bson:"account,omitempty"`

	Name  string `json:"name" bson:"name"`
	Style string `json:"style" bson:"style"`

	Role   ThreadRole     `json:"role" bson:"role"`
	Status IdentityStatus `json:"status" bson:"status"`

	Thread primitive.ObjectID `json:"thread" bson:"thread"`

	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// new empty identity ptr
func NewEmptyIdentity() *Identity {
	ts := time.Now().UTC()
	return &Identity{
		CreatedAt: &ts,
		UpdatedAt: &ts,
	}
}

// randomize identity values
func (i *Identity) Randomize(userId primitive.ObjectID, role ThreadRole) {
	ts := time.Now().UTC()
	i.ID = primitive.NewObjectID()
	i.Name = GetSlug(8, 10)
	i.Style = GetIdentityStyle()
	i.Status = GetWeightedIdentityStatus()
	i.Account = userId
	i.Role = role
	i.UpdatedAt = &ts
}

// weighted identity role
func (i *Identity) SetWeightedRole() {
	if RandomIntBetween(0, 100) < 95 {
		i.Role = ThreadRoleUser
	} else {
		i.Role = ThreadRoleMod
	}
}

// Generate Identity for a thread
func (s *MongoStore) GenerateThreadIdentity(userId primitive.ObjectID, role ThreadRole) *Identity {
	identity := NewEmptyIdentity()
	identity.Randomize(userId, role)
	return identity
}

// Return a user's identity for a thread - or create one if it doesn't exist
func (s *MongoStore) GetUserThreadIdentity(userId, threadId primitive.ObjectID) *Identity {
	threadIdentity := s.cUserThreadIdentitys[threadId][userId]

	if threadIdentity == nil {
		threadIdentity = s.GenerateThreadIdentity(userId, ThreadRoleUser)
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

// enums
type IdentityStatus string

const (
	IdentityStatusUnknown   IdentityStatus = "unknown"
	IdentityStatusActive    IdentityStatus = "active"
	IdentityStatusSuspended IdentityStatus = "suspended"
	IdentityStatusBanned    IdentityStatus = "banned"
	IdentityStatusDeleted   IdentityStatus = "deleted"
)

func (i IdentityStatus) String() string {
	switch i {
	case IdentityStatusActive:
		return "active"
	case IdentityStatusSuspended:
		return "suspended"
	case IdentityStatusBanned:
		return "banned"
	case IdentityStatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}
