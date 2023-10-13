package main

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	SessionID string             `bson:"session_id" json:"session_id"`
	AccountID primitive.ObjectID `bson:"account_id" json:"account_id"`

	Active bool `bson:"active" json:"active"`

	Expiry    *time.Time `bson:"expiry" json:"expiry"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// session from account id
func NewSessionFromAccount(id primitive.ObjectID) *Session {
	now := time.Now().UTC()
	expiry := time.Now().Add(3600 * time.Second).UTC()
	nid := uuid.New()
	return &Session{
		ID:        primitive.NewObjectID(),
		AccountID: id,
		SessionID: nid.String(),
		Active:    true,
		CreatedAt: &now,
		UpdatedAt: &now,
		Expiry:    &expiry,
	}
}

// tells us if the session has expired or not
func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now().UTC())
}

// Persist Sessions
func (s *MongoStore) PersistSessions() error {
	docs := []interface{}{}

	for _, session := range s.cSessions {
		docs = append(docs, session)
	}

	return s.PersistDocuments(docs, "sessions")
}
