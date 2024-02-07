package main

import (
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

// New identity
func NewIdentity(account, thread primitive.ObjectID, role ThreadRole) *Identity {
	ts := time.Now().UTC()
	return &Identity{
		ID:        primitive.NewObjectID(),
		Account:   account,
		Name:      GetSlug(8, 10),
		Style:     GetIdentityStyle(),
		Role:      role,
		Status:    GetWeightedIdentityStatus(),
		Thread:    thread,
		CreatedAt: &ts,
		UpdatedAt: &ts,
	}
}

// Return a user's identity for a thread - or create one if it doesn't exist
func (s *MongoStore) GetUserThreadIdentity(account, thread primitive.ObjectID) *Identity {
	if identity, ok := s.cUserThreadIdentitys[thread][account]; ok {
		return identity
	} else {
		identity := NewIdentity(account, thread, GetWeightedThreadRole())
		s.cUserThreadIdentitys[thread][account] = identity
		s.cIdentites = append(s.cIdentites, identity)
		return identity
	}
}

// Persist Identitys
func (s *MongoStore) PersistIdentities() error {
	docs := []interface{}{}

	for _, identity := range s.cIdentites {
		docs = append(docs, identity)
	}

	return s.PersistDocuments(docs, "identities")
}

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
