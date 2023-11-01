package main

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// measurements of time
	SECONDS_IN_MINUTE = 60
	MINUTES_IN_HOUR   = 60
	HOURS_IN_DAY      = 24

	DAYS_IN_WEEK = 7
	DAYS_IN_YEAR = 365

	SECONDS_IN_HOUR = SECONDS_IN_MINUTE * MINUTES_IN_HOUR
	SECONDS_IN_DAY  = SECONDS_IN_HOUR * HOURS_IN_DAY
	SECONDS_IN_WEEK = SECONDS_IN_DAY * DAYS_IN_WEEK

	MINUTES_IN_DAY  = MINUTES_IN_HOUR * HOURS_IN_DAY
	MINUTES_IN_WEEK = MINUTES_IN_DAY * DAYS_IN_WEEK

	HOURS_IN_WEEK = HOURS_IN_DAY * DAYS_IN_WEEK
	HOURS_IN_YEAR = HOURS_IN_DAY * DAYS_IN_YEAR

	// permissions
	PUBLIC_SESSION_FIELDS   = []string{"created_at", "updated_at", "deleted_at"}
	PERSONAL_SESSION_FIELDS = []string{"_id", "account_id", "session_id", "expires"}
)

type Session struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	SessionID string             `bson:"session_id" json:"session_id"`
	AccountID primitive.ObjectID `bson:"account_id,omitempty" json:"account_id,omitempty"`

	Account *Account `bson:"account,omitempty" json:"account,omitempty"`

	Expires *time.Time `bson:"expires" json:"expires"`

	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// session from account id
func NewSessionFromAccount(account *Account) *Session {
	now := time.Now().UTC()
	exp := time.Now().Add(time.Duration(SECONDS_IN_DAY) * time.Second).UTC()

	return &Session{
		ID:        primitive.NewObjectID(),
		AccountID: account.ID,
		Account:   account,
		SessionID: uuid.NewString(),
		CreatedAt: &now,
		UpdatedAt: &now,
		Expires:   &exp,
	}
}

// tells us if the session has expired or not
func (s *Session) IsExpired() bool {
	return s.Expires.Before(time.Now().UTC())
}

// Persist Sessions
func (s *MongoStore) PersistSessions() error {
	docs := []interface{}{}

	for _, session := range s.cSessions {
		docs = append(docs, session)
	}

	return s.PersistDocuments(docs, "sessions")
}
